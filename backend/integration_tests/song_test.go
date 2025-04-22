package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"slices"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/stretchr/testify/suite"
)

type SongTestSuite struct {
	IntegrationTestSuite
}

func (s *SongTestSuite) TestCreateSong_ShouldSucceed() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	body, err := json.Marshal(WeAreTheChampionsPayload)
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "POST", "/songs", bytes.NewReader(body), token)

	s.Equal(http.StatusCreated, w.Code)

	var res dto.CreateSongResponse
	err = json.NewDecoder(w.Body).Decode(&res)
	s.Require().NoError(err)
	s.Equal("Song created successfully", res.Message)
	s.NotEmpty(res.SongID)
}

func (s *SongTestSuite) TestCreateSong_ShouldReturnDocuments() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	body, err := json.Marshal(DontStopMeNowPayload)
	s.Require().NoError(err)

	createRes := MakeRequest(s.Router, "POST", "/songs", bytes.NewReader(body), token)

	s.Equal(http.StatusCreated, createRes.Code)

	var createBody dto.CreateSongResponse
	err = json.NewDecoder(createRes.Body).Decode(&createBody)
	s.Require().NoError(err)
	s.NotEmpty(createBody.SongID)

	getDocsRes := MakeRequest(s.Router, "GET", "/songs/"+createBody.SongID+"/documents", nil, "")
	s.Equal(http.StatusOK, getDocsRes.Code)

	var docsResponse struct {
		Data []dto.DocumentResponseItem `json:"data"`
	}
	err = json.NewDecoder(getDocsRes.Body).Decode(&docsResponse)
	s.Require().NoError(err)
	s.Len(docsResponse.Data, 2)

	var foundScore, foundTab bool
	for _, doc := range docsResponse.Data {
		if doc.Type == "score" && slices.Contains(doc.Instrument, "voice") {
			foundScore = true
		}
		if doc.Type == "tablature" && slices.Contains(doc.Instrument, "bass") {
			foundTab = true
		}
	}
	s.True(foundScore)
	s.True(foundTab)
}

func (s *SongTestSuite) TestCreateSong_ShouldFailWithoutJWT() {
	body, err := json.Marshal(BohemianRhapsodyPayload)
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "POST", "/songs", bytes.NewReader(body), "")
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *SongTestSuite) TestCreateSong_ShouldReturn400ForInvalidJSON() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)
	w := MakeRequest(s.Router, "POST", "/songs", bytes.NewReader([]byte(`{"title": "`)), token)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *SongTestSuite) TestCreateSong_ShouldReturn400ForInvalidFields() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	invalidPayload := dto.CreateSongRequest{}
	body, err := json.Marshal(invalidPayload)
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "POST", "/songs", bytes.NewReader(body), token)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *SongTestSuite) TestGetSongs_ShouldIncludeBohemian() {
	w := MakeRequest(s.Router, "GET", "/songs", nil, "")
	s.Equal(http.StatusOK, w.Code)

	var response songListResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	s.Require().NoError(err)

	found := false
	for _, song := range response.Data {
		if song.Title == "Bohemian Rhapsody" && song.Author == "Queen" {
			found = true
			break
		}
	}
	s.True(found)
}

func (s *SongTestSuite) TestGetSongByID_ShouldReturnBohemian() {
	w := MakeRequest(s.Router, "GET", "/songs/queen-001", nil, "")
	s.Equal(http.StatusOK, w.Code)

	var response songDetailResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	s.Require().NoError(err)

	s.Equal("queen-001", response.Data.ID)
	s.Equal("Bohemian Rhapsody", response.Data.Title)
	s.Equal("Queen", response.Data.Author)
	s.Contains(response.Data.Genres, "Rock")
}

func (s *SongTestSuite) TestGetSongByID_ShouldReturn404() {
	w := MakeRequest(s.Router, "GET", "/songs/non-existent-id", nil, "")
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *SongTestSuite) TestUpdateSong_ShouldSucceed() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	title := "Bohemian Rhapsody (Remastered)"
	author := "Freddie Mercury"
	genres := []string{"Rock", "Opera"}

	update := dto.UpdateSongRequest{
		Title:  &title,
		Author: &author,
		Genres: genres,
	}

	body, err := json.Marshal(update)
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "PUT", "/songs/queen-001", bytes.NewReader(body), token)
	s.Equal(http.StatusOK, w.Code)

	getRes := MakeRequest(s.Router, "GET", "/songs/queen-001", nil, "")
	var getBody songDetailResponse
	err = json.NewDecoder(getRes.Body).Decode(&getBody)
	s.Require().NoError(err)
	s.Equal(title, getBody.Data.Title)
	s.Equal(author, getBody.Data.Author)
	s.ElementsMatch(genres, getBody.Data.Genres)
}

func (s *SongTestSuite) TestUpdateSong_ShouldReturn404() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	title := "Ghost Song"
	update := dto.UpdateSongRequest{Title: &title}
	body, err := json.Marshal(update)
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "PUT", "/songs/nonexistent-id", bytes.NewReader(body), token)
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *SongTestSuite) TestUpdateSong_ShouldReturn400ForInvalidJSON() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)
	w := MakeRequest(s.Router, "PUT", "/songs/queen-001", bytes.NewReader([]byte(`{"title":`)), token)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *SongTestSuite) TestDeleteSong_ShouldSucceed() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	body, err := json.Marshal(DontStopMeNowPayload)
	s.Require().NoError(err)

	createRes := MakeRequest(s.Router, "POST", "/songs", bytes.NewReader(body), token)

	var createBody dto.CreateSongResponse
	err = json.NewDecoder(createRes.Body).Decode(&createBody)
	s.Require().NoError(err)
	s.NotEmpty(createBody.SongID)

	deleteRes := MakeRequest(s.Router, "DELETE", "/songs/"+createBody.SongID, nil, token)
	s.Equal(http.StatusOK, deleteRes.Code)
}

func (s *SongTestSuite) TestDeleteSong_ShouldReturn404() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)
	res := MakeRequest(s.Router, "DELETE", "/songs/non-existent-id", nil, token)
	s.Equal(http.StatusNotFound, res.Code)
}

func (s *SongTestSuite) TestDeleteSong_ShouldReturn401WithoutToken() {
	res := MakeRequest(s.Router, "DELETE", "/songs/any-id", nil, "")
	s.Equal(http.StatusUnauthorized, res.Code)
}

func TestSongSuite(t *testing.T) {
	suite.Run(t, new(SongTestSuite))
}

func (s *SongTestSuite) SetupTest() {
	err := ClearTestTables(s.DB)
	s.Require().NoError(err)

	err = SeedTestData(s.DB, s.TimeProvider)
	s.Require().NoError(err)
}
