package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("JWTAuth Middleware")
		// 從 header Authorization 中獲取 Token，去除 Bearer 字串
		tokenString := c.GetHeader("Authorization")
		tokenString = tokenString[7:]
		fmt.Println("tokenString", tokenString)

		// 公鑰位於 /key/public.key ，使用公鑰來驗證 Token
		keyData, err := os.ReadFile("/key/public.key")
		if err != nil {
			fmt.Println("Error reading public key:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error reading public key"})
			c.Abort()
			return
		}

		// 解析 Token 並驗證簽名
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 確認 token 使用的演算法是 RS256
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// 解析 RSA 公鑰
			return jwt.ParseRSAPublicKeyFromPEM(keyData)
		})

		if err != nil || !token.Valid {
			fmt.Println("Token is invalid" + err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func TokenAuthMiddleware() gin.HandlerFunc {
	fmt.Println("TokenAuthMiddleware Middleware")
	return func(c *gin.Context) {
		// 從 Header 中獲取 Token
		tokenString := c.GetHeader("Authorization")

		// 公鑰位於 /key/public.key ，使用公鑰來驗證 Token
		keyData, err := os.ReadFile("/key/public.key")
		if err != nil {
			fmt.Println("Error reading public key:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error reading public key"})
			c.Abort()
			return
		}

		// 解析 Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwt.ParseRSAPublicKeyFromPEM(keyData)
		})

		if err != nil || !token.Valid {
			fmt.Println("Token is invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			c.Abort()
			return
		}

		// 將使用者資訊保存到上下文中
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Println("Token claims is invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token claims is invalid"})
			c.Abort()
			return
		}

		c.Set("username", claims["username"])
		c.Next()
	}
}
