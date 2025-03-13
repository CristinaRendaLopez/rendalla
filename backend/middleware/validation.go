package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

// Validator instance
var validate = validator.New()

// ValidateRequest middleware to validate incoming JSON requests
func ValidateRequest(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON to the provided struct
		if err := c.ShouldBindJSON(obj); err != nil {
			logrus.WithError(err).Warn("Invalid request payload")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			c.Abort()
			return
		}

		// Validate the struct using validator.v10
		if err := validate.Struct(obj); err != nil {
			logrus.WithError(err).Warn("Validation failed")

			// Extract validation errors
			validationErrors := make(map[string]string)
			for _, err := range err.(validator.ValidationErrors) {
				validationErrors[err.Field()] = err.Tag()
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationErrors,
			})
			c.Abort()
			return
		}

		// Proceed if validation passes
		c.Next()
	}
}
