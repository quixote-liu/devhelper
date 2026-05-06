package api

import (
	"strings"

	"devhelper/internal/utils"

	"github.com/gin-gonic/gin"
)

var jwtSecret string

func SetJWTSecret(s string) {
	jwtSecret = s
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			utils.Unauthorized(c, "authorization header required")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := utils.ParseToken(tokenStr, jwtSecret)
		if err != nil || claims.Type != "access" {
			utils.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("user_role")
		if role != "admin" {
			utils.Forbidden(c, "admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}

func CORSMiddleware(origins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origins)
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
