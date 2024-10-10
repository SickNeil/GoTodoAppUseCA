// main.go
package main

import (
	"auth/handlers"
	"auth/usecases"
	"database/sql"
	"fmt"
	"log"
	"os"

	"auth/interfaces"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

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

	userRepo := interfaces.NewPostgresUserRepository(postgresDB)
	jwtAuth := usecases.NewJWTAuth()
	userUseCase := usecases.NewUserUseCase(*userRepo, jwtAuth)

	// 初始化 Handler
	authHandler := handlers.NewAuthHandler(userUseCase)

	// 初始化 Gin
	router := gin.Default()

	// 設定路由
	router.POST("/login", authHandler.Login)
	router.POST("/register", authHandler.Register)
	router.POST("/validate-token", authHandler.IsTokenValid)

	// 啟動服務器
	router.Run(":4000")
}
