// frameworks/http/login_handlers.go
package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-todo-app/usecases"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserProcessHandler struct {
	UseCase *usecases.UserProcessUseCase
}

func NewUserProcessHandler(useCase *usecases.UserProcessUseCase) *UserProcessHandler {
	return &UserProcessHandler{UseCase: useCase}
}

// ShowLoginPage 顯示登入頁面
func (h *UserProcessHandler) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// PerformLogin 處理登入請求
func (h *UserProcessHandler) PerformLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 準備向認證伺服器發送的請求資料
	loginData := map[string]string{
		"username": username,
		"password": password,
	}
	jsonData, err := json.Marshal(loginData)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"Error": "系統錯誤，請稍後再試",
		})
		return
	}

	// 讀取認證伺服器的回應
	body, err := SendToAuthServer(jsonData, "http://auth-server:4000/login")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"Error": "登入失敗，請檢查帳號密碼",
		})
		return
	}

	// 認證成功，取得 Token
	var successResponse map[string]string
	json.Unmarshal(body, &successResponse)
	token := successResponse["token"]

	// 將 Token 保存到 Cookie 中
	c.SetCookie("token", token, 3600, "/", "", false, true)

	// 重定向到主頁
	c.Redirect(http.StatusSeeOther, "/")
}

func SendToAuthServer(param []byte, url string) (response []byte, err error) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(param))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("認證失敗: " + string(body))
	}

	return body, nil
}

func (h *UserProcessHandler) Logout(c *gin.Context) {
	// 清除 Cookie 中的 Token
	c.SetCookie("token", "", -1, "/", "", false, true)

	// 重定向到登入頁面
	c.Redirect(http.StatusSeeOther, "/login")
}

// ShowRegisterPage 顯示註冊頁面
func (h *UserProcessHandler) ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

// PerformRegister 處理註冊請求
func (h *UserProcessHandler) PerformRegister(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	registerData := map[string]string{
		"username": username,
		"password": password,
		"email":    email,
	}
	jsonData, err := json.Marshal(registerData)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"Error": "系統錯誤，請稍後再試",
		})
		return
	}

	// 向認證伺服器發送註冊請求
	_, err = SendToAuthServer(jsonData, "http://auth-server:4000/register")
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"Error": "註冊失敗，請稍後再試 " + err.Error(),
		})
		return
	}

	// 註冊成功，重定向到登入頁面
	c.Redirect(http.StatusSeeOther, "/login")
}
