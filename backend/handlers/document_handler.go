package handlers

import (
	"errors"
	"net/http"

	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type DocumentHandler struct {
	documentService services.DocumentServiceInterface
}

func NewDocumentHandler(documentService services.DocumentServiceInterface) *DocumentHandler {
	return &DocumentHandler{documentService: documentService}
}

func (h *DocumentHandler) GetAllDocumentsBySongIDHandler(c *gin.Context) {
	songID, ok := utils.RequireParam(c, "id")
	if !ok {
		return
	}

	documents, err := h.documentService.GetDocumentsBySongID(songID)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to retrieve documents")
		return
	}

	logrus.WithField("song_id", songID).Info("Documents retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"data": documents})
}

func (h *DocumentHandler) GetDocumentByIDHandler(c *gin.Context) {
	songID, ok := utils.RequireParam(c, "id")
	if !ok {
		return
	}

	docID, ok := utils.RequireParam(c, "doc_id")
	if !ok {
		return
	}

	document, err := h.documentService.GetDocumentByID(songID, docID)
	if err != nil {
		utils.HandleAPIError(c, err, "Document not found")
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id":     songID,
		"document_id": docID,
	}).Info("Document retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"data": document})
}

func (h *DocumentHandler) CreateDocumentHandler(c *gin.Context) {
	var document models.Document
	if err := c.ShouldBindJSON(&document); err != nil {
		logrus.WithError(err).Warn("Invalid JSON payload")
		utils.HandleAPIError(c, utils.ErrValidationFailed, "Invalid JSON payload")
		return
	}

	songID, ok := utils.RequireParam(c, "id")
	if !ok {
		return
	}
	document.SongID = songID

	if err := utils.ValidateDocument(document); err != nil {
		utils.HandleAPIError(c, err, "Invalid document data")
		return
	}

	documentID, err := h.documentService.CreateDocument(document)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to create document")
		return
	}

	logrus.WithFields(logrus.Fields{
		"document_id": documentID,
		"song_id":     document.SongID,
	}).Info("Document created successfully")

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Document created successfully",
		"document_id": documentID,
	})
}

func (h *DocumentHandler) UpdateDocumentHandler(c *gin.Context) {
	songID, ok := utils.RequireParam(c, "id")
	if !ok {
		return
	}
	docID, ok := utils.RequireParam(c, "doc_id")
	if !ok {
		return
	}

	var docUpdate map[string]interface{}
	if err := c.ShouldBindJSON(&docUpdate); err != nil {
		logrus.WithError(err).Warn("Invalid JSON payload")
		utils.HandleAPIError(c, utils.ErrValidationFailed, "Invalid JSON payload")
		return
	}

	if err := utils.ValidateDocumentUpdate(docUpdate); err != nil {
		utils.HandleAPIError(c, err, "Invalid document update data")
		return
	}

	err := h.documentService.UpdateDocument(songID, docID, docUpdate)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to update document")
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id":     songID,
		"document_id": docID,
		"updates":     docUpdate,
	}).Info("Document updated successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

func (h *DocumentHandler) DeleteDocumentHandler(c *gin.Context) {
	songID, ok := utils.RequireParam(c, "id")
	if !ok {
		return
	}

	docID, ok := utils.RequireParam(c, "doc_id")
	if !ok {
		return
	}

	err := h.documentService.DeleteDocument(songID, docID)
	if err != nil {
		message := "Failed to delete document"
		if errors.Is(err, utils.ErrResourceNotFound) {
			message = "Document not found"
		}
		utils.HandleAPIError(c, err, message)
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id":     songID,
		"document_id": docID,
	}).Info("Document deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
