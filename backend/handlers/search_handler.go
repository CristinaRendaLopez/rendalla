package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/guregu/dynamo"
)

// SearchSongsByTitleHandler - GET /songs/search?title=&limit=&next_token=
func SearchSongsByTitleHandler(c *gin.Context) {
	title := c.Query("title")
	limit, nextToken := getPaginationParams(c)

	// Llamada al servicio con paginación
	songs, nextKey, err := services.SearchSongsByTitle(title, limit, nextToken)
	if err != nil {
		log.Printf("Error al buscar canciones por título '%s': %v", title, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Error al realizar la búsqueda de canciones",
				"code":    "INTERNAL_ERROR",
			},
		})
		return
	}

	// Respuesta con paginación
	c.JSON(http.StatusOK, gin.H{
		"data":       songs,
		"next_token": nextKey,
	})
}

// SearchDocumentsByTitleHandler - GET /documents/search?title=&limit=&next_token=
func SearchDocumentsByTitleHandler(c *gin.Context) {
	title := c.Query("title")
	limit, nextToken := getPaginationParams(c)

	// Llamada al servicio con paginación
	documents, nextKey, err := services.SearchDocumentsByTitle(title, limit, nextToken)
	if err != nil {
		log.Printf("Error al buscar documentos por título '%s': %v", title, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Error al realizar la búsqueda de documentos",
				"code":    "INTERNAL_ERROR",
			},
		})
		return
	}

	// Respuesta con paginación
	c.JSON(http.StatusOK, gin.H{
		"data":       documents,
		"next_token": nextKey,
	})
}

// FilterDocumentsByInstrumentHandler - GET /documents/filter?instrument=&limit=&next_token=
func FilterDocumentsByInstrumentHandler(c *gin.Context) {
	instrument := c.Query("instrument")
	limit, nextToken := getPaginationParams(c)

	// Llamada al servicio con paginación
	documents, nextKey, err := services.FilterDocumentsByInstrument(instrument, limit, nextToken)
	if err != nil {
		log.Printf("Error al filtrar documentos por instrumento '%s': %v", instrument, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": "Error al filtrar documentos",
				"code":    "INTERNAL_ERROR",
			},
		})
		return
	}

	// Respuesta con paginación
	c.JSON(http.StatusOK, gin.H{
		"data":       documents,
		"next_token": nextKey,
	})
}

// getPaginationParams - Obtiene los parámetros de paginación desde la URL
func getPaginationParams(c *gin.Context) (int, dynamo.PagingKey) {
	// Obtener el límite de resultados (por defecto: 10)
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Obtener el token de paginación si existe
	var nextToken dynamo.PagingKey
	nextTokenStr := c.Query("next_token")
	if nextTokenStr != "" {
		nextToken = dynamo.PagingKey{nextTokenStr: nil} // Convertir a `PagingKey` (guregu/dynamo)
	}

	return limit, nextToken
}
