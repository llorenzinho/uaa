package main

import (
	"context"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/llorenzinho/goauth/internal/config"
	"github.com/llorenzinho/goauth/internal/database"
	"github.com/llorenzinho/goauth/internal/rest/controllers"
	"github.com/llorenzinho/goauth/internal/rest/middlewares"
	"github.com/llorenzinho/goauth/internal/services"
	"github.com/llorenzinho/goauth/pkg/log"
	"go.uber.org/zap"
)

func gracefulShutdown(server *http.Server, done chan struct{}) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	<-ctx.Done()
	if err := server.Shutdown(context.Background()); err != nil {
		log.Get().Error("Server forced to shutdown", zap.Error(err))
	}
	log.Get().Info("Server Graceful Shutdown Complete")
	done <- struct{}{}
}

func main() {
	done := make(chan struct{})
	config := config.NewAppConfig()
	pool := database.CreatePool(&config.DBConfig)
	defer pool.Close()

	// Services
	userService := services.NewUserService(pool)
	jwkService := services.NewJwksService(pool)

	// Controllers
	userController := controllers.NewUserController(userService)
	jwkController := controllers.NewJwkController(&jwkService)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.LogMiddleware)

	api := r.Group("api")
	v1 := api.Group("v1")
	userApi := v1.Group("users")
	jwkApi := v1.Group("jwk")

	userApi.GET(":id", userController.GetUserByID)
	userApi.POST("", userController.CreateUser)

	jwkApi.GET("", jwkController.HandleListJwk)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.ServerConfig.Host, config.ServerConfig.Port),
		Handler: r,
	}
	l := log.Get()
	l.Info("Starting server", zap.String("addr", server.Addr))

	go gracefulShutdown(server, done)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal("Error while starting server", zap.Error(err))
		}
	}()
	<-done
}
