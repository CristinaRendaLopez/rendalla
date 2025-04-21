package integration_tests

import (
	"os"
	"path/filepath"

	"github.com/CristinaRendaLopez/rendalla-backend/app"
	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	DB           *dynamo.DB
	TimeProvider utils.TimeProvider
	Router       *gin.Engine
}

func (s *IntegrationTestSuite) SetupSuite() {
	path, _ := filepath.Abs("../.env.test")
	if err := godotenv.Overload(path); err != nil {
		logrus.Warnf("Could not load .env.test from path: %s", path)
	}

	bootstrap.LoadConfig()

	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(bootstrap.AWSRegion),
		Endpoint: aws.String("http://localhost:8000"),
	}))
	s.DB = dynamo.New(sess)

	if err := CreateTestTables(s.DB); err != nil {
		s.FailNow("failed to create tables", err)
	}

	s.TimeProvider = &utils.UTCTimeProvider{}

	if err := SeedTestData(s.DB, s.TimeProvider); err != nil {
		s.FailNow("failed to seed test data", err)
	}

	s.Router = app.InitApp(s.DB, app.AppConfig{
		JWTSecret:      os.Getenv("JWT_SECRET"),
		EnableCORS:     false,
		EnableLogger:   false,
		EnableRecovery: true,
	})
}

func (s *IntegrationTestSuite) TearDownSuite() {
	svc, ok := s.DB.Client().(*dynamodb.DynamoDB)
	if !ok {
		logrus.Warn("could not assert client as *dynamodb.DynamoDB for teardown")
		return
	}
	for _, table := range []string{bootstrap.SongTableName, bootstrap.DocumentTableName} {
		_, err := svc.DeleteTable(&dynamodb.DeleteTableInput{TableName: aws.String(table)})
		if err != nil {
			logrus.WithError(err).Warnf("could not delete table %s during teardown", table)
		}
	}
}
