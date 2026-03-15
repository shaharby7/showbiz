package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shaharby7/showbiz/services/api/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: migrate <up|down|status>")
		os.Exit(1)
	}

	cfg := config.Load()
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	ensureMigrationsTable(db)

	switch os.Args[1] {
	case "up":
		migrateUp(db)
	case "down":
		migrateDown(db)
	case "status":
		migrateStatus(db)
	default:
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func ensureMigrationsTable(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("failed to create migrations table: %v", err)
	}
}

func getMigrationFiles(direction string) []string {
	pattern := "migrations/*." + direction + ".sql"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		log.Fatalf("failed to glob migration files: %v", err)
	}
	sort.Strings(matches)
	return matches
}

func getAppliedMigrations(db *sql.DB) map[string]bool {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		log.Fatalf("failed to query migrations: %v", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		rows.Scan(&version)
		applied[version] = true
	}
	return applied
}

func versionFromFile(filename string) string {
	base := filepath.Base(filename)
	parts := strings.SplitN(base, ".", 2)
	return parts[0]
}

func migrateUp(db *sql.DB) {
	files := getMigrationFiles("up")
	applied := getAppliedMigrations(db)

	count := 0
	for _, file := range files {
		version := versionFromFile(file)
		if applied[version] {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("failed to read %s: %v", file, err)
		}

		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.Exec(stmt); err != nil {
				log.Fatalf("failed to execute %s: %v\nStatement: %s", file, err, stmt)
			}
		}

		db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
		fmt.Printf("applied: %s\n", version)
		count++
	}

	if count == 0 {
		fmt.Println("no new migrations to apply")
	} else {
		fmt.Printf("applied %d migration(s)\n", count)
	}
}

func migrateDown(db *sql.DB) {
	files := getMigrationFiles("down")
	applied := getAppliedMigrations(db)

	// Find the last applied migration
	upFiles := getMigrationFiles("up")
	var lastApplied string
	for i := len(upFiles) - 1; i >= 0; i-- {
		v := versionFromFile(upFiles[i])
		if applied[v] {
			lastApplied = v
			break
		}
	}

	if lastApplied == "" {
		fmt.Println("no migrations to rollback")
		return
	}

	// Find the matching down file
	for _, file := range files {
		if versionFromFile(file) == lastApplied {
			content, err := os.ReadFile(file)
			if err != nil {
				log.Fatalf("failed to read %s: %v", file, err)
			}

			statements := strings.Split(string(content), ";")
			for _, stmt := range statements {
				stmt = strings.TrimSpace(stmt)
				if stmt == "" {
					continue
				}
				if _, err := db.Exec(stmt); err != nil {
					log.Fatalf("failed to execute %s: %v\nStatement: %s", file, err, stmt)
				}
			}

			db.Exec("DELETE FROM schema_migrations WHERE version = ?", lastApplied)
			fmt.Printf("rolled back: %s\n", lastApplied)
			return
		}
	}

	fmt.Printf("no down migration found for %s\n", lastApplied)
}

func migrateStatus(db *sql.DB) {
	files := getMigrationFiles("up")
	applied := getAppliedMigrations(db)

	for _, file := range files {
		version := versionFromFile(file)
		status := "pending"
		if applied[version] {
			status = "applied"
		}
		fmt.Printf("[%s] %s\n", status, version)
	}
}
