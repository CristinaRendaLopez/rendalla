package errors

import "errors"

var (
	// Auth
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrTokenGenerationFailed = errors.New("failed to generate authentication token")
	ErrUnauthorized          = errors.New("unauthorized access")

	// Validation & requests
	ErrBadRequest       = errors.New("bad request")
	ErrValidationFailed = errors.New("validation failed")
	ErrTooManyResults   = errors.New("too many results")

	// Domain
	ErrResourceNotFound    = errors.New("resource not found")
	ErrOperationNotAllowed = errors.New("operation not allowed")

	// System
	ErrThroughputExceeded = errors.New("throughput limit exceeded")
	ErrInternalServer     = errors.New("internal server error")
)
