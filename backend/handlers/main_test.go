package handlers_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/bootstrap"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	// Load .env.test
	path, _ := filepath.Abs("../.env.test")
	if err := godotenv.Overload(path); err != nil {
		logrus.Warnf("Could not load .env.test from path: %s", path)
	}

	// Load configuration into bootstrap
	bootstrap.LoadConfig()

	// Run all tests
	os.Exit(m.Run())
}
