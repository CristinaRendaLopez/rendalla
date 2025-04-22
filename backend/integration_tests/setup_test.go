package integration_tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/app"
	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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
	path, _ := filepath.Abs("../.env.test")
	if err := godotenv.Overload(path); err != nil {
		logrus.Warnf("Could not load .env.test from path: %s", path)
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

	TestRouter = app.InitApp(db, app.AppConfig{
		JWTSecret:      os.Getenv("JWT_SECRET"),
		EnableCORS:     false,
		EnableLogger:   false,
		EnableRecovery: true,
	})

	// Execute tests
	code := m.Run()
	os.Exit(code)
}
