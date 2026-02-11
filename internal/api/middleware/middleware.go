package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS middleware para permitir peticiones cross-origin
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Logger middleware personalizado
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate response time
		latency := time.Since(start)

		// Get status
		statusCode := c.Writer.Status()

		// Build query
		if raw != "" {
			path = path + "?" + raw
		}

		// Log format
		log.Printf("[API] %s | %3d | %13v | %15s | %-7s %s",
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			c.ClientIP(),
			c.Request.Method,
			path,
		)
	}
}

// Recovery middleware personalizado
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				log.Printf("[RECOVERY] panic recovered: %v", err)

				// Return error response
				c.JSON(500, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "PANIC",
						"message": "Internal server error",
						"details": fmt.Sprintf("%v", err),
					},
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RateLimiter middleware simple (opcional)
func RateLimiter() gin.HandlerFunc {
	// Implementación simple - en producción usar redis
	return func(c *gin.Context) {
		// Por ahora solo lo dejamos pasar
		c.Next()
	}
}
