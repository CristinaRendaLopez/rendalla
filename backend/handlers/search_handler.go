package handlers

import (
	"net/http"
	"strconv"

	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
)

func SearchSongsByTitleHandler(c *gin.Context) {
	title := c.Query("title")
	limit, nextToken := getPaginationParams(c)

	songs, nextKey, err := services.SearchSongsByTitle(title, limit, nextToken)
	if err != nil {
		utils.HandleAPIError(c, err, "Error searching for songs")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       songs,
		"next_token": nextKey,
	})
}

func SearchDocumentsByTitleHandler(c *gin.Context) {
	title := c.Query("title")
	limit, nextToken := getPaginationParams(c)

	documents, nextKey, err := services.SearchDocumentsByTitle(title, limit, nextToken)
	if err != nil {
		utils.HandleAPIError(c, err, "Error searching for documents")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       documents,
		"next_token": nextKey,
	})
}

func FilterDocumentsByInstrumentHandler(c *gin.Context) {
	instrument := c.Query("instrument")
	limit, nextToken := getPaginationParams(c)

	documents, nextKey, err := services.FilterDocumentsByInstrument(instrument, limit, nextToken)
	if err != nil {
		utils.HandleAPIError(c, err, "Error filtering documents by instrument")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       documents,
		"next_token": nextKey,
	})
}

func getPaginationParams(c *gin.Context) (int, dynamo.PagingKey) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	var nextToken dynamo.PagingKey
	nextTokenStr := c.Query("next_token")
	if nextTokenStr != "" {
		nextToken = dynamo.PagingKey{nextTokenStr: nil}
	}

	return limit, nextToken
}
