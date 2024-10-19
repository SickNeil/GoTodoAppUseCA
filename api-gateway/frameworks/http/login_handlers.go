// frameworks/http/login_handlers.go
package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go-todo-app/usecases"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginProcessHandler struct {
	UseCase *usecases.UserProcessUseCase
}

func NewLoginProcessHandler(useCase *usecases.UserProcessUseCase) *LoginProcessHandler {
	return &LoginProcessHandler{UseCase: useCase}
}

// ShowLoginPage 顯示登入頁面
func (h *LoginProcessHandler) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// PerformLogin 處理登入請求
func (h *LoginProcessHandler) PerformLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	result, err := h.UseCase.Login(username, password)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"Error": "系統錯誤，請稍後再試",
		})
		return
	}

	if result.StatusCode != http.StatusOK {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"Error": "登入失敗，請檢查帳號密碼",
		})
		return
	}

	// 從回應中讀取 Token
	body, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"Error": "系統錯誤，請稍後再試",
		})
		return
	}

	var token map[string]string
	if err := json.Unmarshal(body, &token); err != nil {
		fmt.Println("Error unmarshalling response body: ", err)
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"Error": "系統錯誤，請稍後再試",
		})
		return
	}

	fmt.Println("Token: ", token["token"])

	// 將 Token 保存到 Cookie 中
	c.SetCookie("token", token["token"], 3600, "/", "", false, true)

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

func (h *LoginProcessHandler) Logout(c *gin.Context) {
	// 清除 Cookie 中的 Token
	c.SetCookie("token", "", -1, "/", "", false, true)

	// 重定向到登入頁面
	c.Redirect(http.StatusSeeOther, "/login")
}

// ShowRegisterPage 顯示註冊頁面
func (h *LoginProcessHandler) ShowRegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

// PerformRegister 處理註冊請求
func (h *LoginProcessHandler) PerformRegister(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	email := c.PostForm("email")

	resp, err := h.UseCase.Repo.Register(username, password, email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"Error": "系統錯誤，請稍後再試",
		})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "register.html", gin.H{
				"Error": "系統錯誤，請稍後再試",
			})
			return
		}
	}

	// 註冊成功，重定向到登入頁面
	c.Redirect(http.StatusSeeOther, "/login")
}
