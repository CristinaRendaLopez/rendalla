package router

import (
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RouterOptions struct {
	EnableCORS     bool
	EnableLogger   bool
	EnableRecovery bool
}

func SetupRouter(
	songHandler *handlers.SongHandler,
	documentHandler *handlers.DocumentHandler,
	searchHandler *handlers.SearchHandler,
	authHandler *handlers.AuthHandler,
	opts RouterOptions,
) *gin.Engine {

	// Set up Gin router
	r := gin.New()

	// CORS
	if opts.EnableCORS {
		r.Use(cors.Default())
	}

	// Middleware for structured logging and error handling
	if opts.EnableLogger {
		r.Use(gin.Logger())
	}
	if opts.EnableRecovery {
		r.Use(gin.Recovery())
	}

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

	return r
}
