package main

import (
	"backend/internal/database"
	"backend/internal/handler"
	"fmt"
	"log"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // MySQL 驱动包
)

func main() {
	dsn := "root:20170407@tcp(127.0.0.1:3306)/go_practice?charset=utf8mb4&parseTime=True&loc=Local"
	if err := database.InitDB(dsn); err != nil {
		log.Fatalf("main: Failed to initialize database: %v", err)
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

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", handler.HandleCreateUser)
		userRoutes.GET("/:id", handler.HandleGetUserID)
		userRoutes.GET("", handler.HandleGetAllUsers)
		userRoutes.PUT("/:id", handler.HandleUpdateUser)
		userRoutes.DELETE("/:id", handler.HandleDeleteUser)
	}

	fmt.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("main: Error staring Gin server: %v", err)
	}
}
