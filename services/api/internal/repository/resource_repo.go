package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/shaharby7/showbiz/services/api/internal/model"
)

type ResourceRepo struct {
	db *sql.DB
}

func NewResourceRepo(db *sql.DB) *ResourceRepo {
	return &ResourceRepo{db: db}
}

func (r *ResourceRepo) CreateResource(ctx context.Context, res *model.Resource) error {
	valuesJSON, err := json.Marshal(res.Values)
	if err != nil {
		return fmt.Errorf("marshal values: %w", err)
	}

	query := `INSERT INTO resources (id, name, project_id, connection_id, resource_type, values_json, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())`
	_, err = r.db.ExecContext(ctx, query, res.ID, res.Name, res.ProjectID, res.ConnectionID, res.ResourceType, valuesJSON, res.Status)
	if err != nil {
		return fmt.Errorf("create resource: %w", err)
	}
	return nil
}

func (r *ResourceRepo) GetResourceByID(ctx context.Context, id string) (*model.Resource, error) {
	query := `SELECT id, name, project_id, connection_id, resource_type, values_json, status, created_at, updated_at
		FROM resources WHERE id = ?`
	res := &model.Resource{}
	var valuesJSON []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.Name, &res.ProjectID, &res.ConnectionID, &res.ResourceType, &valuesJSON, &res.Status, &res.CreatedAt, &res.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get resource by id: %w", err)
	}

	if err := json.Unmarshal(valuesJSON, &res.Values); err != nil {
		return nil, fmt.Errorf("unmarshal values: %w", err)
	}

	return res, nil
}

func (r *ResourceRepo) GetResourceByProjectAndName(ctx context.Context, projectID, name string) (*model.Resource, error) {
	query := `SELECT id, name, project_id, connection_id, resource_type, values_json, status, created_at, updated_at
		FROM resources WHERE project_id = ? AND name = ?`
	res := &model.Resource{}
	var valuesJSON []byte
	err := r.db.QueryRowContext(ctx, query, projectID, name).Scan(
		&res.ID, &res.Name, &res.ProjectID, &res.ConnectionID, &res.ResourceType, &valuesJSON, &res.Status, &res.CreatedAt, &res.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get resource by project and name: %w", err)
	}

	if err := json.Unmarshal(valuesJSON, &res.Values); err != nil {
		return nil, fmt.Errorf("unmarshal values: %w", err)
	}

	return res, nil
}

func (r *ResourceRepo) ListResourcesByProjectID(ctx context.Context, projectID string, cursor string, limit int) ([]*model.Resource, string, error) {
	if limit <= 0 {
		limit = 20
	}

	var rows *sql.Rows
	var err error

	if cursor == "" {
		query := `SELECT id, name, project_id, connection_id, resource_type, values_json, status, created_at, updated_at
			FROM resources WHERE project_id = ? ORDER BY id ASC LIMIT ?`
		rows, err = r.db.QueryContext(ctx, query, projectID, limit+1)
	} else {
		query := `SELECT id, name, project_id, connection_id, resource_type, values_json, status, created_at, updated_at
			FROM resources WHERE project_id = ? AND id > ? ORDER BY id ASC LIMIT ?`
		rows, err = r.db.QueryContext(ctx, query, projectID, cursor, limit+1)
	}
	if err != nil {
		return nil, "", fmt.Errorf("list resources: %w", err)
	}
	defer rows.Close()

	var resources []*model.Resource
	for rows.Next() {
		res := &model.Resource{}
		var valuesJSON []byte
		if err := rows.Scan(&res.ID, &res.Name, &res.ProjectID, &res.ConnectionID, &res.ResourceType, &valuesJSON, &res.Status, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, "", fmt.Errorf("scan resource: %w", err)
		}
		if err := json.Unmarshal(valuesJSON, &res.Values); err != nil {
			return nil, "", fmt.Errorf("unmarshal values: %w", err)
		}
		resources = append(resources, res)
	}
	if err := rows.Err(); err != nil {
		return nil, "", fmt.Errorf("list resources rows: %w", err)
	}

	var nextCursor string
	if len(resources) > limit {
		nextCursor = resources[limit-1].ID
		resources = resources[:limit]
	}

	return resources, nextCursor, nil
}

func (r *ResourceRepo) UpdateResource(ctx context.Context, id string, values map[string]interface{}, status string) error {
	valuesJSON, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("marshal values: %w", err)
	}

	query := `UPDATE resources SET values_json = ?, status = ?, updated_at = NOW() WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, valuesJSON, status, id)
	if err != nil {
		return fmt.Errorf("update resource: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update resource rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("resource not found")
	}
	return nil
}

func (r *ResourceRepo) DeleteResource(ctx context.Context, id string) error {
	query := `DELETE FROM resources WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete resource: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete resource rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("resource not found")
	}
	return nil
}
