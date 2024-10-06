// main.go
package main

import (
	"context"
	handlers "go-todo-app/frameworks/http"
	interfaces "go-todo-app/interfaces/todo"
	"go-todo-app/usecases"
	"go-todo-app/utils"
	"text/template"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 初始化 MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://mongo-database:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err) // 如果連接失敗，程式將停止並顯示錯誤
	}

	todoCollection := client.Database("tododb").Collection("todos")

	// 使用 MongoDB 實作
	mongoRepo := interfaces.NewMongoDBTodoRepository(todoCollection)

	// 初始化 Gin
	router := gin.Default()

	// 設置自定義函數到模板引擎
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": utils.FormatAsDate, // 使用 utils 中的 FormatAsDate 函數
	})

	router.LoadHTMLGlob("/root/templates/*")

	// 初始化 Use Case
	todoUseCase := usecases.NewTodoUseCase(mongoRepo)

	// 初始化 HTTP Handler
	todoHandler := handlers.NewTodoHandler(todoUseCase)

	// 設定路由
	router.GET("/", todoHandler.ShowTodos)
	router.POST("/", todoHandler.AddTodo)
	router.POST("/delete/:id", todoHandler.DeleteTodo)

	// 啟動服務器
	router.Run(":3000")
}
