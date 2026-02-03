// +build unit

package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret-key"

func setupJWTService(t *testing.T) *JWTService {
	return NewJWTService(testSecret)
}

func TestExtractToken_Success(t *testing.T) {
	service := setupJWTService(t)

	token, err := service.ExtractToken("Bearer abc123token")

	require.NoError(t, err)
	assert.Equal(t, "abc123token", token)
}

func TestExtractToken_EmptyHeader(t *testing.T) {
	service := setupJWTService(t)

	token, err := service.ExtractToken("")

	assert.ErrorIs(t, err, ErrMissingAuthHeader)
	assert.Empty(t, token)
}

func TestExtractToken_InvalidFormat(t *testing.T) {
	service := setupJWTService(t)

	token, err := service.ExtractToken("InvalidFormat")

	assert.ErrorIs(t, err, ErrInvalidAuthFormat)
	assert.Empty(t, token)
}

func TestGenerateAndValidateToken_Success(t *testing.T) {
	service := setupJWTService(t)
	userID := "user-123"

	tokenString, err := service.GenerateToken(userID, 1*time.Hour)
	require.NoError(t, err)

	claims, err := service.ValidateToken(tokenString)
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	service := setupJWTService(t)

	tokenString, _ := service.GenerateToken("user-123", -1*time.Hour)

	claims, err := service.ValidateToken(tokenString)

	assert.ErrorIs(t, err, ErrTokenExpired)
	assert.Nil(t, claims)
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	wrongService := NewJWTService("wrong-secret")
	tokenString, _ := wrongService.GenerateToken("user-123", 1*time.Hour)

	service := setupJWTService(t)
	claims, err := service.ValidateToken(tokenString)

	assert.ErrorIs(t, err, ErrInvalidSignature)
	assert.Nil(t, claims)
}

func TestValidateToken_MalformedToken(t *testing.T) {
	service := setupJWTService(t)

	claims, err := service.ValidateToken("not-a-jwt-token")

	assert.ErrorIs(t, err, ErrInvalidToken)
	assert.Nil(t, claims)
}
