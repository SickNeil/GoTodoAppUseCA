// main.go
package main

import (
	"context"
	"database/sql"
	"fmt"
	handlers "go-todo-app/frameworks/http"
	loginInterfaces "go-todo-app/interfaces/login"
	todoInterfaces "go-todo-app/interfaces/todo"
	"go-todo-app/usecases"
	"go-todo-app/utils"
	"html/template"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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
	loginRepo := loginInterfaces.NewPostgresLoginRepository(postgresDB)

	// 初始化 UseCase
	todoUseCase := usecases.NewTodoUseCase(todoRepo)
	loginUseCase := usecases.NewLoginUseCase(loginRepo)

	// 初始化 Handler
	todoHandler := handlers.NewTodoHandler(todoUseCase)
	loginHandler := handlers.NewLoginHandler(loginUseCase)

	// 初始化 Gin
	router := gin.Default()

	// 設定 Session 中間件
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// 設置自定義函數到模板引擎
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": utils.FormatAsDate,
	})

	router.LoadHTMLGlob("/root/templates/*")

	// 登入與註冊路由
	router.GET("/login", loginHandler.ShowLoginPage)
	router.POST("/login", loginHandler.PerformLogin)
	router.GET("/register", loginHandler.ShowRegisterPage)
	router.POST("/register", loginHandler.PerformRegister)
	router.GET("/logout", loginHandler.Logout)

	// 需要登入的路由群組
	authorized := router.Group("/")
	authorized.Use(handlers.AuthRequired())
	{
		authorized.GET("/", todoHandler.ShowTodos)
		authorized.POST("/", todoHandler.AddTodo)
		authorized.POST("/delete/:id", todoHandler.DeleteTodo)
	}

	// 啟動服務器
	router.Run(":3000")
}
