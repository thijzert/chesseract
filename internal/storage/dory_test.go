package storage

import (
	"context"
	"testing"
)

func TestInitialise(t *testing.T) {
	var b Backend = &Dory{}

	err := b.Initialise(context.Background())

	if err != nil {
		t.Logf("Error initialising: %v", err)
		t.Fail()
	}
}
