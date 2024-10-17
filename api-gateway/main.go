// main.go
package main

import (
	handlers "go-todo-app/frameworks/http"
	loginInterfaces "go-todo-app/interfaces/login"
	todoInterfaces "go-todo-app/interfaces/todo"
	"go-todo-app/usecases"
	"go-todo-app/utils"
	"html/template"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	// 初始化 Repository
	todoRepo := todoInterfaces.NewJwtTodoRepository()
	loginRepo := loginInterfaces.NewJwtLoginRepository()

	// 初始化 UseCase
	todoUseCase := usecases.NewTodoUseCase(todoRepo)
	loginUseCase := usecases.NewLoginUseCase(loginRepo)

	// 初始化 Handler
	todoHandler := handlers.NewTodoHandler(todoUseCase)
	userProcessHandler := handlers.NewLoginProcessHandler(loginUseCase)

	// 初始化 Gin
	router := gin.Default()

	// 設置自定義函數到模板引擎
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": utils.FormatAsDate,
	})

	router.LoadHTMLGlob("/root/templates/*")

	// 登入與註冊路由
	router.GET("/login", userProcessHandler.ShowLoginPage)
	router.POST("/login", userProcessHandler.PerformLogin)
	router.GET("/register", userProcessHandler.ShowRegisterPage)
	router.POST("/register", userProcessHandler.PerformRegister)
	router.GET("/logout", userProcessHandler.Logout)

	// 需要 JWT 認證的路由群組
	// authorized := router.Group("/")
	// authorized.Use(handlers.JWTAuth())
	// {
	router.GET("/", todoHandler.ShowTodos)
	router.POST("/", todoHandler.AddTodo)
	router.POST("/delete/:id", todoHandler.DeleteTodo)
	// }

	// 啟動服務器
	router.Run(":5000")
}
