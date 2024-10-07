// frameworks/http/login_handlers.go
package handlers

import (
	"go-todo-app/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	UseCase *usecases.LoginUseCase
}

func NewLoginHandler(useCase *usecases.LoginUseCase) *LoginHandler {
	return &LoginHandler{UseCase: useCase}
}

// ShowLoginPage 顯示登入頁面
func (h *LoginHandler) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// PerformLogin 處理登入請求
func (h *LoginHandler) PerformLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	valid, err := h.UseCase.Login(username, password)
	if err != nil || !valid {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"Error": "使用者名稱或密碼錯誤",
		})
		return
	}

	// 登入成功，設定 Session 或 JWT（根據需求實作）
	c.Redirect(http.StatusSeeOther, "/")
}

// ShowRegisterPage 顯示註冊頁面
func (h *LoginHandler) ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

// PerformRegister 處理註冊請求
func (h *LoginHandler) PerformRegister(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	err := h.UseCase.RegisterUser(username, password, email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"Error": "註冊失敗，請重試",
		})
		return
	}

	// 註冊成功，重定向到登入頁面
	c.Redirect(http.StatusSeeOther, "/login")
}
