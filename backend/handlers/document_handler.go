package handlers

import (
	"net/http"

	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type DocumentHandler struct {
	documentService services.DocumentService
}

func NewDocumentHandler(documentService *services.DocumentService) *DocumentHandler {
	return &DocumentHandler{documentService: *documentService}
}

func (h *DocumentHandler) GetAllDocumentsBySongIDHandler(c *gin.Context) {
	songID := c.Param("id")

	documents, err := h.documentService.GetDocumentsBySongID(songID)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to retrieve documents")
		return
	}

	logrus.WithField("song_id", songID).Info("Documents retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"data": documents})
}

func (h *DocumentHandler) GetDocumentByIDHandler(c *gin.Context) {
	docID := c.Param("id")

	document, err := h.documentService.GetDocumentByID(docID)
	if err != nil {
		utils.HandleAPIError(c, err, "Document not found")
		return
	}

	logrus.WithField("document_id", docID).Info("Document retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"data": document})
}

func (h *DocumentHandler) CreateDocumentHandler(c *gin.Context) {
	var document models.Document
	if err := c.ShouldBindJSON(&document); err != nil {
		utils.HandleAPIError(c, err, "Invalid request data")
		return
	}

	document.SongID = c.Param("id")

	_, err := h.documentService.CreateDocument(document)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to create document")
		return
	}

	logrus.WithFields(logrus.Fields{
		"document_id": document.ID,
		"song_id":     document.SongID,
	}).Info("Document created successfully")

	c.JSON(http.StatusCreated, gin.H{"message": "Document created successfully"})
}

func (h *DocumentHandler) UpdateDocumentHandler(c *gin.Context) {
	docID := c.Param("id")
	var docUpdate map[string]interface{}

	if err := c.ShouldBindJSON(&docUpdate); err != nil {
		utils.HandleAPIError(c, err, "Invalid request data")
		return
	}

	err := h.documentService.UpdateDocument(docID, docUpdate)
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

func (h *DocumentHandler) DeleteDocumentHandler(c *gin.Context) {
	docID := c.Param("id")

	err := h.documentService.DeleteDocument(docID)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to delete document")
		return
	}

	logrus.WithField("document_id", docID).Info("Document deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
