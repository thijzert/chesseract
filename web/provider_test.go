package web

import (
	"github.com/thijzert/chesseract/chesseract/game"
	"github.com/thijzert/chesseract/internal/notimplemented"
)

type testProvider struct{}

// NewSession generates a new empty session, and returns a string
// representation of its ID, to be communicated to the client.
func (t testProvider) NewSession() (string, error) {
	return "", notimplemented.Error()
}

// Player returns the player associated with this session
func (t testProvider) Player() (game.Player, error) {
	return game.Player{}, notimplemented.Error()
}

// SetPlayer assigns this player to this session
func (t testProvider) SetPlayer(game.Player) error {
	return notimplemented.Error()
}

// LookupPlayer finds the profile in the database, if it exists
func (t testProvider) LookupPlayer(string) (game.Player, bool, error) {
	return game.Player{}, false, notimplemented.Error()
}

// NewNonce generates a new auth challenge for this player
func (t testProvider) NewNonce(string) (string, error) {
	return "", notimplemented.Error()
}

// ValidateNonce checks if a nonce is valid for this player
func (t testProvider) ValidateNonce(playerName string, nonce string) (bool, error) {
	return false, notimplemented.Error()
}

// NewGame creates a new game with the specified players, and returns its game ID
func (t testProvider) NewGame(ruleset string, playerNames []string) (string, error) {
	return "", notimplemented.Error()
}
