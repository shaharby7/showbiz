package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/showbiz-io/showbiz/services/api/internal/model"
)

type PolicyRepo struct {
	db *sql.DB
}

func NewPolicyRepo(db *sql.DB) *PolicyRepo {
	return &PolicyRepo{db: db}
}

func (r *PolicyRepo) CreatePolicy(ctx context.Context, policy *model.Policy) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO policies (id, name, scope, organization_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())`
	var orgID *string
	if policy.OrganizationID != "" {
		orgID = &policy.OrganizationID
	}
	_, err = tx.ExecContext(ctx, query, policy.ID, policy.Name, policy.Scope, orgID)
	if err != nil {
		return fmt.Errorf("insert policy: %w", err)
	}

	for _, perm := range policy.Permissions {
		parts := strings.SplitN(perm, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid permission format: %s", perm)
		}
		permQuery := `INSERT INTO policy_permissions (id, policy_id, entity, action) VALUES (?, ?, ?, ?)`
		permID := policy.ID + "-" + parts[0] + "-" + parts[1]
		_, err = tx.ExecContext(ctx, permQuery, permID, policy.ID, parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("insert permission %s: %w", perm, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (r *PolicyRepo) GetPolicyByID(ctx context.Context, id string) (*model.Policy, error) {
	query := `SELECT id, name, scope, organization_id, created_at, updated_at
		FROM policies WHERE id = ?`
	policy := &model.Policy{}
	var orgID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&policy.ID, &policy.Name, &policy.Scope, &orgID, &policy.CreatedAt, &policy.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get policy by id: %w", err)
	}
	if orgID.Valid {
		policy.OrganizationID = orgID.String
	}

	perms, err := r.loadPermissions(ctx, id)
	if err != nil {
		return nil, err
	}
	policy.Permissions = perms
	return policy, nil
}

func (r *PolicyRepo) ListGlobalPolicies(ctx context.Context) ([]*model.Policy, error) {
	query := `SELECT id, name, scope, organization_id, created_at, updated_at
		FROM policies WHERE scope = 'global' ORDER BY name ASC`
	return r.listPolicies(ctx, query)
}

func (r *PolicyRepo) ListOrgPolicies(ctx context.Context, orgID string) ([]*model.Policy, error) {
	query := `SELECT id, name, scope, organization_id, created_at, updated_at
		FROM policies WHERE scope = 'organization' AND organization_id = ? ORDER BY name ASC`
	return r.listPolicies(ctx, query, orgID)
}

func (r *PolicyRepo) DeletePolicy(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM policies WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete policy: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete policy rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("policy not found")
	}
	return nil
}

func (r *PolicyRepo) loadPermissions(ctx context.Context, policyID string) ([]string, error) {
	query := `SELECT entity, action FROM policy_permissions WHERE policy_id = ?`
	rows, err := r.db.QueryContext(ctx, query, policyID)
	if err != nil {
		return nil, fmt.Errorf("load permissions: %w", err)
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var entity, action string
		if err := rows.Scan(&entity, &action); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		perms = append(perms, entity+":"+action)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("permissions rows: %w", err)
	}
	if perms == nil {
		perms = []string{}
	}
	return perms, nil
}

func (r *PolicyRepo) listPolicies(ctx context.Context, query string, args ...interface{}) ([]*model.Policy, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list policies: %w", err)
	}
	defer rows.Close()

	var policies []*model.Policy
	for rows.Next() {
		p := &model.Policy{}
		var orgID sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Scope, &orgID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan policy: %w", err)
		}
		if orgID.Valid {
			p.OrganizationID = orgID.String
		}
		perms, err := r.loadPermissions(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		p.Permissions = perms
		policies = append(policies, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list policies rows: %w", err)
	}
	if policies == nil {
		policies = []*model.Policy{}
	}
	return policies, nil
}
