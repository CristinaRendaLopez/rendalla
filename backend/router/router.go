package router

import (
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// RouterOptions allows enabling or disabling middleware features when setting up the router.
type RouterOptions struct {
	EnableCORS     bool
	EnableLogger   bool
	EnableRecovery bool
}

// SetupRouter configures and returns a new Gin router instance.
// It registers all public and protected routes, applying middleware as needed.
//
// Handlers:
//   - songHandler: handles song-related endpoints
//   - documentHandler: handles document-related endpoints
//   - searchHandler: handles search functionality for songs and documents
//   - authHandler: handles authentication endpoints
//
// RouterOptions:
//   - EnableCORS: enables CORS middleware if true
//   - EnableLogger: enables Gin's logging middleware if true
//   - EnableRecovery: enables panic recovery middleware if true
func SetupRouter(
	songHandler *handlers.SongHandler,
	documentHandler *handlers.DocumentHandler,
	searchHandler *handlers.SearchHandler,
	authHandler *handlers.AuthHandler,
	opts RouterOptions,
) *gin.Engine {

	r := gin.New()

	if opts.EnableCORS {
		r.Use(cors.Default())
	}
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
		public.GET("/songs/:song_id", songHandler.GetSongByIDHandler)

		public.GET("/songs/:song_id/documents", documentHandler.GetAllDocumentsBySongIDHandler)
		public.GET("/songs/:song_id/documents/:doc_id", documentHandler.GetDocumentByIDHandler)

		public.GET("/songs/search", searchHandler.ListSongsHandler)
		public.GET("/documents/search", searchHandler.ListDocumentsHandler)

		public.POST("/auth/login", authHandler.LoginHandler)
	}

	// Protected routes (authentication required)
	auth := r.Group("/")
	auth.Use(middleware.JWTAuthMiddleware())
	{
		auth.POST("/songs", songHandler.CreateSongHandler)
		auth.PUT("/songs/:song_id", songHandler.UpdateSongHandler)
		auth.DELETE("/songs/:song_id", songHandler.DeleteSongWithDocumentsHandler)

		auth.POST("/songs/:song_id/documents", documentHandler.CreateDocumentHandler)
		auth.PUT("/songs/:song_id/documents/:doc_id", documentHandler.UpdateDocumentHandler)
		auth.DELETE("/songs/:song_id/documents/:doc_id", documentHandler.DeleteDocumentHandler)

		auth.GET("/auth/me", authHandler.MeHandler)
	}

	return r
}
