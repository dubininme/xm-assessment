package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/dubininme/xm-assessment/internal/infra/auth"
	"github.com/dubininme/xm-assessment/pkg/gen/oapi"
)

type ctxKey string

const UserIDKey ctxKey = "user_id"

type AuthMiddleware struct {
	jwtService *auth.JWTService
}

func NewAuthMiddleware(jwtService *auth.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		tokenString, err := m.jwtService.ExtractToken(authHeader)
		if err != nil {
			writeUnauthorized(w, "missing or invalid authorization header")
			return
		}

		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			writeUnauthorized(w, "missing or invalid authorization header")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func writeUnauthorized(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)

	resp := oapi.Error{
		Code:    oapi.ErrorCodeUnauthorized,
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode JSON response", "error", err)
	}
}
