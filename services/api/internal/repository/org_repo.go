package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/showbiz-io/showbiz/services/api/internal/model"
)

type OrgRepo struct {
	db *sql.DB
}

func NewOrgRepo(db *sql.DB) *OrgRepo {
	return &OrgRepo{db: db}
}

func (r *OrgRepo) CreateOrganization(ctx context.Context, org *model.Organization) error {
	query := `INSERT INTO organizations (id, name, display_name, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())`
	_, err := r.db.ExecContext(ctx, query, org.ID, org.Name, org.DisplayName, org.Active)
	if err != nil {
		return fmt.Errorf("create organization: %w", err)
	}
	return nil
}

func (r *OrgRepo) GetOrganizationByID(ctx context.Context, id string) (*model.Organization, error) {
	query := `SELECT id, name, display_name, active, created_at, updated_at
		FROM organizations WHERE id = ?`
	org := &model.Organization{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&org.ID, &org.Name, &org.DisplayName, &org.Active, &org.CreatedAt, &org.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get organization by id: %w", err)
	}
	return org, nil
}

func (r *OrgRepo) UpdateOrganization(ctx context.Context, id, displayName string) error {
	query := `UPDATE organizations SET display_name = ?, updated_at = NOW() WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, displayName, id)
	if err != nil {
		return fmt.Errorf("update organization: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update organization rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("organization not found")
	}
	return nil
}

func (r *OrgRepo) SetOrganizationActive(ctx context.Context, id string, active bool) error {
	query := `UPDATE organizations SET active = ?, updated_at = NOW() WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, active, id)
	if err != nil {
		return fmt.Errorf("set organization active: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("set organization active rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("organization not found")
	}
	return nil
}

func (r *OrgRepo) ListOrganizations(ctx context.Context, cursor string, limit int) ([]*model.Organization, string, error) {
	if limit <= 0 {
		limit = 20
	}

	var rows *sql.Rows
	var err error

	// Fetch one extra to determine if there are more results
	if cursor == "" {
		query := `SELECT id, name, display_name, active, created_at, updated_at
			FROM organizations ORDER BY id ASC LIMIT ?`
		rows, err = r.db.QueryContext(ctx, query, limit+1)
	} else {
		query := `SELECT id, name, display_name, active, created_at, updated_at
			FROM organizations WHERE id > ? ORDER BY id ASC LIMIT ?`
		rows, err = r.db.QueryContext(ctx, query, cursor, limit+1)
	}
	if err != nil {
		return nil, "", fmt.Errorf("list organizations: %w", err)
	}
	defer rows.Close()

	var orgs []*model.Organization
	for rows.Next() {
		org := &model.Organization{}
		if err := rows.Scan(&org.ID, &org.Name, &org.DisplayName, &org.Active, &org.CreatedAt, &org.UpdatedAt); err != nil {
			return nil, "", fmt.Errorf("scan organization: %w", err)
		}
		orgs = append(orgs, org)
	}
	if err := rows.Err(); err != nil {
		return nil, "", fmt.Errorf("list organizations rows: %w", err)
	}

	var nextCursor string
	if len(orgs) > limit {
		nextCursor = orgs[limit-1].ID
		orgs = orgs[:limit]
	}

	return orgs, nextCursor, nil
}

// DeleteProjectsByOrgID deletes all projects belonging to an organization.
func (r *OrgRepo) DeleteProjectsByOrgID(ctx context.Context, orgID string) error {
	query := `DELETE FROM projects WHERE organization_id = ?`
	_, err := r.db.ExecContext(ctx, query, orgID)
	if err != nil {
		return fmt.Errorf("delete projects by org id: %w", err)
	}
	return nil
}
