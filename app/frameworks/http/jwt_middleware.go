package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtSecret = []byte(os.Getenv("JWT_SECRET"))

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Cookie 中獲取 Token
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// 解析並驗證 Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrInvalidKeyType
			}
			return JwtSecret, nil
		})

		fmt.Println("Token", token)
		if err != nil || !token.Valid {
			fmt.Println("Token is invalid")
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// 將使用者資訊保存到上下文中
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("Token claims is invalid")
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Set("username", claims["username"])
		c.Next()
	}
}

func TokenAuthMiddleware(jwtSecret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Cookie 中獲取 Token
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token not found in cookies"})
			return
		}

		// 解析和驗證 Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			return
		}

		// 將 Token 中的資訊存入上下文，供後續處理使用
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("username", claims["username"])
		}

		c.Next()
	}
}
