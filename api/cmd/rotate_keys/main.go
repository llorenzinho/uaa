package main

import (
	"github.com/llorenzinho/goauth/api/internal/config"
	"github.com/llorenzinho/goauth/api/internal/database"
	"github.com/llorenzinho/goauth/api/internal/services"
)

func main() {
	cfg := config.NewAppConfig()
	pool := database.CreatePool(&cfg.DBConfig)
	defer pool.Close()
	s := services.NewJwksService(pool)
	s.Rotate()
}
