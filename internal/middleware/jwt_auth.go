package middleware

import (
	"strings"

	"baokaobao/internal/config"
	"baokaobao/internal/pkg/jwt"
	"baokaobao/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

var jwtInstance *jwt.JWT

func InitJWT() {
	jwtInstance = jwt.NewJWT(config.GlobalConfig.JWT.Secret, config.GlobalConfig.JWT.ExpireHours)
}

func GetJWT() *jwt.JWT {
	return jwtInstance
}

func MiniProgramAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtInstance.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		if claims.UserType != "mini" {
			response.Forbidden(c, "access denied: mini program only")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("openid", claims.OpenID)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtInstance.ParseToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "invalid token")
			c.Abort()
			return
		}

		if claims.UserType != "admin" {
			response.Forbidden(c, "access denied: admin only")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.OpenID)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}
