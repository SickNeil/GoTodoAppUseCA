// main.go
package main

import (
	"api-gateway/handlers"
	"api-gateway/middlewares"
	"api-gateway/services"
	"api-gateway/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 讀取配置
	config := utils.LoadConfig()

	// 初始化服務
	authService := services.NewAuthService(config.AuthServerURL)
	authHandler := handlers.NewAuthHandler(authService)

	todoHandler, err := handlers.NewTodoHandler(config.TodoAppURL)
	if err != nil {
		log.Fatal("初始化 TodoHandler 失敗：", err)
	}

	// 初始化 Gin
	router := gin.Default()

	// 登入路由
	router.POST("/login", authHandler.Login)

	// 需要驗證的路由
	authorized := router.Group("/todo")
	authorized.Use(middlewares.JWTAuthMiddleware())
	{
		authorized.Any("/*proxyPath", todoHandler.Proxy)
	}

	// 啟動服務
	router.Run(":5000")
}
