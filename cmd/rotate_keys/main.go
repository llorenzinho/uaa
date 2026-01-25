package main

import (
	"github.com/llorenzinho/goauth/internal/config"
	"github.com/llorenzinho/goauth/internal/database"
	"github.com/llorenzinho/goauth/internal/services"
)

func main() {
	cfg := config.NewAppConfig()
	pool := database.CreatePool(&cfg.DBConfig)
	defer pool.Close()
	s := services.NewJwksService(pool)
	s.Rotate()
}
