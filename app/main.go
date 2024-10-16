// main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	handlers "go-todo-app/frameworks/http"
	todoInterfaces "go-todo-app/interfaces/todo"
	"go-todo-app/usecases"
	"go-todo-app/utils"
	"html/template"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 初始化 MongoDB
	mongoClientOptions := options.Client().ApplyURI("mongodb://mongo-database:27017")
	mongoClient, err := mongo.Connect(context.TODO(), mongoClientOptions)
	if err != nil {
		panic(err)
	}
	todoCollection := mongoClient.Database("tododb").Collection("todos")

	// 初始化 PostgreSQL
	dbHost := os.Getenv("POSTGRES_DB_HOST")
	dbUser := os.Getenv("POSTGRES_DB_USER")
	dbPassword := os.Getenv("POSTGRES_DB_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName)

	postgresDB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	// 初始化 Repository
	todoRepo := todoInterfaces.NewMongoDBTodoRepository(todoCollection)

	// 初始化 UseCase
	todoUseCase := usecases.NewTodoUseCase(todoRepo)

	// 初始化 Handler
	todoHandler := handlers.NewTodoHandler(todoUseCase)

	// 初始化 Gin
	router := gin.Default()

	// 設置自定義函數到模板引擎
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": utils.FormatAsDate,
	})

	// 需要 JWT 認證的路由群組
	authorized := router.Group("/")
	authorized.Use(handlers.JWTAuth())
	{
		// 登入與註冊路由
		authorized.GET("/get-all", todoHandler.ShowTodos)
		authorized.POST("/add", todoHandler.AddTodo)
		authorized.DELETE("/delete/:id", todoHandler.DeleteTodo)
	}

	// 啟動服務器
	router.Run(":3000")
}
