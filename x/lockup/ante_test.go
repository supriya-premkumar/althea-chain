package lockup

import (
	"testing"
)

func TestAnteHandler(t *testing.T) {
	// TODO: Create a new ante decorator with default state and test that it
	// blocks a transaction coming from 0x1 but not one coming from 0x0.

	// TODO: Test that messages not of the right type aren't blocked

}

// TODO: Test that empty lock exempt fails validation

// TODO: Test that empty locked message types fails validation

// TODO: Test that nothing is blocked when locked is false
