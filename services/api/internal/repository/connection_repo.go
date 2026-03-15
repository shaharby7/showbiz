package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/shaharby7/showbiz/services/api/internal/model"
)

type ConnectionRepo struct {
	db *sql.DB
}

func NewConnectionRepo(db *sql.DB) *ConnectionRepo {
	return &ConnectionRepo{db: db}
}

func (r *ConnectionRepo) CreateConnection(ctx context.Context, conn *model.Connection) error {
	credentialsBlob, err := encodeCredentials(conn.Credentials)
	if err != nil {
		return fmt.Errorf("encode credentials: %w", err)
	}

	configJSON, err := json.Marshal(conn.Config)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	query := `INSERT INTO connections (id, name, project_id, provider, credentials, config, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`
	_, err = r.db.ExecContext(ctx, query, conn.ID, conn.Name, conn.ProjectID, conn.Provider, credentialsBlob, configJSON)
	if err != nil {
		return fmt.Errorf("create connection: %w", err)
	}
	return nil
}

func (r *ConnectionRepo) GetConnectionByID(ctx context.Context, id string) (*model.Connection, error) {
	query := `SELECT id, name, project_id, provider, config, created_at, updated_at
		FROM connections WHERE id = ?`
	conn := &model.Connection{}
	var configJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conn.ID, &conn.Name, &conn.ProjectID, &conn.Provider, &configJSON, &conn.CreatedAt, &conn.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get connection by id: %w", err)
	}

	if err := json.Unmarshal(configJSON, &conn.Config); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return conn, nil
}

func (r *ConnectionRepo) GetConnectionByIDWithCredentials(ctx context.Context, id string) (*model.Connection, error) {
	query := `SELECT id, name, project_id, provider, credentials, config, created_at, updated_at
		FROM connections WHERE id = ?`
	conn := &model.Connection{}
	var credentialsBlob string
	var configJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conn.ID, &conn.Name, &conn.ProjectID, &conn.Provider, &credentialsBlob, &configJSON, &conn.CreatedAt, &conn.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get connection by id with credentials: %w", err)
	}

	creds, err := decodeCredentials(credentialsBlob)
	if err != nil {
		return nil, fmt.Errorf("decode credentials: %w", err)
	}
	conn.Credentials = creds

	if err := json.Unmarshal(configJSON, &conn.Config); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return conn, nil
}

func (r *ConnectionRepo) ListConnectionsByProjectID(ctx context.Context, projectID string, cursor string, limit int) ([]*model.Connection, string, error) {
	if limit <= 0 {
		limit = 20
	}

	var rows *sql.Rows
	var err error

	if cursor == "" {
		query := `SELECT id, name, project_id, provider, config, created_at, updated_at
			FROM connections WHERE project_id = ? ORDER BY id ASC LIMIT ?`
		rows, err = r.db.QueryContext(ctx, query, projectID, limit+1)
	} else {
		query := `SELECT id, name, project_id, provider, config, created_at, updated_at
			FROM connections WHERE project_id = ? AND id > ? ORDER BY id ASC LIMIT ?`
		rows, err = r.db.QueryContext(ctx, query, projectID, cursor, limit+1)
	}
	if err != nil {
		return nil, "", fmt.Errorf("list connections: %w", err)
	}
	defer rows.Close()

	var conns []*model.Connection
	for rows.Next() {
		conn := &model.Connection{}
		var configJSON []byte
		if err := rows.Scan(&conn.ID, &conn.Name, &conn.ProjectID, &conn.Provider, &configJSON, &conn.CreatedAt, &conn.UpdatedAt); err != nil {
			return nil, "", fmt.Errorf("scan connection: %w", err)
		}
		if err := json.Unmarshal(configJSON, &conn.Config); err != nil {
			return nil, "", fmt.Errorf("unmarshal config: %w", err)
		}
		conns = append(conns, conn)
	}
	if err := rows.Err(); err != nil {
		return nil, "", fmt.Errorf("list connections rows: %w", err)
	}

	var nextCursor string
	if len(conns) > limit {
		nextCursor = conns[limit-1].ID
		conns = conns[:limit]
	}

	return conns, nextCursor, nil
}

func (r *ConnectionRepo) UpdateConnection(ctx context.Context, id string, config map[string]interface{}) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	query := `UPDATE connections SET config = ?, updated_at = NOW() WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, configJSON, id)
	if err != nil {
		return fmt.Errorf("update connection: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update connection rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("connection not found")
	}
	return nil
}

func (r *ConnectionRepo) DeleteConnection(ctx context.Context, id string) error {
	query := `DELETE FROM connections WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete connection: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete connection rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("connection not found")
	}
	return nil
}

func (r *ConnectionRepo) HasResources(ctx context.Context, connectionID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM resources WHERE connection_id = ?)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, connectionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check resources: %w", err)
	}
	return exists, nil
}

func encodeCredentials(creds map[string]interface{}) (string, error) {
	if creds == nil {
		return base64.StdEncoding.EncodeToString([]byte("{}")), nil
	}
	data, err := json.Marshal(creds)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func decodeCredentials(encoded string) (map[string]interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	var creds map[string]interface{}
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	return creds, nil
}
