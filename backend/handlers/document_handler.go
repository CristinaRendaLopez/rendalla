package handlers

import (
	"net/http"

	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetAllDocumentsBySongIDHandler(c *gin.Context) {
	songID := c.Param("id")

	documents, err := services.GetDocumentsBySongID(songID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": songID, "error": err}).Error("Failed to retrieve documents")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve documents"})
		return
	}

	logrus.WithField("song_id", songID).Info("Documents retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"data": documents})
}

func GetDocumentByIDHandler(c *gin.Context) {
	docID := c.Param("id")

	document, err := services.GetDocumentByID(docID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"document_id": docID, "error": err}).Error("Document not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	logrus.WithField("document_id", docID).Info("Document retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"data": document})
}

func CreateDocumentHandler(c *gin.Context) {
	songID := c.Param("id")
	var document models.Document

	if err := c.ShouldBindJSON(&document); err != nil {
		logrus.WithError(err).Warn("Invalid request body for document creation")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	document.SongID = songID

	err := services.CreateDocument(document)
	if err != nil {
		logrus.WithFields(logrus.Fields{"song_id": songID, "error": err}).Error("Failed to create document")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create document"})
		return
	}

	logrus.WithFields(logrus.Fields{"document_id": document.ID, "song_id": songID}).Info("Document created successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "Document created successfully"})
}

func UpdateDocumentHandler(c *gin.Context) {
	docID := c.Param("id")
	var docUpdate map[string]interface{}

	if err := c.ShouldBindJSON(&docUpdate); err != nil {
		logrus.WithFields(logrus.Fields{"document_id": docID, "error": err}).Warn("Invalid request body for document update")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := services.UpdateDocument(docID, docUpdate)
	if err != nil {
		logrus.WithFields(logrus.Fields{"document_id": docID, "error": err}).Error("Failed to update document")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document"})
		return
	}

	logrus.WithField("document_id", docID).Info("Document updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

func DeleteDocumentHandler(c *gin.Context) {
	docID := c.Param("id")

	err := services.DeleteDocument(docID)
	if err != nil {
		logrus.WithFields(logrus.Fields{"document_id": docID, "error": err}).Error("Failed to delete document")
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found or deletion error"})
		return
	}

	logrus.WithField("document_id", docID).Info("Document deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
