package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSongHandlerTest() (*handlers.SongHandler, *mocks.MockSongService) {
	mockService := new(mocks.MockSongService)
	handler := handlers.NewSongHandler(mockService)
	return handler, mockService
}

func DecodeJSONResponse[T any](w *httptest.ResponseRecorder) (T, error) {
	var body T
	err := json.Unmarshal(w.Body.Bytes(), &body)
	return body, err
}

func TestCreateSongHandler(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		setupMock      bool
		mockReturnID   string
		mockReturnErr  error
		expectedCode   int
		expectedSongID string
	}{
		{
			name:           "success",
			input:          SongValidJSON,
			setupMock:      true,
			mockReturnID:   "123",
			mockReturnErr:  nil,
			expectedCode:   http.StatusCreated,
			expectedSongID: "123",
		},
		{
			name:         "invalid JSON",
			input:        SongInvalidJSON,
			setupMock:    false,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:          "validation fails",
			input:         SongInvalidDataJSON,
			setupMock:     true,
			mockReturnID:  "",
			mockReturnErr: errors.ErrValidationFailed,
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:          "internal service error",
			input:         SongValidJSON,
			setupMock:     true,
			mockReturnID:  "",
			mockReturnErr: errors.ErrInternalServer,
			expectedCode:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupSongHandlerTest()

			if tt.setupMock {
				mockService.
					On("CreateSongWithDocuments", mock.Anything, mock.Anything).
					Return(tt.mockReturnID, tt.mockReturnErr)
			}

			c, w := utils.CreateTestContext(http.MethodPost, "/songs", strings.NewReader(tt.input))
			handler.CreateSongHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusCreated {

				response, err := DecodeJSONResponse[dto.CreateSongResponse](w)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSongID, response.SongID)
			}

			mockService.AssertExpectations(t)
		})
	}
}

// func TestGetAllSongsHandler_Success(t *testing.T) {
// 	handler, mockService := setupSongHandlerTest()

// 	expectedSongs := []models.Song{
// 		{ID: "1", Title: "Bohemian Rhapsody", Author: "Queen"},
// 		{ID: "2", Title: "Imagine", Author: "John Lennon"},
// 	}
// 	mockService.On("GetAllSongs").Return(expectedSongs, nil)

// 	c, w := utils.CreateTestContext(http.MethodGet, "/songs", nil)
// 	handler.GetAllSongsHandler(c)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Contains(t, w.Body.String(), "Bohemian Rhapsody")
// 	mockService.AssertExpectations(t)
// }

func TestGetAllSongsHandler_Service_Error(t *testing.T) {
	handler, mockService := setupSongHandlerTest()
	mockService.On("GetAllSongs").Return([]dto.SongResponseItem{}, assert.AnError)

	c, w := utils.CreateTestContext(http.MethodGet, "/songs", nil)
	handler.GetAllSongsHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

// func TestGetSongByIDHandler_Success(t *testing.T) {
// 	handler, mockService := setupSongHandlerTest()
// 	mockService.On("GetSongByID", "1").Return(&models.Song{ID: "1", Title: "Song"}, nil)

// 	c, w := utils.CreateTestContext(http.MethodGet, "/songs/1", nil)
// 	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

// 	handler.GetSongByIDHandler(c)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Contains(t, w.Body.String(), "Song")
// 	mockService.AssertExpectations(t)
// }

// func TestGetSongByIDHandler_NotFound(t *testing.T) {
// 	handler, mockService := setupSongHandlerTest()
// 	mockService.On("GetSongByID", "999").Return((*models.Song)(nil), errors.ErrResourceNotFound)

// 	c, w := utils.CreateTestContext(http.MethodGet, "/songs/999", nil)
// 	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "999"})

// 	handler.GetSongByIDHandler(c)

// 	assert.Equal(t, http.StatusNotFound, w.Code)
// 	mockService.AssertExpectations(t)
// }

// func TestGetSongByIDHandler_MissingID(t *testing.T) {
// 	handler, _ := setupSongHandlerTest()

// 	c, w := utils.CreateTestContext(http.MethodGet, "/songs/", nil)
// 	handler.GetSongByIDHandler(c)

//		assert.Equal(t, http.StatusBadRequest, w.Code)
//	}
func TestUpdateSongHandler_Success(t *testing.T) {
	handler, mockService := setupSongHandlerTest()

	title := "New Title"
	update := dto.UpdateSongRequest{
		Title: &title,
	}

	mockService.On("UpdateSong", "1", update).Return(nil)

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/1", strings.NewReader(`{"title":"New Title"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateSongHandler_MissingID(t *testing.T) {
	handler, _ := setupSongHandlerTest()

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/", strings.NewReader(`{"title": "New Title"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSongHandler_InvalidJSON(t *testing.T) {
	handler, mockService := setupSongHandlerTest()
	mockService.On("UpdateSong", "1", mock.Anything).Return(errors.ErrValidationFailed)

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/1", strings.NewReader(`{"title": 123}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSongHandler_EmptyUpdate(t *testing.T) {
	handler, mockService := setupSongHandlerTest()
	mockService.On("UpdateSong", "1", mock.Anything).Return(errors.ErrValidationFailed)

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/1", strings.NewReader(`{}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSongHandler_InvalidFields(t *testing.T) {
	handler, mockService := setupSongHandlerTest()
	mockService.On("UpdateSong", "1", mock.Anything).Return(errors.ErrValidationFailed)
	c, w := utils.CreateTestContext(http.MethodPut, "/songs/1", strings.NewReader(`{"title":""}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSongHandler_InvalidJSONBinding(t *testing.T) {
	handler, _ := setupSongHandlerTest()

	invalidJSON := `{"title":`

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/1", strings.NewReader(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSongHandler_ServiceError(t *testing.T) {
	handler, mockService := setupSongHandlerTest()

	title := "New Title"
	update := dto.UpdateSongRequest{
		Title: &title,
	}
	mockService.On("UpdateSong", "1", update).Return(errors.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodPut, "/songs/1", strings.NewReader(`{"title":"New Title"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.UpdateSongHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteSongWithDocumentsHandler_Success(t *testing.T) {
	handler, mockService := setupSongHandlerTest()
	mockService.On("DeleteSongWithDocuments", "1").Return(nil)

	c, w := utils.CreateTestContext(http.MethodDelete, "/songs/1", nil)
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "1"})

	handler.DeleteSongWithDocumentsHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteSongWithDocumentsHandler_Service_Error(t *testing.T) {
	handler, mockService := setupSongHandlerTest()
	mockService.On("DeleteSongWithDocuments", "999").Return(errors.ErrInternalServer)

	c, w := utils.CreateTestContext(http.MethodDelete, "/songs/999", nil)
	c.Params = append(c.Params, gin.Param{Key: "song_id", Value: "999"})

	handler.DeleteSongWithDocumentsHandler(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteSongWithDocumentsHandler_MissingID(t *testing.T) {
	handler, _ := setupSongHandlerTest()

	c, w := utils.CreateTestContext(http.MethodDelete, "/songs/", nil)
	handler.DeleteSongWithDocumentsHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
