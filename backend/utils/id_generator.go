package utils

import (
	"github.com/google/uuid"
)

// IDGenerator defines an interface for generating unique IDs.
// Useful for injecting deterministic or mockable ID generators in tests.
type IDGenerator interface {
	NewID() string
}

// UUIDGenerator implements IDGenerator using UUID v4.
type UUIDGenerator struct{}

// NewID returns a new random UUID as a string.
func (g *UUIDGenerator) NewID() string {
	return uuid.New().String()
}
