package utils

import (
	"strings"

	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/gin-gonic/gin"
)

// RequireParam checks that a path param is present and non-empty.
// If missing, it sends a 400 error and returns false.
func RequireParam(c *gin.Context, name string) (string, bool) {
	val := strings.TrimSpace(c.Param(name))
	if val == "" {
		errors.HandleAPIError(c, errors.ErrBadRequest, "Missing "+name)
		return "", false
	}
	return val, true
}

// RequireQuery checks that a query parameter is present and non-empty.
// If missing, it sends a 400 error and returns false.
func RequireQuery(c *gin.Context, key string) (string, bool) {
	val := strings.TrimSpace(c.Query(key))
	if val == "" {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing "+key+" parameter")
		return "", false
	}
	return val, true
}

// ValidateNonEmptyStringField validates that a field in a map is a non-empty string.
func ValidateNonEmptyStringField(update map[string]interface{}, key string) error {
	val, ok := update[key]
	if !ok {
		return nil
	}
	strVal, valid := val.(string)
	if !valid || strings.TrimSpace(strVal) == "" {
		return errors.ErrValidationFailed
	}
	return nil
}

// ValidateNonEmptyStringArrayField validates that a field in a map is a non-empty []string
func ValidateNonEmptyStringArrayField(update map[string]interface{}, key string) error {
	val, ok := update[key]
	if !ok {
		return nil
	}

	array, ok := val.([]interface{})
	if !ok || len(array) == 0 {
		return errors.ErrValidationFailed
	}

	for _, item := range array {
		str, valid := item.(string)
		if !valid || strings.TrimSpace(str) == "" {
			return errors.ErrValidationFailed
		}
	}

	return nil
}

func IsEmptyString(val string) bool {
	return strings.TrimSpace(val) == ""
}
