package auth

import (
	"errors"
	"os"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	IsSystemUser bool `json:"sys"`
	jwt.RegisteredClaims
}

func NewClaims(username string, isSystemUser bool, validity time.Duration) *Claims {
	iss, err := os.Hostname()
	if err != nil {
		iss = "webrd"
	}
	return &Claims{
		IsSystemUser: isSystemUser,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    iss,
			Subject:   username,
			Audience:  []string{iss},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(validity)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
}

func NewClaimsFromToken(tokenString string, secret Secret) (*Claims, error) {
	aud, err := os.Hostname()
	if err != nil {
		aud = "webrd"
	}
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret.Get(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		if !slices.Contains(claims.Audience, aud) {
			return nil, errors.New("invalid audience")
		}
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (c *Claims) Token(secret Secret) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secret.Get())
}

func (c *Claims) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(c.ExpiresAt.Time)
}
