package shared

import (
	"context"

	"github.com/llorenzinho/goauth/shared/queries"
)

type JwksService struct {
	db queries.DBTX
}

func NewJwksService(db *queries.DBTX) JwksService {
	return JwksService{db: *db}
}

func (js *JwksService) ListJWKSEntries() ([]queries.JwkKey, error) {
	return queries.New().ListKeys(context.TODO(), js.db)
}

// func (js *JwksService) GetJWKSKeyById(id string) (*queries.JwkKey, error) {

// 	k, err := queries.New().GetJwksKey(context.TODO(), js.db)
// }
