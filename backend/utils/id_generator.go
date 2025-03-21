package utils

import (
	"github.com/google/uuid"
)

type IDGenerator interface {
	NewID() string
}

type UUIDGenerator struct{}

func (g *UUIDGenerator) NewID() string {
	return uuid.New().String()
}
