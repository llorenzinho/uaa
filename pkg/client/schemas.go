package client

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	UserID string `json:"uid"`
	Email  string `json:"email"`
	// Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JwkKey struct {
	Kid          string `json:"kid,omitempty" validate:"required"`
	PublicKeyPem string `json:"publicKeyPem" validate:"required"`
	Algorithm    string `json:"algorithm" validate:"required"`
}

type JwksResponse struct {
	Keys []JwkKey `validate:"required,min=1,dive"`
}
