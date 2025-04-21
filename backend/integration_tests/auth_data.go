package integration_tests

import "github.com/CristinaRendaLopez/rendalla-backend/dto"

var ValidLogin = dto.LoginRequest{
	Username: "test",
	Password: "test",
}

var InvalidUsernameLogin = dto.LoginRequest{
	Username: "wronguser",
	Password: "test",
}

var InvalidPasswordLogin = dto.LoginRequest{
	Username: "test",
	Password: "wrongpass",
}

var InvalidJSONLogin = `{"username":`
