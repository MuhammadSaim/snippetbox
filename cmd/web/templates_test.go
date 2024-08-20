package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// Initialize a new time object and pass it to the function
	tm := time.Date(2024, 04, 05, 15, 15, 15, 0, time.UTC)
	hd := humanDate(tm)

	// Check that the func output matches to the value
	if hd != "05 Apr 2024 at 15:15" {
		t.Errorf("got %q; want %q", hd, "05 Apr 2024 at 15:15")
	}
}
