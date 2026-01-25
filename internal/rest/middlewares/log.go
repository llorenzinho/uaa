package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/llorenzinho/goauth/internal/log"
	"go.uber.org/zap"
)

var l zap.Logger

func init() {
	l = *log.Get()
}

func LogMiddleware(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery
	method := c.Request.Method

	c.Next()

	elapsed := time.Since(start)
	status := c.Writer.Status()
	l.Info(
		fmt.Sprintf("[%d] [%s] %s?%s", status, method, path, query),
		zap.String("method", method),
		zap.Uint16("status", uint16(status)),
		zap.String("client", c.ClientIP()),
		zap.Duration("elapsed", elapsed),
	)
}
