package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/CristinaRendaLopez/rendalla-backend/dto"
	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	IntegrationTestSuite
}

func (s *AuthTestSuite) SetupTest() {}

func (s *AuthTestSuite) TestLogin_ShouldSucceed() {
	body, err := json.Marshal(ValidLogin)
	s.Require().NoError(err)

	res := MakeRequest(s.Router, "POST", "/auth/login", bytes.NewReader(body), "")
	s.Equal(http.StatusOK, res.Code)

	var resp dto.AuthResponse
	err = json.NewDecoder(res.Body).Decode(&resp)
	s.Require().NoError(err)
	s.NotEmpty(resp.Token)
}

func (s *AuthTestSuite) TestLogin_ShouldFail_InvalidUsername() {
	body, err := json.Marshal(InvalidUsernameLogin)
	s.Require().NoError(err)

	res := MakeRequest(s.Router, "POST", "/auth/login", bytes.NewReader(body), "")
	s.Equal(http.StatusUnauthorized, res.Code)
}

func (s *AuthTestSuite) TestLogin_ShouldFail_InvalidPassword() {
	body, err := json.Marshal(InvalidPasswordLogin)
	s.Require().NoError(err)

	res := MakeRequest(s.Router, "POST", "/auth/login", bytes.NewReader(body), "")
	s.Equal(http.StatusUnauthorized, res.Code)
}

func (s *AuthTestSuite) TestLogin_ShouldFail_InvalidPayload() {
	res := MakeRequest(s.Router, "POST", "/auth/login", strings.NewReader(InvalidJSONLogin), "")
	s.Equal(http.StatusBadRequest, res.Code)
}

func (s *AuthTestSuite) TestMe_ShouldReturnUser_WithValidToken() {
	username := ValidLogin.Username

	token, err := GenerateTestJWT(username)
	s.Require().NoError(err)

	res := MakeRequest(s.Router, "GET", "/auth/me", nil, token)
	s.Equal(http.StatusOK, res.Code)

	var me dto.MeResponse
	err = json.NewDecoder(res.Body).Decode(&me)
	s.Require().NoError(err)
	s.Equal(username, me.Username)
	s.Equal("admin", me.Role)
}

func (s *AuthTestSuite) TestMe_ShouldReturn401_WithoutToken() {
	res := MakeRequest(s.Router, "GET", "/auth/me", nil, "")
	s.Equal(http.StatusUnauthorized, res.Code)
}

func (s *AuthTestSuite) TestMe_ShouldReturn401_WithInvalidToken() {
	invalidToken := "this.is.not.valid"

	res := MakeRequest(s.Router, "GET", "/auth/me", nil, invalidToken)
	s.Equal(http.StatusUnauthorized, res.Code)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
