package storage

import "testing"

func TestInitialise(t *testing.T) {
	var b Backend = &Dory{}

	err := b.Initialise()

	if err != nil {
		t.Logf("Error initialising: %v", err)
		t.Fail()
	}
}
