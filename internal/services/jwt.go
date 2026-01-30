package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/llorenzinho/goauth/internal/database"
	"github.com/llorenzinho/goauth/pkg/client"
	"github.com/llorenzinho/goauth/pkg/log"
	"go.uber.org/zap"
)

type JwtService struct {
	l *zap.Logger
	q *database.Queries
	p *database.DBTX
	c context.Context
}

func NewJwtService(p *database.DBTX) *JwtService {
	return &JwtService{
		l: log.Get(),
		q: database.New(),
		c: context.TODO(),
		p: p,
	}
}

func (s *JwtService) MakeJWT(u *database.User) (string, error) {
	if u == nil {
		return "", errors.New("User cannot be null")
	}

	claims := client.JwtClaims{
		UserID: u.ID.String(),
		Email:  u.Email,
		Roles:  make([]string, 0),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "uaa", // TODO: make configurable
			ID:        fmt.Sprintf("user:%s", u.ID.String()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // TODO: make configurable
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   u.ID.String(),
		},
	}

	k, err := s.q.GetActiveJwk(s.c, *s.p)
	if err != nil {
		return "", nil
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t.Header["kid"] = k.Kid

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(k.PrivateKeyPem))
	if err != nil {
		return "", err
	}

	tokenString, err := t.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
