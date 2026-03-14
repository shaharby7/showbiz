package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/showbiz-io/showbiz/services/api/internal/middleware"
	"github.com/showbiz-io/showbiz/services/api/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input service.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if input.Email == "" || input.Password == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email and password are required")
		return
	}

	user, err := h.svc.Register(r.Context(), input)
	if err != nil {
		if err.Error() == "email already registered" {
			Error(w, http.StatusConflict, "EMAIL_EXISTS", "Email is already registered")
			return
		}
		slog.Error("register failed", "error", err)
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to register user")
		return
	}

	JSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if input.Email == "" || input.Password == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Email and password are required")
		return
	}

	accessToken, refreshToken, err := h.svc.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
			return
		}
		if err.Error() == "account is deactivated" {
			Error(w, http.StatusForbidden, "ACCOUNT_DEACTIVATED", "Account is deactivated")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to login")
		return
	}

	JSON(w, http.StatusOK, map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	if input.RefreshToken == "" {
		Error(w, http.StatusBadRequest, "VALIDATION_ERROR", "Refresh token is required")
		return
	}

	accessToken, refreshToken, err := h.svc.RefreshToken(r.Context(), input.RefreshToken)
	if err != nil {
		if err.Error() == "invalid refresh token" || err.Error() == "refresh token expired" {
			Error(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired refresh token")
			return
		}
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to refresh token")
		return
	}

	JSON(w, http.StatusOK, map[string]string{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	email := middleware.GetUserEmail(r.Context())
	if email == "" {
		Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	user, err := h.svc.GetCurrentUser(r.Context(), email)
	if err != nil {
		Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get user")
		return
	}

	JSON(w, http.StatusOK, user)
}
