package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
)

// ExtractPaginationParams parses pagination-related query parameters from the request context.
//
// Supported query parameters:
//   - limit: max number of results to return (defaults to 10, minimum 1)
//   - next_token: optional pagination token for DynamoDB paging
//
// Returns:
//   - limit as an integer
//   - nextToken as a dynamo.PagingKey (empty if not provided)
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
