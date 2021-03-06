// Package internal contains implementation details.
package internal

import (
	"errors"
	"time"
)

// maxTimeout is the maximum assignable timeout
const maxTimeout = int64(120)

// MakeTimeout converts a timeout to time.Duration. This function will
// fail if the timeout is negative or too big.
func MakeTimeout(timeout int64) (time.Duration, error) {
	if timeout < 0 || timeout > maxTimeout {
		return time.Duration(0), errors.New("timeout is negative or too large")
	}
	duration := time.Duration(timeout) * time.Second
	return duration, nil
}
