package main

import (
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/handler"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("警告: .env ファイルの読み込みに失敗しました。環境変数を直接使用します。")
	}

	if err := config.LoadConfig(); err != nil {
		log.Fatalf("main: 設定のロードに失敗しました: %v", err)
	}

	// 2. データベース接続を初期化 (設定からDSNを使用)
	if err := database.InitDB(config.AppConfig.DatabaseDSN); err != nil {
		log.Fatalf("main: データベースの初期化に失敗しました: %v", err)
	}

	defer func() {
		if err := database.DB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		} else {
			fmt.Println("Database connection closed.")
		}
	}()

	router := gin.Default()

	router.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

	api := router.Group("/api")
	{
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/login", handler.HandleLogin)
		}

		userRoutes := api.Group("/users")
		{
			userRoutes.POST("", handler.HandleCreateUser)
			userRoutes.GET("/:id", handler.HandleGetUserID)
			userRoutes.GET("", handler.HandleGetAllUsers)
			userRoutes.PUT("/:id", handler.HandleUpdateUser)
			userRoutes.DELETE("/:id", handler.HandleDeleteUser)
		}
	}

	port := config.AppConfig.ServerPort
	fmt.Println("Starting server on :8080")
	if err := router.Run(port); err != nil {
		log.Fatalf("main: Error staring Gin server: %v", err)
	}
}
