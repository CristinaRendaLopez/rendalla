package main

import (
	"fmt"
	"os"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Could not load .env file, using default values")
	}

	// Initialize configuration and database
	bootstrap.LoadConfig()
	bootstrap.InitDB()

	// Set up Gin router
	r := gin.Default()

	// Enable CORS
	r.Use(cors.Default())

	// Middleware for structured logging and error handling
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Rendalla backend is running successfully",
		})
	})

	// Public routes (no authentication required)
	public := r.Group("/")
	{
		public.GET("/songs", handlers.GetAllSongsHandler)
		public.GET("/songs/:id", handlers.GetSongByIDHandler)
		public.GET("/songs/:id/documents", handlers.GetAllDocumentsBySongIDHandler)
		public.GET("/documents/:id", handlers.GetDocumentByIDHandler)
		public.GET("/songs/search", handlers.SearchSongsByTitleHandler)
		public.GET("/documents/search", handlers.SearchDocumentsByTitleHandler)
		public.GET("/documents/filter", handlers.FilterDocumentsByInstrumentHandler)
		public.POST("/auth/login", handlers.LoginHandler)
	}

	// Protected routes (authentication required)
	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.POST("/songs", handlers.CreateSongHandler)
		auth.PUT("/songs/:id", handlers.UpdateSongHandler)
		auth.DELETE("/songs/:id", handlers.DeleteSongWithDocumentsHandler)

		auth.POST("/songs/:id/documents", handlers.CreateDocumentHandler)
		auth.PUT("/documents/:id", handlers.UpdateDocumentHandler)
		auth.DELETE("/documents/:id", handlers.DeleteDocumentHandler)

		auth.GET("/auth/me", handlers.MeHandler)
	}

	// Get and validate port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logrus.Infof("Rendalla backend is running on port %s", port)
	r.Run(fmt.Sprintf(":%s", port))
}
