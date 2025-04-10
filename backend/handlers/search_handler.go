package handlers

import (
	"net/http"

	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SearchHandler handles HTTP requests for searching songs and documents.
// It delegates the business logic to the SearchServiceInterface.
type SearchHandler struct {
	searchService services.SearchServiceInterface
}

// NewSearchHandler returns a new instance of SearchHandler.
func NewSearchHandler(searchService services.SearchServiceInterface) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

// ListSongsHandler handles GET /songs/search.
// Supports filtering by title and sorting/pagination options.
func (h *SearchHandler) ListSongsHandler(c *gin.Context) {
	title := c.Query("title")
	sortField := c.Query("sort")
	sortOrder := c.Query("order")
	limit, nextToken := utils.ExtractPaginationParams(c)

	songs, nextKey, err := h.searchService.ListSongs(title, sortField, sortOrder, limit, nextToken)
	if err != nil {
		errors.HandleAPIError(c, err, "Error listing songs")
		return
	}

	logrus.WithFields(logrus.Fields{
		"title":      title,
		"sort":       sortField,
		"order":      sortOrder,
		"limit":      limit,
		"next_token": nextToken,
	}).Info("Listed songs with filters")

	c.JSON(http.StatusOK, gin.H{
		"data":       songs,
		"next_token": nextKey,
	})
}

// ListDocumentsHandler handles GET /documents/search.
// Supports filtering by title, instrument, and type, as well as sorting and pagination.
func (h *SearchHandler) ListDocumentsHandler(c *gin.Context) {
	title := c.Query("title")
	instrument := c.Query("instrument")
	docType := c.Query("type")
	sortField := c.Query("sort")
	sortOrder := c.Query("order")
	limit, nextToken := utils.ExtractPaginationParams(c)

	documents, nextKey, err := h.searchService.ListDocuments(title, instrument, docType, sortField, sortOrder, limit, nextToken)
	if err != nil {
		errors.HandleAPIError(c, err, "Error listing documents")
		return
	}

	logrus.WithFields(logrus.Fields{
		"title":      title,
		"instrument": instrument,
		"type":       docType,
		"sort":       sortField,
		"order":      sortOrder,
		"limit":      limit,
		"next_token": nextToken,
	}).Info("Listed documents with filters")

	c.JSON(http.StatusOK, gin.H{
		"data":       documents,
		"next_token": nextKey,
	})
}
