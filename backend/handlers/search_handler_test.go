package handlers_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/errors"
	"github.com/CristinaRendaLopez/rendalla-backend/handlers"
	"github.com/CristinaRendaLopez/rendalla-backend/mocks"
	"github.com/CristinaRendaLopez/rendalla-backend/models"
	"github.com/CristinaRendaLopez/rendalla-backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSearchHandlerTest() (*handlers.SearchHandler, *mocks.MockSearchService) {
	mockService := new(mocks.MockSearchService)
	handler := handlers.NewSearchHandler(mockService)
	return handler, mockService
}

func TestListSongsHandler(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		mockTitle    string
		mockSort     string
		mockOrder    string
		mockReturn   []models.Song
		mockNext     interface{}
		mockErr      error
		expectedCode int
		expectedBody []string
	}{
		{
			name:         "filter by title",
			query:        "title=love",
			mockTitle:    "love",
			mockReturn:   []models.Song{SongLoveOfMyLife},
			expectedCode: http.StatusOK,
			expectedBody: []string{"Love of My Life"},
		},
		{
			name:         "sort by title desc",
			query:        "title=love&sort=title&order=desc",
			mockTitle:    "love",
			mockSort:     "title",
			mockOrder:    "desc",
			mockReturn:   []models.Song{SongSomebodyToLove, SongLoveOfMyLife},
			expectedCode: http.StatusOK,
			expectedBody: []string{"Somebody to Love", "Love of My Life"},
		},
		{
			name:         "empty result",
			query:        "title=nothing",
			mockTitle:    "nothing",
			mockReturn:   []models.Song{},
			expectedCode: http.StatusOK,
			expectedBody: []string{`"data":[]`},
		},
		{
			name:         "next_token included",
			query:        "next_token=abc",
			mockNext:     "abc",
			mockReturn:   []models.Song{SongOneVision},
			expectedCode: http.StatusOK,
			expectedBody: []string{"One Vision", "next_token"},
		},
		{
			name:         "service error",
			query:        "title=error",
			mockTitle:    "error",
			mockErr:      errors.ErrInternalServer,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupSearchHandlerTest()

			mockService.On("ListSongs", tt.mockTitle, tt.mockSort, tt.mockOrder, 10, mock.Anything).
				Return(tt.mockReturn, tt.mockNext, tt.mockErr)

			path := "/songs/search"
			if tt.query != "" {
				path += "?" + tt.query
			}

			c, w := utils.CreateTestContext(http.MethodGet, path, nil)
			c.Request.URL.RawQuery = tt.query

			handler.ListSongsHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)
			for _, s := range tt.expectedBody {
				assert.Contains(t, w.Body.String(), s)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestListDocumentsHandler(t *testing.T) {
	tests := []struct {
		name         string
		query        string
		params       []string
		mockReturn   []models.Document
		mockNext     interface{}
		mockErr      error
		expectedCode int
		expectedIDs  []string
	}{
		{
			name:         "filter by title",
			query:        "title=queen",
			params:       []string{"queen", "", "", "", ""},
			mockReturn:   []models.Document{DocSheetMusicGuitar},
			expectedCode: http.StatusOK,
			expectedIDs:  []string{"1"},
		},
		{
			name:         "filter by instrument",
			query:        "instrument=Piano",
			params:       []string{"", "Piano", "", "", ""},
			mockReturn:   []models.Document{DocSheetMusicPiano},
			expectedCode: http.StatusOK,
			expectedIDs:  []string{"2"},
		},
		{
			name:         "filter by type",
			query:        "type=tablature",
			params:       []string{"", "", "tablature", "", ""},
			mockReturn:   []models.Document{DocTablatureGuitar},
			expectedCode: http.StatusOK,
			expectedIDs:  []string{"3"},
		},
		{
			name:         "combined filters with sorting",
			query:        "title=love&instrument=Violin&type=sheet_music&sort=title&order=asc",
			params:       []string{"love", "Violin", "sheet_music", "title", "asc"},
			mockReturn:   []models.Document{DocViolinLoveOfMyLife, DocViolinSomebodyToLove},
			expectedCode: http.StatusOK,
			expectedIDs:  []string{"3", "4"},
		},
		{
			name:         "sort by created_at desc",
			query:        "sort=created_at&order=desc",
			params:       []string{"", "", "", "created_at", "desc"},
			mockReturn:   []models.Document{DocUnderPressure, DocInnuendo},
			expectedCode: http.StatusOK,
			expectedIDs:  []string{"5", "6"},
		},
		{
			name:         "empty result",
			query:        "title=none",
			params:       []string{"none", "", "", "", ""},
			mockReturn:   []models.Document{},
			expectedCode: http.StatusOK,
			expectedIDs:  []string{},
		},
		{
			name:         "service error",
			query:        "title=queen",
			params:       []string{"queen", "", "", "", ""},
			mockErr:      errors.ErrInternalServer,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, mockService := setupSearchHandlerTest()

			mockService.On("ListDocuments",
				tt.params[0], tt.params[1], tt.params[2], tt.params[3], tt.params[4], 10, mock.Anything,
			).Return(tt.mockReturn, tt.mockNext, tt.mockErr)

			path := "/documents/search"
			if tt.query != "" {
				path += "?" + tt.query
			}

			c, w := utils.CreateTestContext(http.MethodGet, path, nil)
			c.Request.URL.RawQuery = tt.query

			handler.ListDocumentsHandler(c)

			assert.Equal(t, tt.expectedCode, w.Code)

			if tt.expectedCode == http.StatusOK {
				var response struct {
					Data []models.Document `json:"data"`
				}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				var resultIDs []string
				for _, doc := range response.Data {
					resultIDs = append(resultIDs, doc.ID)
				}
				assert.ElementsMatch(t, tt.expectedIDs, resultIDs)
			}

			mockService.AssertExpectations(t)
		})
	}
}
