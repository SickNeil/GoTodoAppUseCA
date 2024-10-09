// interfaces/handlers/auth_handlers.go
package handlers

import (
	"auth/entities"
	"auth/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UseCase usecases.UserUseCase
}

func NewAuthHandler(useCase usecases.UserUseCase) *AuthHandler {
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
