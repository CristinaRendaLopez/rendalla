package utils

import (
	"strings"

	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/gin-gonic/gin"
)

// RequireParam checks that a path param is present and non-empty.
// If missing, it sends a 400 error and returns false.
func RequireParam(c *gin.Context, name string) (string, bool) {
	val := strings.TrimSpace(c.Param(name))
	if val == "" {
		HandleAPIError(c, ErrBadRequest, "Missing "+name)
		return "", false
	}
	return val, true
}

// RequireQuery checks that a query parameter is present and non-empty.
// If missing, it sends a 400 error and returns false.
func RequireQuery(c *gin.Context, key string) (string, bool) {
	val := strings.TrimSpace(c.Query(key))
	if val == "" {
		HandleAPIError(c, ErrValidationFailed, "Missing "+key+" parameter")
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
		return ErrValidationFailed
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
		return ErrValidationFailed
	}

	for _, item := range array {
		str, valid := item.(string)
		if !valid || strings.TrimSpace(str) == "" {
			return ErrValidationFailed
		}
	}

	return nil
}

func IsEmptyString(val string) bool {
	return strings.TrimSpace(val) == ""
}

func ValidateSong(song models.Song) error {
	if IsEmptyString(song.Title) || len(song.Title) < 3 {
		return ErrValidationFailed
	}
	if IsEmptyString(song.Author) {
		return ErrValidationFailed
	}
	if len(song.Genres) == 0 {
		return ErrValidationFailed
	}
	for _, g := range song.Genres {
		if len(g) < 3 {
			return ErrValidationFailed
		}
	}
	return nil
}

func ValidateDocument(doc models.Document) error {
	if IsEmptyString(doc.Type) {
		return ErrValidationFailed
	}
	if IsEmptyString(doc.PDFURL) {
		return ErrValidationFailed
	}
	if len(doc.Instrument) == 0 {
		return ErrValidationFailed
	}
	for _, inst := range doc.Instrument {
		if IsEmptyString(inst) {
			return ErrValidationFailed
		}
	}
	return nil
}

func ValidateSongAndDocuments(song models.Song, documents []models.Document) error {
	if err := ValidateSong(song); err != nil {
		return err
	}

	for _, doc := range documents {
		if err := ValidateDocument(doc); err != nil {
			return err
		}
	}

	return nil
}

func ValidateDocumentUpdate(update map[string]interface{}) error {
	if err := ValidateNonEmptyStringField(update, "type"); err != nil {
		return err
	}
	if err := ValidateNonEmptyStringField(update, "pdf_url"); err != nil {
		return err
	}
	if err := ValidateNonEmptyStringArrayField(update, "instrument"); err != nil {
		return err
	}
	return nil
}

func ValidateSongUpdate(update map[string]interface{}) error {
	if len(update) == 0 {
		return ErrValidationFailed
	}
	if err := ValidateNonEmptyStringField(update, "title"); err != nil {
		return err
	}
	if err := ValidateNonEmptyStringField(update, "author"); err != nil {
		return err
	}
	if err := ValidateNonEmptyStringArrayField(update, "genres"); err != nil {
		return err
	}
	return nil
}
