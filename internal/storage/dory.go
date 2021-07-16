package storage

import (
	"sync"
)

// The Dory storage backend implements a storage backend that forgets everything as soon as the program stops.
type Dory struct {
	mu sync.RWMutex
}

func (d *Dory) Initialise() error {
	d.mu.Lock()
	// Just keep swimming
	defer d.mu.Unlock()

	return nil
}
