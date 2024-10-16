// frameworks/http/handlers.go
package handlers

import (
	"fmt"
	"go-todo-app/usecases"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type TodoHandler struct {
	UseCase *usecases.TodoUseCase
}

func NewTodoHandler(useCase *usecases.TodoUseCase) *TodoHandler {
	return &TodoHandler{UseCase: useCase}
}

func (h *TodoHandler) ShowTodos(c *gin.Context) {
	// 從 cookie 的 JWT Token 中取得使用者名稱及email，public key 位於 /key/public.key
	tokenString, err := c.Cookie("token")
	fmt.Println("handlers showtodos tokenString", tokenString)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	keyData, err := os.ReadFile("/key/public.key")
	if err != nil {
		fmt.Println("Error reading public key:", err)
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(keyData)
	})

	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.String(http.StatusUnauthorized, "invalid token")
		return
	}

	username := claims["username"].(string)
	email := claims["email"].(string)
	fmt.Println("username", username)
	fmt.Println("email", email)

	fmt.Println("app server show todos")
	fmt.Println("app server show todos tokenString", tokenString)
	h.UseCase.SetJWT(h.UseCase.Repo, tokenString)
	todos, err := h.UseCase.GetTodos()
	if err != nil {
		fmt.Println("app server show todos error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 渲染 HTML 模板，並傳入待辦事項和使用者名稱
	c.HTML(http.StatusOK, "todos.html", gin.H{
		"tasks":    todos,
		"Username": username,
		"Email":    email,
	})
}

func (h *TodoHandler) AddTodo(c *gin.Context) {
	task := c.PostForm("task")
	if task == "" {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	tokenString, err := c.Cookie("token")

	h.UseCase.SetJWT(h.UseCase.Repo, tokenString)
	err = h.UseCase.AddTodo(task)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id := c.Param("id")

	tokenString, err := c.Cookie("token")

	h.UseCase.SetJWT(h.UseCase.Repo, tokenString)
	err = h.UseCase.DeleteTodo(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}
