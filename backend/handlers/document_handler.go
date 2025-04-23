package handlers

import (
	stdErrors "errors"
	"net/http"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// DocumentHandler handles HTTP requests related to documents (scores or tablatures).
// It delegates business logic to the DocumentServiceInterface.
type DocumentHandler struct {
	documentService services.DocumentServiceInterface
	fileService     utils.FileUploader
}

// NewDocumentHandler returns a new instance of DocumentHandler.
func NewDocumentHandler(documentService services.DocumentServiceInterface, fileService utils.FileUploader) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
		fileService:     fileService,
	}
}

// CreateDocumentHandler handles POST /songs/:song_id/documents.
// Validates the request and creates a new document linked to a song.
func (h *DocumentHandler) CreateDocumentHandler(c *gin.Context) {

	songID, ok := utils.RequireParam(c, "song_id")
	if !ok {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing parameter: song_id")
		return
	}

	docType := c.PostForm("type")
	instruments := c.PostFormArray("instrument[]")

	fileHeader, err := c.FormFile("pdf")
	if err != nil {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing or invalid PDF file")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Unable to open uploaded PDF")
		return
	}
	defer file.Close()

	pdfURL, err := h.fileService.UploadPDFToS3(file, fileHeader, songID)
	if err != nil {
		errors.HandleAPIError(c, err, "Failed to upload PDF")
		return
	}

	req := dto.CreateDocumentRequest{
		Type:       docType,
		Instrument: instruments,
		PDFURL:     pdfURL,
		SongID:     songID,
	}

	documentID, err := h.documentService.CreateDocument(req)
	if err != nil {
		errors.HandleAPIError(c, err, "Failed to create document")
		return
	}

	logrus.WithFields(logrus.Fields{
		"document_id": documentID,
		"song_id":     songID,
	}).Info("Document created successfully")

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Document created successfully",
		"document_id": documentID,
	})
}

// GetAllDocumentsBySongIDHandler handles GET /songs/:song_id/documents.
// Retrieves all documents associated with a specific song.
func (h *DocumentHandler) GetAllDocumentsBySongIDHandler(c *gin.Context) {
	songID, ok := utils.RequireParam(c, "song_id")
	if !ok {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing parameter: song_id")
		return
	}

	documents, err := h.documentService.GetDocumentsBySongID(songID)
	if err != nil {
		errors.HandleAPIError(c, err, "Failed to retrieve documents")
		return
	}

	logrus.WithField("song_id", songID).Info("Documents retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"data": documents})
}

// GetDocumentByIDHandler handles GET /songs/:song_id/documents/:doc_id.
// Retrieves a single document by song ID and document ID.
func (h *DocumentHandler) GetDocumentByIDHandler(c *gin.Context) {
	songID, ok := utils.RequireParam(c, "song_id")
	if !ok {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing parameter: song_id")
		return
	}

	docID, ok := utils.RequireParam(c, "doc_id")
	if !ok {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing parameter: doc_id")
		return
	}

	document, err := h.documentService.GetDocumentByID(songID, docID)
	if err != nil {
		message := "Failed to retrieve document"
		if stdErrors.Is(err, errors.ErrResourceNotFound) {
			message = "Document not found"
		}
		errors.HandleAPIError(c, err, message)
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id":     songID,
		"document_id": docID,
	}).Info("Document retrieved successfully")
	c.JSON(http.StatusOK, gin.H{"data": document})
}

// UpdateDocumentHandler handles PUT /songs/:song_id/documents/:doc_id.
// Applies updates to a specific document.
func (h *DocumentHandler) UpdateDocumentHandler(c *gin.Context) {
	songID, ok := utils.RequireParam(c, "song_id")
	if !ok {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing parameter: song_id")
		return
	}
	docID, ok := utils.RequireParam(c, "doc_id")
	if !ok {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing parameter: doc_id")
		return
	}

	var docUpdate dto.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&docUpdate); err != nil {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Invalid JSON payload")
		return
	}

	if err := dto.ValidateUpdateDocumentRequest(docUpdate); err != nil {
		errors.HandleAPIError(c, err, "Invalid document update data")
		return
	}

	err := h.documentService.UpdateDocument(songID, docID, docUpdate)
	if err != nil {
		errors.HandleAPIError(c, err, "Failed to update document")
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id":     songID,
		"document_id": docID,
		"updates":     docUpdate,
	}).Info("Document updated successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Document updated successfully"})
}

// DeleteDocumentHandler handles DELETE /songs/:song_id/documents/:doc_id.
// Deletes a specific document linked to a song.
func (h *DocumentHandler) DeleteDocumentHandler(c *gin.Context) {
	songID, ok := utils.RequireParam(c, "song_id")
	if !ok {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing parameter: song_id")
		return
	}

	docID, ok := utils.RequireParam(c, "doc_id")
	if !ok {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing parameter: doc_id")
		return
	}

	err := h.documentService.DeleteDocument(songID, docID)
	if err != nil {
		message := "Failed to delete document"
		if stdErrors.Is(err, errors.ErrResourceNotFound) {
			message = "Document not found"
		}
		errors.HandleAPIError(c, err, message)
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id":     songID,
		"document_id": docID,
	}).Info("Document deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}
