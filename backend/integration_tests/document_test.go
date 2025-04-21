package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/stretchr/testify/suite"
)

type DocumentTestSuite struct {
	IntegrationTestSuite
}

func (s *DocumentTestSuite) SetupTest() {
	err := ClearTestTables(s.DB)
	s.Require().NoError(err)

	err = SeedTestData(s.DB, s.TimeProvider)
	s.Require().NoError(err)
}

func (s *DocumentTestSuite) TestGetDocumentsBySongID_ShouldReturnSeededDocuments() {
	w := MakeRequest(s.Router, "GET", "/songs/queen-001/documents", nil, "")
	s.Equal(http.StatusOK, w.Code)

	var response DocumentListResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	s.Require().NoError(err)
	s.Len(response.Data, 2)
}

func (s *DocumentTestSuite) TestGetDocumentsBySongID_ShouldReturnEmptyList() {
	w := MakeRequest(s.Router, "GET", "/songs/non-existent-id/documents", nil, "")
	s.Equal(http.StatusOK, w.Code)

	var response DocumentListResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	s.Require().NoError(err)
	s.Empty(response.Data)
}

func (s *DocumentTestSuite) TestGetDocumentByID_ShouldReturnSeededDocument() {
	w := MakeRequest(s.Router, "GET", "/songs/queen-001/documents/doc-br-piano", nil, "")
	s.Equal(http.StatusOK, w.Code)

	var response DocumentDetailResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	s.Require().NoError(err)
	s.Equal("doc-br-piano", response.Data.ID)
}

func (s *DocumentTestSuite) TestGetDocumentByID_ShouldReturn404() {
	w := MakeRequest(s.Router, "GET", "/songs/queen-001/documents/non-existent-id", nil, "")
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *DocumentTestSuite) TestCreateDocument_ShouldSucceedWithJWT() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	body, err := json.Marshal(ViolinScore)
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "POST", "/songs/queen-001/documents", bytes.NewReader(body), token)
	s.Equal(http.StatusCreated, w.Code)

	var res struct {
		Message    string `json:"message"`
		DocumentID string `json:"document_id"`
	}
	err = json.NewDecoder(w.Body).Decode(&res)
	s.Require().NoError(err)
	s.Equal("Document created successfully", res.Message)
	s.NotEmpty(res.DocumentID)
}

func (s *DocumentTestSuite) TestCreateDocument_ShouldReturn401WithoutToken() {

	body, err := json.Marshal(FluteScore)
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "POST", "/songs/queen-001/documents", bytes.NewReader(body), "")
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *DocumentTestSuite) TestCreateDocument_ShouldReturn400ForInvalidJSON() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	body, err := json.Marshal(InvalidJSONDocument)
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "POST", "/songs/queen-001/documents", bytes.NewReader(body), token)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *DocumentTestSuite) TestCreateDocument_ShouldReturn400ForInvalidFields() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	body, err := json.Marshal(InvalidFieldsDocument)
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "POST", "/songs/queen-001/documents", bytes.NewReader(body), token)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *DocumentTestSuite) TestUpdateDocument_ShouldSucceed() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	body, err := json.Marshal(TablatureUpdate)
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "PUT", "/songs/queen-001/documents/doc-br-piano", bytes.NewReader(body), token)
	s.Equal(http.StatusOK, w.Code)

	getRes := MakeRequest(s.Router, "GET", "/songs/queen-001/documents/doc-br-piano", nil, "")
	s.Equal(http.StatusOK, getRes.Code)

	var getBody struct {
		Data dto.DocumentResponseItem `json:"data"`
	}
	err = json.NewDecoder(getRes.Body).Decode(&getBody)
	s.Require().NoError(err)
	s.Equal(TablatureUpdate.Type, getBody.Data.Type)
	s.ElementsMatch(TablatureUpdate.Instrument, getBody.Data.Instrument)
	s.Equal(TablatureUpdate.PDFURL, getBody.Data.PDFURL)
	s.Equal(TablatureUpdate.AudioURL, getBody.Data.AudioURL)
}

func (s *DocumentTestSuite) TestUpdateDocument_ShouldReturn400ForInvalidJSON() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	body, err := json.Marshal(InvalidJSONDocument)
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "PUT", "/songs/queen-001/documents/doc-br-voice", bytes.NewReader(body), token)
	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *DocumentTestSuite) TestUpdateDocument_ShouldReturn404IfNotExists() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	body, err := json.Marshal(TablatureUpdate)
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "PUT", "/songs/queen-001/documents/nonexistent-id", bytes.NewReader(body), token)
	s.Equal(http.StatusNotFound, w.Code)
}

func (s *DocumentTestSuite) TestDeleteDocument_ShouldSucceed() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	body, err := json.Marshal(FluteScore)
	s.Require().NoError(err)

	createRes := MakeRequest(s.Router, "POST", "/songs/queen-001/documents", bytes.NewReader(body), token)
	var createBody struct {
		DocumentID string `json:"document_id"`
	}
	err = json.NewDecoder(createRes.Body).Decode(&createBody)
	s.Require().NoError(err)
	s.NotEmpty(createBody.DocumentID)

	deleteRes := MakeRequest(s.Router, "DELETE", "/songs/queen-001/documents/"+createBody.DocumentID, nil, token)
	s.Equal(http.StatusOK, deleteRes.Code)
}

func (s *DocumentTestSuite) TestDeleteDocument_ShouldReturn401WithoutToken() {
	w := MakeRequest(s.Router, "DELETE", "/songs/queen-001/documents/doc-br-voice", nil, "")
	s.Equal(http.StatusUnauthorized, w.Code)
}

func (s *DocumentTestSuite) TestDeleteDocument_ShouldReturn404IfNotExists() {
	token, err := GenerateTestJWT("admin")
	s.Require().NoError(err)

	w := MakeRequest(s.Router, "DELETE", "/songs/queen-001/documents/nonexistent-id", nil, token)
	s.Equal(http.StatusNotFound, w.Code)
}

func TestDocumentSuite(t *testing.T) {
	suite.Run(t, new(DocumentTestSuite))
}
