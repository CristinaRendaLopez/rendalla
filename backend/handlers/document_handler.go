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
	docID, ok := utils.RequireParam(c, "id")
	if !ok {
		return
	}

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
		utils.HandleAPIError(c, utils.ErrValidationFailed, "Invalid request data")
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
	docID, ok := utils.RequireParam(c, "id")
	if !ok {
		return
	}

	var docUpdate map[string]interface{}
	if err := c.ShouldBindJSON(&docUpdate); err != nil {
		utils.HandleAPIError(c, utils.ErrValidationFailed, "Invalid request data")
		return
	}

	if err := utils.ValidateDocumentUpdate(docUpdate); err != nil {
		utils.HandleAPIError(c, err, "Invalid document update data")
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
	docID, ok := utils.RequireParam(c, "id")
	if !ok {
		return
	}

	err := h.documentService.DeleteDocument(docID)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to delete document")
		return
	}

	logrus.WithField("document_id", docID).Info("Document deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
