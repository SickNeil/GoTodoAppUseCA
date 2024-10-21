// interfaces/handlers/auth_handlers.go
package handlers

import (
	"auth/entities"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UseCase entities.IUserUseCase
}

func NewAuthHandler(useCase entities.IUserUseCase) *AuthHandler {
	return &AuthHandler{
		UseCase: useCase,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供正確的登入資訊"})
		return
	}

	token, err := h.UseCase.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var registerRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請提供正確的註冊資訊"})
		return
	}

	user := &entities.User{
		Username: registerRequest.Username,
		Password: registerRequest.Password,
		Email:    registerRequest.Email,
	}

	if err := h.UseCase.Register(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "註冊成功"})
}

func (h *AuthHandler) IsTokenValid(c *gin.Context) {
	fmt.Println("Checking token")
	// 從標頭中獲取 Authorization 標頭的內容
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		fmt.Println("No token provided")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供授權 Token"})
		c.Abort()
		return
	}

	// 檢查標頭是否包含 Bearer 前綴
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		fmt.Println("Invalid token format")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的授權 Token"})
		c.Abort()
		return
	}

	// 提取 Bearer 之後的 Token 部分
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		fmt.Println("Empty token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "授權 Token 是空的"})
		c.Abort()
		return
	}

	// 使用 UseCase 驗證 Token 是否有效
	isValid, err := h.UseCase.IsTokenValid(token)
	if err != nil {
		fmt.Println("Error validating token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "驗證 Token 時出錯"})
		c.Abort()
		return
	}

	if !isValid {
		fmt.Println("Token is invalid")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token 無效"})
		c.Abort()
		return
	}

	// 驗證成功後繼續執行
	c.Next()
}
