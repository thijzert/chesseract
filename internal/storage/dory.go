package storage

import (
	"fmt"
	"sync"

	"github.com/thijzert/chesseract/chesseract/game"
)

func init() {
	RegisterBackend("dory", func(string) (Backend, error) {
		return &Dory{}, nil
	})
}

// The Dory storage backend implements a storage backend that forgets everything as soon as the program stops.
type Dory struct {
	mu sync.RWMutex

	// sessions stores all sessions
	sessions map[SessionID]Session

	// players stores all players
	players map[PlayerID]game.Player

	// games stores all past and active games
	games map[GameID]game.Game
}

func (d *Dory) String() string {
	return "dummy storage backend; no data is sav- ooooh, what's that?"
}

func (d *Dory) Initialise() error {
	d.mu.Lock()
	// Just keep swimming
	defer d.mu.Unlock()

	d.sessions = make(map[SessionID]Session)
	d.players = make(map[PlayerID]game.Player)
	d.games = make(map[GameID]game.Game)

	return nil
}

func (d *Dory) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.sessions = nil
	d.players = nil
	d.games = nil

	return nil
}

// Transaction starts a transaction before running the function.
// If it returns any error, the transaction is rolled back. If it
// returns nil, it is committed.
// Transaction only returns an error if committing or rolling back fails.
func (d *Dory) Transaction(f func() error) error {
	// Transactions are just not supported
	err := f()
	if err != nil {
		return fmt.Errorf("transactions are not supported; unable to roll back")
	}

	return nil
}

// NewSession creates a new session
func (d *Dory) NewSession() (SessionID, Session, error) {
	sessionID := NewSessionID()
	defaultSession := Session{}

	return sessionID, defaultSession, d.StoreSession(sessionID, defaultSession)
}

// GetSession retrieves a session from the store
func (d *Dory) GetSession(id SessionID) (Session, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rv, ok := d.sessions[id]
	if !ok {
		return rv, errNotPresent
	}

	return rv, nil
}

// StoreSession updates a modified Session in the datastore
func (d *Dory) StoreSession(id SessionID, sess Session) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.sessions[id] = sess

	return nil
}

// NewPlayer creates a new player
func (d *Dory) NewPlayer() (PlayerID, game.Player, error) {
	id := NewPlayerID()
	player := game.Player{}

	return id, player, d.StorePlayer(id, player)
}

// GetPlayer retrieves a player from the store
func (d *Dory) GetPlayer(id PlayerID) (game.Player, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rv, ok := d.players[id]
	if !ok {
		return rv, errNotPresent
	}

	return rv, nil
}

// StorePlayer updates a modified Player in the datastore
func (d *Dory) StorePlayer(id PlayerID, player game.Player) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.players[id] = player

	return nil
}

// NewGame creates a new game
func (d *Dory) NewGame() (GameID, game.Game, error) {
	id := NewGameID()
	match := game.Game{}

	return id, match, d.StoreGame(id, match)
}

// GetGame retrieves a game from the store
func (d *Dory) GetGame(id GameID) (game.Game, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	rv, ok := d.games[id]
	if !ok {
		return rv, errNotPresent
	}

	return rv, nil
}

// StoreGame updates a modified Game in the datastore
func (d *Dory) StoreGame(id GameID, match game.Game) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.games[id] = match

	return nil
}
