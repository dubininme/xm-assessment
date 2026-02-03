package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/dubininme/xm-assessment/pkg/gen/oapi"
)

type TokenGenerator interface {
	GenerateToken(userID string, expiresIn time.Duration) (string, error)
}

type AuthHandler struct {
	tokenGenerator TokenGenerator
}

func NewAuthHandler(tokenGenerator TokenGenerator) *AuthHandler {
	return &AuthHandler{
		tokenGenerator: tokenGenerator,
	}
}

type GenerateTokenRequest struct {
	UserID   string `json:"user_id"`
	Password string `json:"password"`
}

type GenerateTokenResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	var req GenerateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, oapi.ErrorCodeBadRequest, "invalid request body")
		return
	}

	if len(req.UserID) == 0 {
		writeErr(w, http.StatusBadRequest, oapi.ErrorCodeBadRequest, "user_id is required")
		return
	}

	if len(req.Password) == 0 {
		writeErr(w, http.StatusBadRequest, oapi.ErrorCodeBadRequest, "password is required")
		return
	}

	// Simple authentication check for demo purposes
	// In production, this would check against a user database with hashed passwords
	if req.Password != "demo-password-123" {
		writeErr(w, http.StatusUnauthorized, oapi.ErrorCodeUnauthorized, "invalid credentials")
		return
	}

	token, err := h.tokenGenerator.GenerateToken(req.UserID, time.Hour)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, oapi.ErrorCodeInternalError, "failed to generate token")
		return
	}

	resp := GenerateTokenResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
