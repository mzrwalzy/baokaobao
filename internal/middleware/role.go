package middleware

import (
	"baokaobao/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// RoleMiddleware creates a middleware that checks if the admin user's role
// is in the allowed roles list. Must be used after AdminAuth middleware.
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			response.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if roleStr == allowed {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "无权限访问")
		c.Abort()
	}
}
