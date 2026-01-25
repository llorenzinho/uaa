package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/llorenzinho/goauth/internal/database"
	"github.com/llorenzinho/goauth/internal/log"
	"go.uber.org/zap"
)

type JwksService struct {
	l *zap.Logger
	p *pgxpool.Pool
}

func NewJwksService(pool *pgxpool.Pool) JwksService {
	return JwksService{
		l: log.Get(),
		p: pool,
	}
}

func (j *JwksService) generatePEMKeyPairs() *rsa.PrivateKey {
	private, err := rsa.GenerateKey(rand.Reader, 2024)
	j.l.Info("Generating KEY PAIRS")
	if err != nil {
		j.l.Fatal("Unable to generate private key", zap.Error(err))
	}
	return private
}

func (j *JwksService) pEMEncode(k *rsa.PrivateKey) ([]byte, []byte) {
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(k)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	l := log.Get()
	if err != nil {
		l.Fatal("Error while PEM ENCODING")
	}
	pkb := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	pukb := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	privateKeyPemBytes := pem.EncodeToMemory(pkb)
	publicKeyPemBytes := pem.EncodeToMemory(pukb)
	return privateKeyPemBytes, publicKeyPemBytes
}

func (j *JwksService) Rotate() {
	key := j.generatePEMKeyPairs()
	privPem, pubPem := j.pEMEncode(key)
	q := database.New()

	kid := base64.RawURLEncoding.EncodeToString(pubPem)
	ctx, cls := context.WithTimeout(context.Background(), time.Second*3)
	defer cls()
	q.CreateNewRs256Key(ctx, j.p, database.CreateNewRs256KeyParams{
		Kid:           kid,
		PrivateKeyPem: string(privPem),
		PublicKeyPem:  string(pubPem),
		ExpiresAt:     pgtype.Timestamptz{Valid: true, Time: time.Now().Add(time.Hour * 24 * 7)}, // Make Configurable
	})
	delKeys, err := q.DeleteExpiredKey(ctx, j.p)
	if err != nil {
		j.l.Fatal("Failed ")
	}

	delKeysStrings := make([]string, len(delKeys))
	for i, el := range delKeys {
		delKeysStrings[i] = el.Kid
	}
	j.l.Info("Successfuly deleted expired keys", zap.Strings("keys", delKeysStrings))
	q.ActivateKey(ctx, j.p, kid)
}
