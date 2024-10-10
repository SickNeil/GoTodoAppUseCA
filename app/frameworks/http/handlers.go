// frameworks/http/handlers.go
package handlers

import (
	"fmt"
	"go-todo-app/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	UseCase *usecases.TodoUseCase
}

type AddTodoRequest struct {
	Task string `json:"task" binding:"required"`
}

func NewTodoHandler(useCase *usecases.TodoUseCase) *TodoHandler {
	return &TodoHandler{UseCase: useCase}
}

func (h *TodoHandler) ShowTodos(c *gin.Context) {
	fmt.Println("app server show todos")
	todos, err := h.UseCase.GetTodos()
	if err != nil {
		fmt.Println("app server show todos error", err)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("app server todos", todos)
	c.JSON(http.StatusOK, gin.H{
		"tasks": todos,
	})
}

func (h *TodoHandler) AddTodo(c *gin.Context) {
	var req AddTodoRequest

	// 解析 JSON 請求
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 檢查 task 是否為空
	if req.Task == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task is required"})
		return
	}
	err := h.UseCase.AddTodo(req.Task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "task added"})
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	err := h.UseCase.DeleteTodo(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}
