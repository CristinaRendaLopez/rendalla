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

type SongRequest struct {
	Title     string            `json:"title" binding:"required,min=3"`
	Author    string            `json:"author" binding:"required"`
	Genres    []string          `json:"genres" binding:"required,dive,min=3"`
	Documents []models.Document `json:"documents,omitempty"` // Optional field
}

func GetAllSongsHandler(c *gin.Context) {
	songs, err := services.GetAllSongs()
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to retrieve songs")
		return
	}

	logrus.Info("Fetched all songs successfully")
	c.JSON(http.StatusOK, gin.H{"data": songs})
}

func GetSongByIDHandler(c *gin.Context) {
	id := c.Param("id")

	song, err := services.GetSongByID(id)
	if err != nil {
		utils.HandleAPIError(c, err, "Song not found")
		return
	}

	logrus.WithField("song_id", id).Info("Fetched song successfully")
	c.JSON(http.StatusOK, gin.H{"data": song})
}

func CreateSongHandler(c *gin.Context) {
	var req SongRequest

	middleware.ValidateRequest(&req)(c)
	if c.IsAborted() {
		return
	}

	song := models.Song{
		Title:  req.Title,
		Author: req.Author,
		Genres: req.Genres,
	}

	err := services.CreateSongWithDocuments(song, req.Documents)
	if err != nil {
		logrus.WithError(err).Error("Failed to create song with documents")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create song"})
		return
	}

	logrus.WithFields(logrus.Fields{"title": req.Title, "author": req.Author}).Info("Song created successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "Song created successfully"})
}

func UpdateSongHandler(c *gin.Context) {
	id := c.Param("id")
	var songUpdate map[string]interface{}

	middleware.ValidateRequest(&songUpdate)(c)
	if c.IsAborted() {
		return
	}

	err := services.UpdateSong(id, songUpdate)
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

func DeleteSongWithDocumentsHandler(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteSongWithDocuments(id)
	if err != nil {
		utils.HandleAPIError(c, err, "Failed to delete song")
		return
	}

	logrus.WithField("song_id", id).Info("Song and associated documents deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
