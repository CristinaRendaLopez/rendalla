package app

import (
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/router"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
)

// AppConfig agrupa dependencias que podrías querer inyectar
type AppConfig struct {
	JWTSecret      string
	EnableCORS     bool
	EnableLogger   bool
	EnableRecovery bool
}

func InitApp(db *dynamo.DB, cfg AppConfig) *gin.Engine {
	// Initialize repositories
	documentRepo := repository.NewDynamoDocumentRepository(db)
	songRepo := repository.NewDynamoSongRepository(db, documentRepo)
	searchRepo := repository.NewDynamoSearchRepository(db, documentRepo)
	authRepo := repository.NewAWSAuthRepository()

	// Initialize services
	idGen := &utils.UUIDGenerator{}
	timeProvider := &utils.UTCTimeProvider{}
	clock := &utils.RealClock{}
	tokenGen := &utils.JWTTokenGenerator{Secret: []byte(cfg.JWTSecret)}

	songService := services.NewSongService(songRepo, documentRepo, idGen, timeProvider)
	documentService := services.NewDocumentService(documentRepo, idGen, timeProvider)
	searchService := services.NewSearchService(searchRepo)
	authService := services.NewAuthService(authRepo, clock, tokenGen)

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
