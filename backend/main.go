package main

import (
	"fmt"
	"os"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/router"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
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

	router := router.SetupRouter(songHandler, documentHandler, searchHandler, authHandler)

	// Get and validate port
	port := bootstrap.AppPort
	if port == "" {
		port = "8080"
	}

	logrus.Infof("Rendalla backend is running on port %s", port)
	router.Run(fmt.Sprintf(":%s", port))
}
