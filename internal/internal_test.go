package internal

import (
	"testing"
	"time"
)

// TestMakeTimeoutOK checks whether MakeTimeout works in the common case.
func TestMakeTimeoutOK(t *testing.T) {
	v, err := MakeTimeout(10)
	if err != nil {
		t.Fatal(err)
	}
	if v != 10*time.Second {
		t.Fatal("The timeout was incorrectly set")
	}
}

// TestMakeTimeoutNegative checks whether MakeTimeout correctly fails
// when the timeout is negative.
func TestMakeTimeoutNegative(t *testing.T) {
	_, err := MakeTimeout(-1)
	if err == nil {
		t.Fatal("We were expecting an error here")
	}
}

// TestMakeTimeoutTooLarge checks whether MakeTimeout correctly fails
// when the timeout is too large.
func TestMakeTimeoutTooLarge(t *testing.T) {
	_, err := MakeTimeout(maxTimeout + 1)
	if err == nil {
		t.Fatal("We were expecting an error here")
	}
}
