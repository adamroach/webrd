package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/adamroach/webrd/pkg/config"
)

type StaticAuthenticator struct {
	userPass map[string]string
	secret   Secret
	config   *config.Auth
}

func NewStaticAuthenticator(config *config.Auth) *StaticAuthenticator {
	secret, err := NewHmacSecret(config.HmacKey)
	if err != nil {
		panic(err)
	}
	userPass := make(map[string]string)
	for _, user := range config.Users {
		userPass[user.Username] = user.Password
	}
	return &StaticAuthenticator{
		userPass: userPass,
		secret:   secret,
		config:   config,
	}
}

func (a *StaticAuthenticator) Authenticate(username, password string) (token string, err error) {
	if pass, ok := a.userPass[username]; ok {
		if pass == password {
			claims := NewClaims(username, false, time.Duration(a.config.TokenValidityHours)*time.Hour)
			token, err = claims.Token(a.secret)
			if err != nil {
				return "", fmt.Errorf("failed to sign token: %w", err)
			}
			return token, nil
		}
	}
	return "", errors.New("invalid username or password")
}

func (a *StaticAuthenticator) ValidateToken(token string) (username string, err error) {
	claims, err := NewClaimsFromToken(token, a.secret)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}
	if claims.IsSystemUser {
		return "", errors.New("unexpected system user token")
	}
	if claims.IsExpired() {
		return "", errors.New("token is expired")
	}
	return claims.Subject, nil
}
