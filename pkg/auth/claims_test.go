package auth_test

import (
	"testing"
	"time"

	"github.com/adamroach/webrd/mock"
	"github.com/adamroach/webrd/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClaims(t *testing.T) {
	username := "testuser"
	isSystemUser := true
	validity := time.Hour

	claims := auth.NewClaims(username, isSystemUser, validity)

	assert.Equal(t, username, claims.Subject)
	assert.Equal(t, isSystemUser, claims.IsSystemUser)
	assert.NotNil(t, claims.ExpiresAt)
	assert.WithinDuration(t, time.Now().Add(validity), claims.ExpiresAt.Time, time.Second)
}

func TestNewClaimsFromToken(t *testing.T) {
	secret := mock.NewSecret(t)
	secret.On("Get").Return([]byte("test-string"))
	username := "testuser"
	isSystemUser := true
	validity := time.Hour

	claims := auth.NewClaims(username, isSystemUser, validity)
	tokenString, err := claims.Token(secret)
	require.NoError(t, err)

	parsedClaims, err := auth.NewClaimsFromToken(tokenString, secret)
	require.NoError(t, err)

	assert.Equal(t, claims.Subject, parsedClaims.Subject)
	assert.Equal(t, claims.IsSystemUser, parsedClaims.IsSystemUser)
	assert.Equal(t, claims.ExpiresAt, parsedClaims.ExpiresAt)
}

func TestNewClaimsFromToken_InvalidToken(t *testing.T) {
	secret := mock.NewSecret(t)
	invalidToken := "invalid.token.string"

	_, err := auth.NewClaimsFromToken(invalidToken, secret)
	assert.Error(t, err)
}

func TestNewClaimsFromToken_InvalidAudience(t *testing.T) {
	secret := mock.NewSecret(t)
	secret.On("Get").Return([]byte("test-string"))
	username := "testuser"
	isSystemUser := true
	validity := time.Hour

	claims := auth.NewClaims(username, isSystemUser, validity)
	claims.Audience = []string{"invalid-audience"}
	tokenString, err := claims.Token(secret)
	require.NoError(t, err)

	_, err = auth.NewClaimsFromToken(tokenString, secret)
	assert.Error(t, err)
	assert.Equal(t, "invalid audience", err.Error())
}

func TestToken(t *testing.T) {
	secret := mock.NewSecret(t)
	secret.On("Get").Return([]byte("test-string"))
	username := "testuser"
	isSystemUser := true
	validity := time.Hour

	claims := auth.NewClaims(username, isSystemUser, validity)
	tokenString, err := claims.Token(secret)
	require.NoError(t, err)

	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (any, error) {
		return secret.Get(), nil
	})
	require.NoError(t, err)
	assert.True(t, token.Valid)
}

func TestIsExpired(t *testing.T) {
	claims := auth.NewClaims("testuser", false, -time.Hour)
	assert.True(t, claims.IsExpired())

	claims = auth.NewClaims("testuser", false, time.Hour)
	assert.False(t, claims.IsExpired())
}
