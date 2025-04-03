package main

import (
	"fmt"
	"os"

	"github.com/CristinaRendaLopez/rendalla-backend/app"
	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
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

	// Get and validate port
	port := bootstrap.AppPort
	if port == "" {
		port = "8080"
	}

	logrus.Infof("Rendalla backend is running on port %s", port)
	app.Run(fmt.Sprintf(":%s", port))
}
