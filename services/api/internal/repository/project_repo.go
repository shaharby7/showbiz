package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shaharby7/showbiz/services/api/internal/model"
)

type ProjectRepo struct {
	db *sql.DB
}

func NewProjectRepo(db *sql.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) CreateProject(ctx context.Context, project *model.Project) error {
	query := `INSERT INTO projects (id, name, organization_id, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())`
	_, err := r.db.ExecContext(ctx, query, project.ID, project.Name, project.OrganizationID, project.Description)
	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	return nil
}

func (r *ProjectRepo) GetProjectByID(ctx context.Context, id string) (*model.Project, error) {
	query := `SELECT id, name, organization_id, description, created_at, updated_at
		FROM projects WHERE id = ?`
	p := &model.Project{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.OrganizationID, &p.Description, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}
	return p, nil
}

func (r *ProjectRepo) GetProjectByOrgIDAndName(ctx context.Context, orgID, name string) (*model.Project, error) {
	query := `SELECT id, name, organization_id, description, created_at, updated_at
		FROM projects WHERE organization_id = ? AND name = ?`
	p := &model.Project{}
	err := r.db.QueryRowContext(ctx, query, orgID, name).Scan(
		&p.ID, &p.Name, &p.OrganizationID, &p.Description, &p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get project by org id and name: %w", err)
	}
	return p, nil
}

func (r *ProjectRepo) ListProjectsByOrgID(ctx context.Context, orgID string, cursor string, limit int) ([]*model.Project, string, error) {
	if limit <= 0 {
		limit = 20
	}

	var rows *sql.Rows
	var err error

	// Fetch one extra to determine if there are more results
	if cursor == "" {
		query := `SELECT id, name, organization_id, description, created_at, updated_at
			FROM projects WHERE organization_id = ? ORDER BY id ASC LIMIT ?`
		rows, err = r.db.QueryContext(ctx, query, orgID, limit+1)
	} else {
		query := `SELECT id, name, organization_id, description, created_at, updated_at
			FROM projects WHERE organization_id = ? AND id > ? ORDER BY id ASC LIMIT ?`
		rows, err = r.db.QueryContext(ctx, query, orgID, cursor, limit+1)
	}
	if err != nil {
		return nil, "", fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []*model.Project
	for rows.Next() {
		p := &model.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.OrganizationID, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, "", fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, p)
	}
	if err := rows.Err(); err != nil {
		return nil, "", fmt.Errorf("list projects rows: %w", err)
	}

	var nextCursor string
	if len(projects) > limit {
		nextCursor = projects[limit-1].ID
		projects = projects[:limit]
	}

	return projects, nextCursor, nil
}

func (r *ProjectRepo) UpdateProject(ctx context.Context, id, description string) error {
	query := `UPDATE projects SET description = ?, updated_at = NOW() WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, description, id)
	if err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update project rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("project not found")
	}
	return nil
}

func (r *ProjectRepo) DeleteProject(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Cascade: delete dependent rows then the project itself
	cascadeQueries := []string{
		`DELETE FROM policy_attachments WHERE project_id = ?`,
		`DELETE FROM resources WHERE project_id = ?`,
		`DELETE FROM connections WHERE project_id = ?`,
		`DELETE FROM projects WHERE id = ?`,
	}

	for _, q := range cascadeQueries {
		if _, err := tx.ExecContext(ctx, q, id); err != nil {
			return fmt.Errorf("delete project cascade: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (r *ProjectRepo) CountProjectsByOrgID(ctx context.Context, orgID string) (int, error) {
	query := `SELECT COUNT(*) FROM projects WHERE organization_id = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, orgID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count projects by org id: %w", err)
	}
	return count, nil
}
