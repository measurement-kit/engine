package internal

import (
	"errors"
	"time"
)

func MakeDuration(timeout int64) (time.Duration, error) {
	const maxTimeout = int64(120)
	if timeout < 0 || timeout > maxTimeout {
		return time.Duration(0), errors.New("timeout is negative or too large")
	}
	duration := time.Duration(timeout) * time.Second
	return duration, nil
}
