package utils

import (
	"time"
)

type TimeProvider interface {
	Now() string
}

type UTCTimeProvider struct{}

func (t *UTCTimeProvider) Now() string {
	return time.Now().UTC().Format(time.RFC3339)
}
