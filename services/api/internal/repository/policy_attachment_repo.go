package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shaharby7/showbiz/services/api/internal/model"
)

type PolicyAttachmentRepo struct {
	db *sql.DB
}

func NewPolicyAttachmentRepo(db *sql.DB) *PolicyAttachmentRepo {
	return &PolicyAttachmentRepo{db: db}
}

func (r *PolicyAttachmentRepo) AttachPolicy(ctx context.Context, attachment *model.PolicyAttachment) error {
	query := `INSERT INTO policy_attachments (id, project_id, user_email, policy_id, created_at)
		VALUES (?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query, attachment.ID, attachment.ProjectID, attachment.UserEmail, attachment.PolicyID)
	if err != nil {
		return fmt.Errorf("attach policy: %w", err)
	}
	return nil
}

func (r *PolicyAttachmentRepo) DetachPolicy(ctx context.Context, projectID, userEmail, policyID string) error {
	query := `DELETE FROM policy_attachments WHERE project_id = ? AND user_email = ? AND policy_id = ?`
	result, err := r.db.ExecContext(ctx, query, projectID, userEmail, policyID)
	if err != nil {
		return fmt.Errorf("detach policy: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("detach policy rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("attachment not found")
	}
	return nil
}

func (r *PolicyAttachmentRepo) ListAttachmentsByProject(ctx context.Context, projectID string) ([]*model.PolicyAttachment, error) {
	query := `SELECT id, project_id, user_email, policy_id, created_at
		FROM policy_attachments WHERE project_id = ? ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("list attachments by project: %w", err)
	}
	defer rows.Close()

	var attachments []*model.PolicyAttachment
	for rows.Next() {
		a := &model.PolicyAttachment{}
		if err := rows.Scan(&a.ID, &a.ProjectID, &a.UserEmail, &a.PolicyID, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan attachment: %w", err)
		}
		attachments = append(attachments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list attachments rows: %w", err)
	}
	if attachments == nil {
		attachments = []*model.PolicyAttachment{}
	}
	return attachments, nil
}

func (r *PolicyAttachmentRepo) ListAttachmentsByUser(ctx context.Context, userEmail, projectID string) ([]*model.PolicyAttachment, error) {
	query := `SELECT id, project_id, user_email, policy_id, created_at
		FROM policy_attachments WHERE user_email = ? AND project_id = ? ORDER BY created_at ASC`
	rows, err := r.db.QueryContext(ctx, query, userEmail, projectID)
	if err != nil {
		return nil, fmt.Errorf("list attachments by user: %w", err)
	}
	defer rows.Close()

	var attachments []*model.PolicyAttachment
	for rows.Next() {
		a := &model.PolicyAttachment{}
		if err := rows.Scan(&a.ID, &a.ProjectID, &a.UserEmail, &a.PolicyID, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan attachment: %w", err)
		}
		attachments = append(attachments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list attachments by user rows: %w", err)
	}
	if attachments == nil {
		attachments = []*model.PolicyAttachment{}
	}
	return attachments, nil
}

func (r *PolicyAttachmentRepo) GetUserPermissions(ctx context.Context, userEmail, projectID string) ([]string, error) {
	query := `SELECT DISTINCT CONCAT(pp.entity, ':', pp.action)
		FROM policy_attachments pa
		JOIN policies p ON pa.policy_id = p.id
		JOIN policy_permissions pp ON pp.policy_id = p.id
		WHERE pa.user_email = ? AND pa.project_id = ?`
	rows, err := r.db.QueryContext(ctx, query, userEmail, projectID)
	if err != nil {
		return nil, fmt.Errorf("get user permissions: %w", err)
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var perm string
		if err := rows.Scan(&perm); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		perms = append(perms, perm)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("user permissions rows: %w", err)
	}
	if perms == nil {
		perms = []string{}
	}
	return perms, nil
}
