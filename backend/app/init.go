package app

import (
	"os"

	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/router"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
)

// AppConfig defines configuration options for initializing the application.
// Includes toggles for middleware and required secrets.
type AppConfig struct {
	JWTSecret      string
	EnableCORS     bool
	EnableLogger   bool
	EnableRecovery bool
}

// InitApp initializes all application components and returns a fully configured Gin router.
//
// Components initialized:
//   - Repositories: DynamoDB implementations for songs, documents, search, and authentication
//   - Services: business logic layers wired with required dependencies
//   - Handlers: HTTP controllers connected to services
//   - Router: sets up routes and middleware with the configured handlers
//
// Parameters:
//   - db: a DynamoDB connection
//   - cfg: configuration struct for middleware and secrets
//
// Returns:
//   - a *gin.Engine instance ready to serve HTTP requests
func InitApp(db *dynamo.DB, cfg AppConfig) *gin.Engine {

	// Initialize repositories
	documentRepo := repository.NewDynamoDocumentRepository(db)
	songRepo := repository.NewDynamoSongRepository(db, documentRepo)
	searchRepo := repository.NewDynamoSearchRepository(db, documentRepo)
	authRepo := repository.NewAWSAuthRepository(os.Getenv("ENV"))

	// Initialize services
	idGen := &utils.UUIDGenerator{}
	timeProvider := &utils.UTCTimeProvider{}
	tokenGen := &utils.JWTTokenGenerator{Secret: []byte(cfg.JWTSecret)}

	songService := services.NewSongService(songRepo, documentRepo, idGen, timeProvider)
	documentService := services.NewDocumentService(documentRepo, songRepo, idGen, timeProvider)
	searchService := services.NewSearchService(searchRepo)
	authService := services.NewAuthService(authRepo, timeProvider, tokenGen)

	// Initialize handlers
	songHandler := handlers.NewSongHandler(songService)
	documentHandler := handlers.NewDocumentHandler(documentService)
	searchHandler := handlers.NewSearchHandler(searchService)
	authHandler := handlers.NewAuthHandler(authService)

	// Router
	return router.SetupRouter(songHandler, documentHandler, searchHandler, authHandler, router.RouterOptions{
		EnableCORS:     cfg.EnableCORS,
		EnableLogger:   cfg.EnableLogger,
		EnableRecovery: cfg.EnableRecovery,
	})
}
