package errors

import (
	"errors"
	"net/http"
)

func MapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, ErrBadRequest), errors.Is(err, ErrValidationFailed):
		return http.StatusBadRequest
	case errors.Is(err, ErrResourceNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrOperationNotAllowed):
		return http.StatusForbidden
	case errors.Is(err, ErrThroughputExceeded):
		return http.StatusTooManyRequests
	case errors.Is(err, ErrInvalidCredentials), errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func MapErrorToMessage(err error) string {
	switch {
	case errors.Is(err, ErrBadRequest):
		return "The request was malformed or missing required parameters"
	case errors.Is(err, ErrValidationFailed):
		return "Some of the provided data is invalid"
	case errors.Is(err, ErrResourceNotFound):
		return "The requested resource was not found"
	case errors.Is(err, ErrOperationNotAllowed):
		return "You are not allowed to perform this operation"
	case errors.Is(err, ErrThroughputExceeded):
		return "Too many requests, please try again later"
	case errors.Is(err, ErrInvalidCredentials):
		return "Invalid username or password"
	case errors.Is(err, ErrUnauthorized):
		return "You must be authenticated to perform this action"
	case errors.Is(err, ErrTokenGenerationFailed):
		return "Could not generate access token"
	case errors.Is(err, ErrInternalServer):
		return "An internal server error occurred"
	default:
		return "An unexpected error occurred"
	}
}

func MapErrorToCode(err error) string {
	switch {
	case errors.Is(err, ErrBadRequest):
		return "bad_request"
	case errors.Is(err, ErrValidationFailed):
		return "validation_failed"
	case errors.Is(err, ErrResourceNotFound):
		return "resource_not_found"
	case errors.Is(err, ErrOperationNotAllowed):
		return "operation_not_allowed"
	case errors.Is(err, ErrThroughputExceeded):
		return "throughput_exceeded"
	case errors.Is(err, ErrInvalidCredentials):
		return "invalid_credentials"
	case errors.Is(err, ErrUnauthorized):
		return "unauthorized"
	case errors.Is(err, ErrTokenGenerationFailed):
		return "token_generation_failed"
	case errors.Is(err, ErrInternalServer):
		return "internal_server_error"
	default:
		return "unknown_error"
	}
}
