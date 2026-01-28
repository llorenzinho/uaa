package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/llorenzinho/goauth/api/internal/config"
	"github.com/llorenzinho/goauth/api/internal/log"
	"go.uber.org/zap"
)

func CreatePool(cfg *config.DBConfig) *pgxpool.Pool {
	l := log.Get()
	ctx, close := context.WithTimeout(context.Background(), time.Second*3)
	defer close()
	pool, err := pgxpool.New(ctx, cfg.ConnectionString)
	if err != nil {
		l.Fatal("Failed to create Database Pool", zap.Error(err))
	}

	err = pool.Ping(ctx)
	if err != nil {
		l.Fatal("Failed to Ping Database", zap.Error(err))
	}
	return pool
}
