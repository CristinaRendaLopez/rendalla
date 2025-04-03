package integration_tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SongListItem struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type SongListResponse struct {
	Data []SongListItem `json:"data"`
}

type SongDetail struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Author string   `json:"author"`
	Genres []string `json:"genres"`
}

type SongDetailResponse struct {
	Data SongDetail `json:"data"`
}

func TestGETSongs_ShouldReturnSeededQueenSong(t *testing.T) {
	w := MakeRequest("GET", "/songs", nil, "")
	assert.Equal(t, http.StatusOK, w.Code)

	var response SongListResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	found := false
	for _, song := range response.Data {
		if song.Title == "Bohemian Rhapsody" && song.Author == "Queen" {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestGetSongByID_ShouldReturnQueenSong(t *testing.T) {
	w := MakeRequest("GET", "/songs/queen-001", nil, "")
	assert.Equal(t, http.StatusOK, w.Code)

	var response SongDetailResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "queen-001", response.Data.ID)
	assert.Equal(t, "Bohemian Rhapsody", response.Data.Title)
	assert.Equal(t, "Queen", response.Data.Author)
	assert.Contains(t, response.Data.Genres, "Rock")
}

func TestGetSongByID_ShouldReturn404IfNotExists(t *testing.T) {
	w := MakeRequest("GET", "/songs/non-existent-id", nil, "")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCreateSong_ShouldSucceedWithJWT(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{
		"title": "Imagine",
		"author": "John Lennon",
		"genres": ["Pop"],
		"documents": [
			{
				"type": "partitura",
				"instrument": ["piano"],
				"pdf_url": "https://test-bucket/imagine.pdf"
			}
		]
	}`

	w := MakeRequest("POST", "/songs", strings.NewReader(payload), token)
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Message string `json:"message"`
		SongID  string `json:"song_id"`
	}
	fmt.Println("BODY:", w.Body.String())
	err = json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Song created successfully", response.Message)
	assert.NotEmpty(t, response.SongID)
}

func TestCreateSong_ShouldFailWithoutJWT(t *testing.T) {
	payload := `{
		"title": "No Auth Song",
		"author": "Anon",
		"genres": ["Jazz"],
		"documents": [
			{
				"type": "partitura",
				"instrument": ["saxophone"],
				"pdf_url": "https://test-bucket/noauth.pdf"
			}
		]
	}`

	w := MakeRequest("POST", "/songs", strings.NewReader(payload), "")
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateSong_ShouldReturn400ForInvalidJSON(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{"title": "Oops"`

	w := MakeRequest("POST", "/songs", strings.NewReader(payload), token)
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateSong_ShouldReturn400ForInvalidFields(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{
		"title": "",
		"author": "",
		"genres": [],
		"documents": []
	}`

	w := MakeRequest("POST", "/songs", strings.NewReader(payload), token)
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSong_ShouldSucceed(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{
		"title": "Bohemian Rhapsody (Remastered)"
	}`

	w := MakeRequest("PUT", "/songs/queen-001", strings.NewReader(payload), token)
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateSong_ShouldReturn404(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{"title": "Ghost Song"}`

	w := MakeRequest("PUT", "/songs/nonexistent-id", strings.NewReader(payload), token)
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSong_ShouldReturn400ForInvalidJSON(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{"title":`

	w := MakeRequest("PUT", "/songs/queen-001", strings.NewReader(payload), token)
	w.Header().Set("Content-Type", "application/json")

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteSong_ShouldSucceed(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{
		"title": "Temporary Song",
		"author": "Test Artist",
		"genres": ["Test"],
		"documents": []
	}`

	createRes := MakeRequest("POST", "/songs", strings.NewReader(payload), token)
	createRes.Header().Set("Content-Type", "application/json")

	var createBody struct {
		SongID string `json:"song_id"`
	}
	err = json.NewDecoder(createRes.Body).Decode(&createBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, createBody.SongID)

	deleteRes := MakeRequest("DELETE", "/songs/"+createBody.SongID, nil, token)
	assert.Equal(t, http.StatusOK, deleteRes.Code)
}

func TestDeleteSong_ShouldAlsoDeleteDocuments(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	payload := `{
		"title": "Song to Delete",
		"author": "Doc Killer",
		"genres": ["Jazz"],
		"documents": [
			{
				"type": "tablatura",
				"instrument": ["bass"],
				"pdf_url": "https://test.com/delete_me.pdf"
			}
		]
	}`

	createRes := MakeRequest("POST", "/songs", strings.NewReader(payload), token)
	createRes.Header().Set("Content-Type", "application/json")

	var createBody struct {
		SongID string `json:"song_id"`
	}
	err = json.NewDecoder(createRes.Body).Decode(&createBody)
	assert.NoError(t, err)
	assert.NotEmpty(t, createBody.SongID)

	getDocsRes := MakeRequest("GET", "/songs/"+createBody.SongID+"/documents", nil, "")
	assert.Equal(t, http.StatusOK, getDocsRes.Code)

	var docsResponse struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	err = json.NewDecoder(getDocsRes.Body).Decode(&docsResponse)
	assert.NoError(t, err)
	assert.NotEmpty(t, docsResponse.Data)

	docID := docsResponse.Data[0].ID

	deleteRes := MakeRequest("DELETE", "/songs/"+createBody.SongID, nil, token)
	assert.Equal(t, http.StatusOK, deleteRes.Code)

	docGetRes := MakeRequest("GET", "/documents/"+docID, nil, "")
	assert.Equal(t, http.StatusNotFound, docGetRes.Code)
}

func TestDeleteSong_ShouldReturn404IfNotExists(t *testing.T) {
	token, err := GenerateTestJWT("admin")
	assert.NoError(t, err)

	res := MakeRequest("DELETE", "/songs/non-existent-id", nil, token)
	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestDeleteSong_ShouldReturn401WithoutToken(t *testing.T) {
	res := MakeRequest("DELETE", "/songs/any-id", nil, "")
	assert.Equal(t, http.StatusUnauthorized, res.Code)
}
