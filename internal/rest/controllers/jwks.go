package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/llorenzinho/goauth/internal/services"
	"github.com/llorenzinho/goauth/pkg/client"
	"github.com/llorenzinho/goauth/pkg/log"
	"go.uber.org/zap"
)

type JwkController struct {
	s *services.JwksService
	l *zap.Logger
}

func NewJwkController(s *services.JwksService) *JwkController {
	return &JwkController{s: s, l: log.Get()}
}

func (jc *JwkController) HandleListJwk(c *gin.Context) {
	jwks, err := jc.s.ListValidJwkKeys()

	if err != nil {
		jc.l.Error("Failed to get jwk list", zap.Error(err))
		c.JSON(http.StatusInternalServerError, nil)
	}

	c.JSON(http.StatusOK, client.JwksResponse{Keys: jwks})
}
