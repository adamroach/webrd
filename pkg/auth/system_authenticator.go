package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/adamroach/webrd/pkg/config"
)

type SystemAuthenticator struct {
	passwordChecker PasswordChecker
	config          *config.Auth
	secret          Secret
}

func NewSystemAuthenticator(config *config.Auth) *SystemAuthenticator {
	secret, err := NewHmacSecret(config.HmacKey)
	if err != nil {
		panic(fmt.Errorf("cannot read hmac key: %v", err))
	}
	return &SystemAuthenticator{
		passwordChecker: NewPasswordChecker(),
		config:          config,
		secret:          secret,
	}
}

func (a *SystemAuthenticator) Authenticate(username, password string) (token string, err error) {
	if a.passwordChecker.CurrentUser() != username {
		return "", errors.New("invalid username")
	}
	if !a.passwordChecker.CheckPassword(username, password) {
		return "", errors.New("invalid password")
	}
	claims := NewClaims(username, true, time.Duration(a.config.TokenValidityHours)*time.Hour)
	token, err = claims.Token(a.secret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return token, nil
}

func (a *SystemAuthenticator) ValidateToken(token string) (username string, err error) {
	claims, err := NewClaimsFromToken(token, a.secret)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}
	if !claims.IsSystemUser {
		return "", errors.New("not a system user token")
	}
	if claims.IsExpired() {
		return "", errors.New("token is expired")
	}
	if claims.Subject != a.passwordChecker.CurrentUser() {
		return "", errors.New("invalid username")
	}
	return claims.Subject, nil
}
