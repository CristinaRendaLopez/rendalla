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

// SongHandler handles HTTP requests related to song operations.
// It delegates the business logic to the SongServiceInterface.
type SongHandler struct {
	songService services.SongServiceInterface
}

// NewSongHandler returns a new instance of SongHandler.
func NewSongHandler(songService services.SongServiceInterface) *SongHandler {
	return &SongHandler{songService: songService}
}

// CreateSongHandler handles POST /songs.
// Delegates validation and creation to the service layer.
func (h *SongHandler) CreateSongHandler(c *gin.Context) {
	var req dto.CreateSongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Invalid JSON payload")
		return
	}

	songID, err := h.songService.CreateSongWithDocuments(req)
	if err != nil {
		errors.HandleAPIError(c, err, "Failed to create song")
		return
	}

	logrus.WithFields(logrus.Fields{"title": req.Title, "author": req.Author}).Info("Song created successfully")
	c.JSON(http.StatusCreated, gin.H{
		"message": "Song created successfully",
		"song_id": songID,
	})
}

// GetAllSongsHandler handles GET /songs.
// Returns all songs stored in the system.
func (h *SongHandler) GetAllSongsHandler(c *gin.Context) {
	songs, err := h.songService.GetAllSongs()
	if err != nil {
		errors.HandleAPIError(c, err, "Failed to retrieve songs")
		return
	}

	logrus.Info("Fetched all songs successfully")
	c.JSON(http.StatusOK, gin.H{"data": songs})
}

// GetSongByIDHandler handles GET /songs/:song_id.
// Retrieves a song by its unique ID.
func (h *SongHandler) GetSongByIDHandler(c *gin.Context) {
	id, ok := utils.RequireParam(c, "song_id")
	if !ok {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing parameter: song_id")
		return
	}

	song, err := h.songService.GetSongByID(id)
	if err != nil {
		msg := "Failed to retrieve song"
		if stdErrors.Is(err, errors.ErrResourceNotFound) {
			msg = "Song not found"
		}
		errors.HandleAPIError(c, err, msg)
		return
	}

	logrus.WithField("song_id", id).Info("Fetched song successfully")
	c.JSON(http.StatusOK, gin.H{"data": song})
}

// UpdateSongHandler handles PUT /songs/:song_id.
// Delegates validation and update logic to the service layer.
func (h *SongHandler) UpdateSongHandler(c *gin.Context) {
	id, ok := utils.RequireParam(c, "song_id")
	if !ok {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing parameter: song_id")
		return
	}

	var songUpdate dto.UpdateSongRequest
	if err := c.ShouldBindJSON(&songUpdate); err != nil {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Invalid JSON payload")
		return
	}

	if err := dto.ValidateUpdateSongRequest(songUpdate); err != nil {
		errors.HandleAPIError(c, err, "Invalid update payload")
		return
	}

	if err := h.songService.UpdateSong(id, songUpdate); err != nil {
		errors.HandleAPIError(c, err, "Failed to update song")
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id": id,
		"updates": songUpdate,
	}).Info("Song updated successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

// DeleteSongWithDocumentsHandler handles DELETE /songs/:song_id.
// Deletes the song and all documents linked to it.
func (h *SongHandler) DeleteSongWithDocumentsHandler(c *gin.Context) {
	id, ok := utils.RequireParam(c, "song_id")
	if !ok {
		errors.HandleAPIError(c, errors.ErrValidationFailed, "Missing parameter: song_id")
		return
	}

	err := h.songService.DeleteSongWithDocuments(id)
	if err != nil {
		errors.HandleAPIError(c, err, "Failed to delete song")
		return
	}

	logrus.WithField("song_id", id).Info("Song and associated documents deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
