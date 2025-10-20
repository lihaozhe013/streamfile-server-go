package middleware

import (
	"net/http"
	"simple-server/src/backend/config"
	"simple-server/src/backend/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityMiddleware is a security middleware
func SecurityMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check path security
		if strings.Contains(c.Request.URL.Path, "..") {
			utils.SendError(c, http.StatusBadRequest, "Invalid path")
			c.Abort()
			return
		}

		// Add security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")

		c.Next()
	}
}

// CORSMiddleware is a CORS middleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
