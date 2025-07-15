package main

import (
	"fmt"
	"log"
	"os"

	"news-fetcher/db"
	"news-fetcher/handlers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("📰 Starting News Fetcher API server...")

	// Load environment variables from .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, relying on environment variables")
	}

	if os.Getenv("NEWS_API_KEY") == "" {
		log.Fatal("❌ NEWS_API_KEY is not set in environment")
	}

	// Initialize database
	db.InitDB()

	// Create Gin router and define routes
	r := gin.Default()
	r.POST("/news/fetch", handlers.FetchNewsHandler)

	// Run server on port 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
