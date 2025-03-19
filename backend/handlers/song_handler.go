package handlers

import (
	"net/http"

	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/services"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SongHandler struct {
	songService *services.SongService
}

func NewSongHandler(songService *services.SongService) *SongHandler {
	return &SongHandler{songService: songService}
}

type SongRequest struct {
	Title     string            `json:"title" binding:"required,min=3"`
	Author    string            `json:"author" binding:"required"`
	Genres    []string          `json:"genres" binding:"required,dive,min=3"`
	Documents []models.Document `json:"documents,omitempty"`
}

func (h *SongHandler) GetAllSongsHandler(c *gin.Context) {
	songs, err := h.songService.GetAllSongs()
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to retrieve songs")
		return
	}

	logrus.Info("Fetched all songs successfully")
	c.JSON(http.StatusOK, gin.H{"data": songs})
}

func (h *SongHandler) GetSongByIDHandler(c *gin.Context) {
	id := c.Param("id")

	song, err := h.songService.GetSongByID(id)
	if err != nil {
		utils.HandleAPIError(c, err, "Song not found")
		return
	}

	logrus.WithField("song_id", id).Info("Fetched song successfully")
	c.JSON(http.StatusOK, gin.H{"data": song})
}

func (h *SongHandler) CreateSongHandler(c *gin.Context) {
	var req SongRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleAPIError(c, err, "Invalid song data")
		return
	}

	song := models.Song{
		Title:  req.Title,
		Author: req.Author,
		Genres: req.Genres,
	}

	_, err := h.songService.CreateSongWithDocuments(song, req.Documents)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to create song")
		return
	}

	logrus.WithFields(logrus.Fields{"title": req.Title, "author": req.Author}).Info("Song created successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "Song created successfully"})
}

func (h *SongHandler) UpdateSongHandler(c *gin.Context) {
	id := c.Param("id")
	var songUpdate map[string]interface{}

	if err := c.ShouldBindJSON(&songUpdate); err != nil {
		utils.HandleAPIError(c, err, "Invalid request data")
		return
	}

	err := h.songService.UpdateSong(id, songUpdate)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to update song")
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id": id,
		"updates": songUpdate,
	}).Info("Song updated successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

func (h *SongHandler) DeleteSongWithDocumentsHandler(c *gin.Context) {
	id := c.Param("id")

	err := h.songService.DeleteSongWithDocuments(id)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to delete song")
		return
	}

	logrus.WithField("song_id", id).Info("Song and associated documents deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
