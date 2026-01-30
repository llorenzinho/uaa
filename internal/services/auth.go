package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/llorenzinho/goauth/internal/database"
	"github.com/llorenzinho/goauth/internal/rest/dtos"
	"github.com/llorenzinho/goauth/pkg/log"
	"go.uber.org/zap"
)

type AuthService struct {
	c context.Context
	p *pgxpool.Pool
	q database.Queries
	l *zap.Logger
}

func NewAuthService(p *pgxpool.Pool) *AuthService {
	return &AuthService{
		c: context.TODO(),
		p: p,
		q: *database.New(),
		l: log.Get(),
	}
}

func generateAuthorizationCode(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// Usiamo URLEncoding per evitare caratteri come '+' o '/' che rompono gli URL
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (as *AuthService) CreateAuthCode(u *database.User, p *dtos.AuthorizationCodeQueryParams) (string, error) {
	c, err := generateAuthorizationCode(32)
	if err != nil {
		return "", nil
	}

	codeParams := &database.CreateAuthorizationCodeParams{
		Code:          c,
		UserID:        u.ID,
		ClientID:      p.CLientId,
		RedirectUri:   p.RedirectUri,
		Scope:         &p.Scope,
		CodeChallenge: nil,
		ExpiresAt:     pgtype.Timestamptz{Valid: true, Time: time.Now().Add(5 * time.Minute)},
	}

	authCode, err := as.q.CreateAuthorizationCode(as.c, as.p, *codeParams)
	if err != nil {
		return "", nil
	}
	return authCode.Code, nil

}
