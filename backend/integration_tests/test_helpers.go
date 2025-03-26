package integration_tests

import (
	"os"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/router"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	TestRouter *gin.Engine
)

func TestMain(m *testing.M) {

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Could not load .env file, using default values")
	}

	// Initialize configuration and database
	bootstrap.LoadConfig()

	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(bootstrap.AWSRegion),
		Endpoint: aws.String("http://localhost:8000"),
	}))
	db := dynamo.New(sess)

	// Initialize repositories
	documentRepo := repository.NewDynamoDocumentRepository(db)
	songRepo := repository.NewDynamoSongRepository(db, documentRepo)
	searchRepo := repository.NewDynamoSearchRepository(db)
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

	// Crear handlers
	songHandler := handlers.NewSongHandler(songService)
	docHandler := handlers.NewDocumentHandler(documentService)
	searchHandler := handlers.NewSearchHandler(searchService)
	authHandler := handlers.NewAuthHandler(authService)

	// Inicializar router para todos los tests
	TestRouter = router.SetupRouter(songHandler, docHandler, searchHandler, authHandler)

	// Ejecutar los tests
	code := m.Run()
	os.Exit(code)
}
