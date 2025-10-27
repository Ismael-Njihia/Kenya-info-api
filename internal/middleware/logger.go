package middleware

import (
	"time"

	"github.com/Ismael-Njihia/Kenya-info-api/internal/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger logs all incoming HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		logger.Log.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("ip", clientIP),
		)
	}
}
