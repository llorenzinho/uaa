package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/llorenzinho/goauth/pkg/log"
	"go.uber.org/zap"
)

type UaaClient struct {
	c            *http.Client
	s            KeyStore
	timeout      time.Duration
	ctx          context.Context
	baseEndpoint string
	l            *zap.Logger
}

func NewUaaClient(baseEndpoint string, opts ...UaaOption) *UaaClient {
	cl := &UaaClient{
		c:            &http.Client{},
		s:            &InMemoryKeystore{},
		timeout:      time.Second * 5, // Default timeout
		ctx:          context.Background(),
		baseEndpoint: baseEndpoint,
		l:            log.Get(),
	}

	for _, o := range opts {
		o(cl)
	}

	cl.c.Timeout = cl.timeout
	return cl
}

func (uc *UaaClient) listKeys() ([]JwkKey, error) {
	if err := uc.ctx.Err(); err != nil {
		uc.l.Error("Context error", zap.Error(err))
		return nil, err
	}
	res, err := uc.c.Get(fmt.Sprintf("%s/api/v1/jwk", uc.baseEndpoint))
	if err != nil {
		uc.l.Error("Failed fetching data", zap.Error(err))
		return nil, err
	}
	// Parse result
	var jwks JwksResponse
	if err := json.NewDecoder(res.Body).Decode(&jwks); err != nil {
		uc.l.Error("Error while decoding jswon", zap.Error(err))
		return nil, err
	}

	if err := validator.New().Struct(jwks); err != nil {
		uc.l.Error("Validation Error", zap.Error(err))
		return nil, err
	}
	return jwks.Keys, nil

}

func (uc *UaaClient) RefreshJwks() {
	uc.l.Info("Jwks refresh triggered")
	keys, err := uc.listKeys()
	if err != nil {
		return
	}
	uc.s.Clean()
	for _, k := range keys {
		uc.s.Set(&k)
	}
}

func (uc *UaaClient) ValidateToken(token string) (*JwtClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(
		token,
		&JwtClaims{},
		func(t *jwt.Token) (any, error) {

			kid, ok := t.Header["kid"].(string)
			if !ok {
				return nil, errors.New("kid not present")
			}

			if !uc.s.Exist(kid) {
				// try fetching data
				uc.RefreshJwks()
				if !uc.s.Exist(kid) {
					return nil, errors.New("provided kid does not exist")
				}
			}

			jwk, err := uc.s.Get(kid)
			if err != nil {
				return nil, err
			}

			if t.Method.Alg() != jwk.Algorithm {
				return nil, errors.New("invalid algorithm")
			}

			return jwk.PublicKeyPem, nil // this automatically validate the token sign
		},
	)

	if err != nil {
		return nil, err
	}

	// 4. Token valido?
	if !jwtToken.Valid {
		return nil, errors.New("token non valido")
	}

	claims, ok := jwtToken.Claims.(*JwtClaims)
	if !ok {
		return nil, errors.New("claims non valide")
	}

	// 5. Issuer
	if claims.Issuer != "auth.myapp.com" {
		return nil, errors.New("issuer non valido")
	}

	// // 6. Audience
	// if !claims.Audience.Contains("api.myapp.com") {
	// 	return nil, errors.New("audience non valida")
	// }

	// 7. Claim applicativi
	if claims.UserID == "" {
		return nil, errors.New("user id mancante")
	}

	return claims, nil
}
