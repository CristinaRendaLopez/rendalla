package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
)

func ExtractPaginationParams(c *gin.Context) (int, dynamo.PagingKey) {
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
