package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("JWTAuth Middleware")
		// 從 Cookie 中獲取 Token
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		// 創建一個新的 HTTP POST 請求
		req, err := http.NewRequest("POST", os.Getenv("AUTH_SERVER_URL")+"/validate-token", nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "請求錯誤"})
			return
		}

		// 將 Token 放入 Authorization 標頭
		req.Header.Set("Authorization", "Bearer "+tokenString)
		fmt.Println("header", req.Header)

		// 發送 HTTP 請求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token 無效"})
			return
		}

		c.Next()
	}
}

func AttachJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Cookie 中獲取 JWT
		tokenString, err := c.Cookie("token")
		if err != nil {
			// 如果 Cookie 中沒有 JWT，返回 401 錯誤
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供 JWT Token"})
			c.Abort()
			return
		}

		// 將 JWT Token 附加到 Authorization 標頭中
		c.Request.Header.Set("Authorization", "Bearer "+tokenString)

		// 打印調試信息，確認 JWT 已附加到標頭
		fmt.Println("JWT Token 附加到 Authorization 標頭: Bearer " + tokenString)

		// 繼續執行後續的請求處理
		c.Next()
	}
}
