package handlers

import (
	"net/http"

	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SearchHandler struct {
	searchService services.SearchServiceInterface
}

func NewSearchHandler(searchService services.SearchServiceInterface) *SearchHandler {
	return &SearchHandler{searchService: searchService}
}

func (h *SearchHandler) ListSongsHandler(c *gin.Context) {
	title := c.Query("title")
	sortField := c.DefaultQuery("sort", "created_at")
	sortOrder := c.DefaultQuery("order", "desc")
	limit, nextToken := utils.ExtractPaginationParams(c)

	songs, nextKey, err := h.searchService.ListSongs(title, sortField, sortOrder, limit, nextToken)
	if err != nil {
		utils.HandleAPIError(c, err, "Error listing songs")
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

func (h *SearchHandler) SearchDocumentsByTitleHandler(c *gin.Context) {
	title, ok := utils.RequireQuery(c, "title")
	if !ok {
		return
	}

	limit, nextToken := utils.ExtractPaginationParams(c)

	documents, nextKey, err := h.searchService.SearchDocumentsByTitle(title, limit, nextToken)
	if err != nil {
		utils.HandleAPIError(c, err, "Error searching for documents")
		return
	}

	logrus.WithFields(logrus.Fields{
		"title":      title,
		"limit":      limit,
		"next_token": nextToken,
	}).Info("Searched documents by title")

	c.JSON(http.StatusOK, gin.H{
		"data":       documents,
		"next_token": nextKey,
	})
}

func (h *SearchHandler) FilterDocumentsByInstrumentHandler(c *gin.Context) {
	instrument, ok := utils.RequireQuery(c, "instrument")
	if !ok {
		return
	}

	limit, nextToken := utils.ExtractPaginationParams(c)

	documents, nextKey, err := h.searchService.FilterDocumentsByInstrument(instrument, limit, nextToken)
	if err != nil {
		utils.HandleAPIError(c, err, "Error filtering documents by instrument")
		return
	}

	logrus.WithFields(logrus.Fields{
		"instrument": instrument,
		"limit":      limit,
		"next_token": nextToken,
	}).Info("Filtered documents by instrument")

	c.JSON(http.StatusOK, gin.H{
		"data":       documents,
		"next_token": nextKey,
	})
}
