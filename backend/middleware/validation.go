package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

var validate = validator.New()

// ValidateRequest is a middleware function that validates incoming JSON requests
// against a given struct using go-playground/validator. If validation fails,
// it returns a 400 response with details and aborts the request.
func ValidateRequest(obj interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		newObj := obj
		if err := c.ShouldBindJSON(newObj); err != nil {
			logrus.WithError(err).Warn("Invalid request payload")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			c.Abort()
			return
		}

		if err := validate.Struct(newObj); err != nil {
			logrus.WithError(err).Warn("Validation failed")

			validationErrors := []map[string]string{}
			for _, err := range err.(validator.ValidationErrors) {
				validationErrors = append(validationErrors, map[string]string{
					"field": err.Field(),
					"error": err.Tag(),
				})
			}

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": validationErrors,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
