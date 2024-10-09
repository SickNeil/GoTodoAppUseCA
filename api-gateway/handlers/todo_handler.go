// handlers/todo_handler.go
package handlers

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	TodoServiceURL *url.URL
}

func NewTodoHandler(todoServiceURL string) (*TodoHandler, error) {
	parsedURL, err := url.Parse(todoServiceURL)
	if err != nil {
		return nil, err
	}
	return &TodoHandler{TodoServiceURL: parsedURL}, nil
}

func (h *TodoHandler) Proxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(h.TodoServiceURL)
	proxy.ServeHTTP(c.Writer, c.Request)
}
