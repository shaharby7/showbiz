package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/showbiz-io/showbiz/services/api/internal/auth"
	"github.com/showbiz-io/showbiz/services/api/internal/model"
	"github.com/showbiz-io/showbiz/services/api/internal/repository"
)

type AuthService struct {
	userRepo  *repository.UserRepo
	tokenRepo *repository.TokenRepo
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepo, tokenRepo *repository.TokenRepo, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		jwtSecret: jwtSecret,
	}
}

type RegisterInput struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	DisplayName    string `json:"displayName"`
	OrganizationID string `json:"organizationId"`
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*model.User, error) {
	existing, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("check existing user: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Email:          input.Email,
		PasswordHash:   string(hash),
		OrganizationID: input.OrganizationID,
		DisplayName:    input.DisplayName,
		EmailVerified:  true, // auto-verify for now
		Active:         true,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Re-fetch to get server-set timestamps
	created, err := s.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("fetch created user: %w", err)
	}
	return created, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return "", "", fmt.Errorf("invalid credentials")
	}
	if !user.Active {
		return "", "", fmt.Errorf("account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	accessToken, err = auth.GenerateAccessToken(email, s.jwtSecret)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	plainRefresh, hashedRefresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	rt := &model.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    email,
		TokenHash: hashedRefresh,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.tokenRepo.CreateRefreshToken(ctx, rt); err != nil {
		return "", "", fmt.Errorf("store refresh token: %w", err)
	}

	return accessToken, plainRefresh, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (newAccessToken, newRefreshToken string, err error) {
	hash := auth.HashToken(refreshToken)

	stored, err := s.tokenRepo.GetRefreshTokenByHash(ctx, hash)
	if err != nil {
		return "", "", fmt.Errorf("get refresh token: %w", err)
	}
	if stored == nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}
	if stored.ExpiresAt.Before(time.Now()) {
		return "", "", fmt.Errorf("refresh token expired")
	}

	// Delete the old token (single-use rotation)
	if err := s.tokenRepo.DeleteRefreshTokensByUser(ctx, stored.UserID); err != nil {
		return "", "", fmt.Errorf("delete old tokens: %w", err)
	}

	newAccessToken, err = auth.GenerateAccessToken(stored.UserID, s.jwtSecret)
	if err != nil {
		return "", "", fmt.Errorf("generate access token: %w", err)
	}

	plainRefresh, hashedRefresh, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}

	rt := &model.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    stored.UserID,
		TokenHash: hashedRefresh,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := s.tokenRepo.CreateRefreshToken(ctx, rt); err != nil {
		return "", "", fmt.Errorf("store refresh token: %w", err)
	}

	return newAccessToken, plainRefresh, nil
}

func (s *AuthService) GetCurrentUser(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}
