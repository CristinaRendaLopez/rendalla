package integration_tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/stretchr/testify/suite"
)

type SearchTestSuite struct {
	IntegrationTestSuite
}

func (s *SearchTestSuite) SetupTest() {
	err := ClearTestTables(s.DB)
	s.Require().NoError(err)

	err = SeedTestData(s.DB, s.TimeProvider)
	s.Require().NoError(err)
}

func (s *SearchTestSuite) TestSearchSongs_ShouldReturnAll() {
	res := MakeRequest(s.Router, "GET", "/songs/search", nil, "")
	s.Equal(http.StatusOK, res.Code)

	var body struct {
		Data []dto.SongResponseItem `json:"data"`
	}
	err := json.NewDecoder(res.Body).Decode(&body)
	s.Require().NoError(err)
	s.Len(body.Data, 2)
}

func (s *SearchTestSuite) TestSearchSongs_ByTitle() {
	res := MakeRequest(s.Router, "GET", "/songs/search?title=Bohemian", nil, "")
	s.Equal(http.StatusOK, res.Code)

	var body struct {
		Data []dto.SongResponseItem `json:"data"`
	}
	err := json.NewDecoder(res.Body).Decode(&body)
	s.Require().NoError(err)

	s.Len(body.Data, 1)
	s.Equal("Bohemian Rhapsody", body.Data[0].Title)
}

func (s *SearchTestSuite) TestSearchSongs_NoMatch() {
	res := MakeRequest(s.Router, "GET", "/songs/search?title=NothingHere", nil, "")
	s.Equal(http.StatusOK, res.Code)

	var body struct {
		Data []dto.SongResponseItem `json:"data"`
	}
	err := json.NewDecoder(res.Body).Decode(&body)
	s.Require().NoError(err)

	s.Len(body.Data, 0)
}

func (s *SearchTestSuite) TestSearchDocuments_ShouldReturnAll() {
	res := MakeRequest(s.Router, "GET", "/documents/search", nil, "")
	s.Equal(http.StatusOK, res.Code)

	var body struct {
		Data []dto.DocumentResponseItem `json:"data"`
	}
	err := json.NewDecoder(res.Body).Decode(&body)
	s.Require().NoError(err)
	s.Len(body.Data, 3)
}

func (s *SearchTestSuite) TestSearchDocuments_ByType() {
	res := MakeRequest(s.Router, "GET", "/documents/search?type=score", nil, "")
	s.Equal(http.StatusOK, res.Code)

	var body struct {
		Data []dto.DocumentResponseItem `json:"data"`
	}
	err := json.NewDecoder(res.Body).Decode(&body)
	s.Require().NoError(err)

	s.Len(body.Data, 2)
	for _, doc := range body.Data {
		s.Equal("score", doc.Type)
	}
}

func (s *SearchTestSuite) TestSearchDocuments_ByInstrument() {
	res := MakeRequest(s.Router, "GET", "/documents/search?instrument=guitar", nil, "")
	s.Equal(http.StatusOK, res.Code)

	var body struct {
		Data []dto.DocumentResponseItem `json:"data"`
	}
	err := json.NewDecoder(res.Body).Decode(&body)
	s.Require().NoError(err)

	s.Len(body.Data, 1)
	s.Contains(body.Data[0].Instrument, "guitar")
}

func (s *SearchTestSuite) TestSearchDocuments_CombinedFilters() {
	res := MakeRequest(s.Router, "GET", "/documents/search?instrument=piano&type=score", nil, "")
	s.Equal(http.StatusOK, res.Code)

	var body struct {
		Data []dto.DocumentResponseItem `json:"data"`
	}
	err := json.NewDecoder(res.Body).Decode(&body)
	s.Require().NoError(err)

	s.Len(body.Data, 1)
	s.Equal("doc-br-piano", body.Data[0].ID)
}

func (s *SearchTestSuite) TestSearchDocuments_NoMatch() {
	res := MakeRequest(s.Router, "GET", "/documents/search?instrument=harp", nil, "")
	s.Equal(http.StatusOK, res.Code)

	var body struct {
		Data []dto.DocumentResponseItem `json:"data"`
	}
	err := json.NewDecoder(res.Body).Decode(&body)
	s.Require().NoError(err)

	s.Len(body.Data, 0)
}

func TestSearchSuite(t *testing.T) {
	suite.Run(t, new(SearchTestSuite))
}
