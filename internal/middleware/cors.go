package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"baokaobao/internal/config"
)

func isOriginAllowed(origin string, allowed []string) bool {
	if origin == "" {
		return false
	}
	for _, o := range allowed {
		if o == origin {
			return true
		}
	}
	return false
}

func CORS() gin.HandlerFunc {
	defaultOrigins := []string{"http://localhost:5173", "http://localhost:8080"}

	return func(c *gin.Context) {
		allowedOrigins := config.GlobalConfig.App.AllowedOrigins
		if len(allowedOrigins) == 0 {
			allowedOrigins = defaultOrigins
		}

		origin := c.Request.Header.Get("Origin")
		if isOriginAllowed(origin, allowedOrigins) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}