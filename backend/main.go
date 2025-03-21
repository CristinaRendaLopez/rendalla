package main

import (
	"fmt"
	"os"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/middleware"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
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

	// Initialize repositories
	songRepo := repository.NewDynamoSongRepository(repository.NewDynamoDocumentRepository())
	documentRepo := repository.NewDynamoDocumentRepository()
	searchRepo := repository.NewDynamoSearchRepository(bootstrap.DB)
	authRepo := repository.NewAWSAuthRepository()

	// Initialize services
	idGen := &utils.UUIDGenerator{}
	timeProvider := &utils.UTCTimeProvider{}
	clock := &utils.RealClock{}
	tokenGen := &utils.JWTTokenGenerator{Secret: []byte(os.Getenv("JWT_SECRET"))}

	songService := services.NewSongService(songRepo, documentRepo, idGen, timeProvider)
	documentService := services.NewDocumentService(documentRepo, idGen, timeProvider)
	searchService := services.NewSearchService(searchRepo)
	authService := services.NewAuthService(authRepo, clock, tokenGen)

	// Initialize handlers
	songHandler := handlers.NewSongHandler(songService)
	documentHandler := handlers.NewDocumentHandler(documentService)
	searchHandler := handlers.NewSearchHandler(searchService)
	authHandler := handlers.NewAuthHandler(authService)

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
		public.GET("/songs", songHandler.GetAllSongsHandler)
		public.GET("/songs/:id", songHandler.GetSongByIDHandler)
		public.GET("/songs/:id/documents", documentHandler.GetAllDocumentsBySongIDHandler)
		public.GET("/documents/:id", documentHandler.GetDocumentByIDHandler)
		public.GET("/songs/search", searchHandler.SearchSongsByTitleHandler)
		public.GET("/documents/search", searchHandler.SearchDocumentsByTitleHandler)
		public.GET("/documents/filter", searchHandler.FilterDocumentsByInstrumentHandler)
		public.POST("/auth/login", authHandler.LoginHandler)
	}

	// Protected routes (authentication required)
	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.POST("/songs", songHandler.CreateSongHandler)
		auth.PUT("/songs/:id", songHandler.UpdateSongHandler)
		auth.DELETE("/songs/:id", songHandler.DeleteSongWithDocumentsHandler)

		auth.POST("/songs/:id/documents", documentHandler.CreateDocumentHandler)
		auth.PUT("/documents/:id", documentHandler.UpdateDocumentHandler)
		auth.DELETE("/documents/:id", documentHandler.DeleteDocumentHandler)

		auth.GET("/auth/me", authHandler.MeHandler)
	}

	// Get and validate port
	port := bootstrap.AppPort
	if port == "" {
		port = "8080"
	}

	logrus.Infof("Rendalla backend is running on port %s", port)
	r.Run(fmt.Sprintf(":%s", port))
}
