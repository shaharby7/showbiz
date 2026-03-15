package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shaharby7/showbiz/services/api/internal/model"
)

type MemberRepo struct {
	db *sql.DB
}

func NewMemberRepo(db *sql.DB) *MemberRepo {
	return &MemberRepo{db: db}
}

func (r *MemberRepo) ListMembers(ctx context.Context, orgID string) ([]*model.User, error) {
	query := `SELECT email, password_hash, organization_id, display_name, email_verified, active, created_at, updated_at
		FROM users WHERE organization_id = ?`
	rows, err := r.db.QueryContext(ctx, query, orgID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(&u.Email, &u.PasswordHash, &u.OrganizationID, &u.DisplayName, &u.EmailVerified, &u.Active, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan member: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list members rows: %w", err)
	}
	return users, nil
}

func (r *MemberRepo) AddMember(ctx context.Context, orgID, email string) error {
	query := `UPDATE users SET organization_id = ?, updated_at = NOW() WHERE email = ?`
	result, err := r.db.ExecContext(ctx, query, orgID, email)
	if err != nil {
		return fmt.Errorf("add member: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("add member rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *MemberRepo) RemoveMember(ctx context.Context, orgID, email string) error {
	query := `UPDATE users SET organization_id = '', updated_at = NOW() WHERE email = ? AND organization_id = ?`
	result, err := r.db.ExecContext(ctx, query, email, orgID)
	if err != nil {
		return fmt.Errorf("remove member: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("remove member rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("member not found in organization")
	}
	return nil
}
