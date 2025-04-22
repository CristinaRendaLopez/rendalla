package dto

import "github.com/CristinaRendaLopez/rendalla-backend/repository"

func ToAuthCredentials(dto LoginRequest) repository.AuthCredentials {
	return repository.AuthCredentials{
		Username: dto.Username,
		Password: dto.Password,
	}
}
