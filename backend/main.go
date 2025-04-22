package main

import (
	"fmt"
	"os"

	"github.com/CristinaRendaLopez/rendalla-backend/app"
	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
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

	app := app.InitApp(bootstrap.DB, app.AppConfig{
		JWTSecret:      os.Getenv("JWT_SECRET"),
		EnableCORS:     true,
		EnableLogger:   true,
		EnableRecovery: true,
	})

	if os.Getenv("LAMBDA_TASK_ROOT") != "" {
		logrus.Info("Running in AWS Lambda mode")
		lambdaAdapter := ginadapter.New(app)
		lambda.Start(func(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			return lambdaAdapter.Proxy(req)
		})
	} else {
		// Get and validate port
		port := bootstrap.AppPort
		if port == "" {
			port = "8080"
		}

		logrus.Infof("Rendalla backend is running on port %s", port)
		app.Run(fmt.Sprintf(":%s", port))
	}
}
