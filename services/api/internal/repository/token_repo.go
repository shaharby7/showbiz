package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/showbiz-io/showbiz/services/api/internal/model"
)

type TokenRepo struct {
	db *sql.DB
}

func NewTokenRepo(db *sql.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

func (r *TokenRepo) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES (?, ?, ?, ?, NOW())`
	_, err := r.db.ExecContext(ctx, query,
		token.ID, token.UserID, token.TokenHash, token.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}
	return nil
}

func (r *TokenRepo) GetRefreshTokenByHash(ctx context.Context, hash string) (*model.RefreshToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, created_at
		FROM refresh_tokens WHERE token_hash = ?`
	token := &model.RefreshToken{}
	err := r.db.QueryRowContext(ctx, query, hash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get refresh token by hash: %w", err)
	}
	return token, nil
}

func (r *TokenRepo) DeleteRefreshTokensByUser(ctx context.Context, email string) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, email)
	if err != nil {
		return fmt.Errorf("delete refresh tokens by user: %w", err)
	}
	return nil
}

func (r *TokenRepo) DeleteExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("delete expired tokens: %w", err)
	}
	return nil
}
