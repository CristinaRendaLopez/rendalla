package integration_tests

import (
	"os"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/repository"
	"github.com/CristinaRendaLopez/rendalla-backend/router"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	TestRouter       *gin.Engine
	TestTimeProvider utils.TimeProvider
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

	if err := CreateTestTables(db); err != nil {
		logrus.Fatal("Could not create test tables:", err)
	}

	TestTimeProvider = &utils.UTCTimeProvider{}

	if err := SeedTestData(db, TestTimeProvider); err != nil {
		logrus.Fatal("Failed to seed test data:", err)
	}

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

func CreateTestTables(db *dynamo.DB) error {
	svc := db.Client()

	tables := []struct {
		Name       string
		PrimaryKey string
	}{
		{bootstrap.SongTableName, "id"},
		{bootstrap.DocumentTableName, "id"},
	}

	for _, t := range tables {
		_, err := svc.CreateTable(&dynamodb.CreateTableInput{
			TableName: aws.String(t.Name),
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String(t.PrimaryKey),
					KeyType:       aws.String("HASH"), // Partition key
				},
			},
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String(t.PrimaryKey),
					AttributeType: aws.String("S"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(5),
				WriteCapacityUnits: aws.Int64(5),
			},
		})

		if err != nil && !isTableAlreadyExists(err) {
			logrus.WithFields(logrus.Fields{"table": t.Name, "error": err}).Error("Failed to create table")
			return err
		}

		logrus.WithField("table", t.Name).Info("Table is ready")
	}

	return nil
}

func isTableAlreadyExists(err error) bool {
	if err == nil {
		return false
	}
	return dynamodb.ErrCodeResourceInUseException == aws.StringValue(aws.String(err.Error()))
}

func SeedTestData(db *dynamo.DB, timeProvider utils.TimeProvider) error {
	now := timeProvider.Now()

	song := models.Song{
		ID:         "queen-001",
		Title:      "Bohemian Rhapsody",
		Author:     "Queen",
		Genres:     []string{"Rock", "Progressive"},
		YoutubeURL: "https://youtube.com/watch?v=fJ9rUzIMcZQ",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	documents := []models.Document{
		{
			ID:         "doc-br-piano",
			SongID:     song.ID,
			Type:       "partitura",
			Instrument: []string{"piano"},
			PDFURL:     "https://s3.test/bohemian_rhapsody_piano.pdf",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		{
			ID:         "doc-br-voice",
			SongID:     song.ID,
			Type:       "tablatura",
			Instrument: []string{"voz"},
			PDFURL:     "https://s3.test/bohemian_rhapsody_voz.pdf",
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	if err := db.Table(bootstrap.SongTableName).Put(song).Run(); err != nil {
		logrus.WithError(err).Error("Failed to seed Queen song")
		return err
	}

	for _, doc := range documents {
		if err := db.Table(bootstrap.DocumentTableName).Put(doc).Run(); err != nil {
			logrus.WithError(err).WithField("doc_id", doc.ID).Error("Failed to seed document")
			return err
		}
	}

	logrus.Info("Test data seeded successfully")
	return nil
}
