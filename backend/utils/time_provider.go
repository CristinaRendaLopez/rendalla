package utils

import (
	"time"
)

// TimeProvider abstracts access to the current time.
// It allows services to inject a testable clock and supports both
// human-readable and Unix timestamp formats.
type TimeProvider interface {
	// Now returns the current time as an RFC3339-formatted string.
	Now() string

	// NowUnix returns the current time as a Unix timestamp (seconds since epoch).
	NowUnix() int64
}

// UTCTimeProvider implements TimeProvider using the system clock in UTC.
type UTCTimeProvider struct{}

// Now returns the current UTC time as an RFC3339 string.
func (t *UTCTimeProvider) Now() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// NowUnix returns the current UTC time as a Unix timestamp.
func (t *UTCTimeProvider) NowUnix() int64 {
	return time.Now().UTC().Unix()
}
