package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	APIPort    string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("SHOWBIZ_DB_HOST", "localhost"),
		DBPort:     getEnv("SHOWBIZ_DB_PORT", "3306"),
		DBUser:     getEnv("SHOWBIZ_DB_USER", "showbiz"),
		DBPassword: getEnv("SHOWBIZ_DB_PASSWORD", "showbiz_dev"),
		DBName:     getEnv("SHOWBIZ_DB_NAME", "showbiz"),
		JWTSecret:  getEnv("SHOWBIZ_JWT_SECRET", "dev-secret-do-not-use-in-production"),
		APIPort:    getEnv("SHOWBIZ_API_PORT", "8080"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
