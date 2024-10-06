// frameworks/http/handlers.go
package handlers

import (
	"go-todo-app/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	UseCase *usecases.TodoUseCase
}

func NewTodoHandler(useCase *usecases.TodoUseCase) *TodoHandler {
	return &TodoHandler{UseCase: useCase}
}

func (h *TodoHandler) ShowTodos(c *gin.Context) {
	todos, err := h.UseCase.GetTodos()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// 使用 c.HTML 渲染模板，並將待辦事項列表傳遞到模板中
	c.HTML(http.StatusOK, "todos.html", gin.H{
		"tasks": todos,
	})
}

func (h *TodoHandler) AddTodo(c *gin.Context) {
	task := c.PostForm("task")
	if task == "" {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	err := h.UseCase.AddTodo(task)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
	id := c.Param("id")
	err := h.UseCase.DeleteTodo(id)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}
