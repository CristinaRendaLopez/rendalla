package errors

import "errors"

var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrTokenGenerationFailed = errors.New("failed to generate authentication token")
	ErrBadRequest            = errors.New("bad request")
	ErrValidationFailed      = errors.New("validation failed")
	ErrResourceNotFound      = errors.New("resource not found")
	ErrOperationNotAllowed   = errors.New("operation not allowed")
	ErrThroughputExceeded    = errors.New("throughput limit exceeded")
	ErrInternalServer        = errors.New("internal server error")
	ErrUnauthorized          = errors.New("unauthorized access")
)
