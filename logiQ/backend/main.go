package main

import (
	"backend/database"
	"backend/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	log.Println("Starting Quiz backend service...")
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, fallback to environment variables")
	}
	database.ConnectMongoDB()
	database.ConnectRedis()
	r := routes.SetupRoutes()

	log.Println("Server is running on port 3000...")
	if err := r.Run("0.0.0.0:3000"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
