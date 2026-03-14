package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/showbiz-io/showbiz/services/api/internal/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (email, password_hash, organization_id, display_name, email_verified, active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`
	_, err := r.db.ExecContext(ctx, query,
		user.Email, user.PasswordHash, user.OrganizationID, user.DisplayName, user.EmailVerified, user.Active,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT email, password_hash, organization_id, display_name, email_verified, active, created_at, updated_at
		FROM users WHERE email = ?`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.Email, &user.PasswordHash, &user.OrganizationID, &user.DisplayName,
		&user.EmailVerified, &user.Active, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, email string, displayName string) error {
	query := `UPDATE users SET display_name = ?, updated_at = NOW() WHERE email = ?`
	result, err := r.db.ExecContext(ctx, query, displayName, email)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("update user rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepo) SetUserActive(ctx context.Context, email string, active bool) error {
	query := `UPDATE users SET active = ?, updated_at = NOW() WHERE email = ?`
	result, err := r.db.ExecContext(ctx, query, active, email)
	if err != nil {
		return fmt.Errorf("set user active: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("set user active rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
