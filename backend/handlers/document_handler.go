package handlers

import (
	"net/http"

	"github.com/CristinaRendaLopez/rendalla-backend/middleware"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetAllDocumentsBySongIDHandler(c *gin.Context) {
	songID := c.Param("id")

	documents, err := services.GetDocumentsBySongID(songID)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to retrieve documents")
		return
	}

	logrus.WithField("song_id", songID).Info("Documents retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"data": documents})
}

func GetDocumentByIDHandler(c *gin.Context) {
	docID := c.Param("id")

	document, err := services.GetDocumentByID(docID)
	if err != nil {
		utils.HandleAPIError(c, err, "Document not found")
		return
	}

	logrus.WithField("document_id", docID).Info("Document retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"data": document})
}

func CreateDocumentHandler(c *gin.Context) {
	songID := c.Param("id")
	var document models.Document

	middleware.ValidateRequest(&document)(c)
	if c.IsAborted() {
		return
	}

	document.SongID = songID

	err := services.CreateDocument(document)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to create document")
		return
	}

	logrus.WithFields(logrus.Fields{
		"document_id": document.ID,
		"song_id":     songID,
	}).Info("Document created successfully")

	c.JSON(http.StatusCreated, gin.H{"message": "Document created successfully"})
}

func UpdateDocumentHandler(c *gin.Context) {
	docID := c.Param("id")
	var docUpdate map[string]interface{}

	middleware.ValidateRequest(&docUpdate)(c)
	if c.IsAborted() {
		return
	}

	err := services.UpdateDocument(docID, docUpdate)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to update document")
		return
	}

	logrus.WithFields(logrus.Fields{
		"document_id": docID,
		"updates":     docUpdate,
	}).Info("Document updated successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

func DeleteDocumentHandler(c *gin.Context) {
	docID := c.Param("id")

	err := services.DeleteDocument(docID)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to delete document")
		return
	}

	logrus.WithField("document_id", docID).Info("Document deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
