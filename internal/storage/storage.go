package storage

type Backend interface {
	// Initialises sets up the Backend for use
	Initialise() error
}
