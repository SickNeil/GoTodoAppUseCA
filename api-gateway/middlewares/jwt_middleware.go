// middlewares/jwt_middleware.go
package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供 Token"})
			return
		}
		// 此處可選擇驗證 Token 的格式，或將驗證留給 todo-app
		c.Next()
	}
}
