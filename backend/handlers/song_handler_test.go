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

func TestGetAllSongsHandler(t *testing.T) {
	tests := []struct {
		name         string
		mockSongs    []dto.SongResponseItem
		mockError    error
		expectedCode int
		expectedData []dto.SongResponseItem
	}{
		{
			name:         "success with songs",
			mockSongs:    []dto.SongResponseItem{{ID: "1", Title: "Don't Stop Me Now", Author: "Queen", Genres: []string{"rock"}}},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedData: []dto.SongResponseItem{{ID: "1", Title: "Don't Stop Me Now", Author: "Queen", Genres: []string{"rock"}}},
		},
		{
			name:         "success with empty list",
			mockSongs:    []dto.SongResponseItem{},
			mockError:    nil,
			expectedCode: http.StatusOK,
			expectedData: []dto.SongResponseItem{},
		},
		{
			name:         "internal service error",
			mockSongs:    nil,
			mockError:    errors.ErrInternalServer,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupSongHandlerTest()

			mockService.
				On("GetAllSongs").
				Return(tt.mockSongs, tt.mockError)

			c, w := utils.CreateTestContext(http.MethodGet, "/songs", nil)
			handler.GetAllSongsHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				response := struct {
					Data []dto.SongResponseItem `json:"data"`
				}{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, response.Data)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetSongByIDHandler(t *testing.T) {
	tests := []struct {
		name         string
		songID       string
		setupParam   bool
		mockResult   dto.SongResponseItem
		mockError    error
		expectedCode int
	}{
		{
			name:       "success",
			songID:     "1",
			setupParam: true,
			mockResult: dto.SongResponseItem{
				ID:     "1",
				Title:  "Radio Ga Ga",
				Author: "Queen",
				Genres: []string{"pop", "rock"},
			},
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "missing song_id parameter",
			setupParam:   false,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "song not found",
			songID:       "999",
			setupParam:   true,
			mockError:    errors.ErrResourceNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "internal service error",
			songID:       "1",
			setupParam:   true,
			mockError:    errors.ErrInternalServer,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupSongHandlerTest()

			if tt.setupParam {
				mockService.
					On("GetSongByID", tt.songID).
					Return(tt.mockResult, tt.mockError)
			}

			reqPath := "/songs"
			if tt.setupParam {
				reqPath = "/songs/" + tt.songID
			}

			c, w := utils.CreateTestContext(http.MethodGet, reqPath, nil)
			if tt.setupParam {
				c.Params = []gin.Param{{Key: "song_id", Value: tt.songID}}
			}

			handler.GetSongByIDHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response struct {
					Data dto.SongResponseItem `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResult, response.Data)
			}

			if tt.setupParam {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestUpdateSongHandler(t *testing.T) {
	tests := []struct {
		name         string
		songID       string
		setupParam   bool
		body         string
		expectedCode int
		mockError    error
	}{
		{
			name:       "successfully updates song title and genres",
			songID:     "1",
			setupParam: true,
			body: `{
				"title": "We Are The Champions",
				"genres": ["rock", "anthem"]
			}`,
			expectedCode: http.StatusOK,
			mockError:    nil,
		},
		{
			name:         "missing song_id param",
			setupParam:   false,
			body:         `{ "title": "Somebody to Love" }`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid JSON payload",
			songID:       "1",
			setupParam:   true,
			body:         `{ "title": "Love of My Life",`, // missing bracket
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "empty update payload",
			songID:       "1",
			setupParam:   true,
			body:         `{}`,
			expectedCode: http.StatusBadRequest,
			mockError:    errors.ErrValidationFailed,
		},
		{
			name:         "invalid title too short",
			songID:       "1",
			setupParam:   true,
			body:         `{ "title": "a" }`,
			expectedCode: http.StatusBadRequest,
			mockError:    errors.ErrValidationFailed,
		},
		{
			name:         "song not found",
			songID:       "999",
			setupParam:   true,
			body:         `{ "author": "Queen" }`,
			expectedCode: http.StatusNotFound,
			mockError:    errors.ErrResourceNotFound,
		},
		{
			name:         "internal server error",
			songID:       "1",
			setupParam:   true,
			body:         `{ "genres": ["rock"] }`,
			expectedCode: http.StatusInternalServerError,
			mockError:    errors.ErrInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupSongHandlerTest()

			if tt.setupParam && tt.mockError != errors.ErrValidationFailed {
				var update dto.UpdateSongRequest
				_ = json.Unmarshal([]byte(tt.body), &update)

				mockService.
					On("UpdateSong", tt.songID, update).
					Return(tt.mockError)
			}

			path := "/songs"
			if tt.setupParam {
				path += "/" + tt.songID
			}

			c, w := utils.CreateTestContext(http.MethodPut, path, strings.NewReader(tt.body))
			if tt.setupParam {
				c.Params = []gin.Param{{Key: "song_id", Value: tt.songID}}
			}

			handler.UpdateSongHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.setupParam && tt.expectedCode != http.StatusBadRequest {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestDeleteSongWithDocumentsHandler(t *testing.T) {
	tests := []struct {
		name         string
		songID       string
		setupParam   bool
		mockError    error
		expectedCode int
	}{
		{
			name:         "successfully deletes song",
			songID:       "1",
			setupParam:   true,
			mockError:    nil,
			expectedCode: http.StatusOK,
		},
		{
			name:         "missing song_id param",
			setupParam:   false,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "song not found",
			songID:       "999",
			setupParam:   true,
			mockError:    errors.ErrResourceNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "internal server error",
			songID:       "1",
			setupParam:   true,
			mockError:    errors.ErrInternalServer,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupSongHandlerTest()

			if tt.setupParam && tt.expectedCode != http.StatusBadRequest {
				mockService.
					On("DeleteSongWithDocuments", tt.songID).
					Return(tt.mockError)
			}

			path := "/songs"
			if tt.setupParam {
				path += "/" + tt.songID
			}

			c, w := utils.CreateTestContext(http.MethodDelete, path, nil)
			if tt.setupParam {
				c.Params = []gin.Param{{Key: "song_id", Value: tt.songID}}
			}

			handler.DeleteSongWithDocumentsHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.setupParam && tt.expectedCode != http.StatusBadRequest {
				mockService.AssertExpectations(t)
			}
		})
	}
}
