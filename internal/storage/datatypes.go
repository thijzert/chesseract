package storage

import (
	crand "crypto/rand"
	"fmt"
	mrand "math/rand"
)

func randomInt64() uint64 {
	// Start with a sort-of-random number
	var rv uint64 = mrand.Uint64()

	// Mix some of the good stuff in there
	buf := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	n, _ := crand.Read(buf)
	for i, c := range buf[:n] {
		rv ^= uint64(c) << (8 * i)
	}

	return rv
}

type SessionID [4]uint64

// NewSessionID generates a new SessionID. The probability of colliding with a
// previously generated SessionID should be around 2^-128.
func NewSessionID() SessionID {
	return SessionID{randomInt64(), randomInt64(), randomInt64(), randomInt64()}
}

func (s SessionID) IsEmpty() bool {
	return s[0] == 0 && s[1] == 0 && s[2] == 0 && s[3] == 0
}

func (s SessionID) String() string {
	return fmt.Sprintf("%016x-%016x-%016x-%016x", s[0], s[1], s[2], s[3])
}

func ParseSessionID(str string) (SessionID, error) {
	var s SessionID
	_, err := fmt.Sscanf(str, "%x-%x-%x-%x", &s[0], &s[1], &s[2], &s[3])
	return s, err
}

type Session struct {
	PlayerID PlayerID
}

// NewPlayerID generates a new PlayerID. The probability of colliding with a
// previously generated PlayerID should be around 2^-64.
type PlayerID [2]uint64

func NewPlayerID() PlayerID {
	return PlayerID{randomInt64(), randomInt64()}
}

// NewGameID generates a new GameID. The probability of colliding with a
// previously generated GameID should be around 2^-64.
type GameID [2]uint64

func NewGameID() GameID {
	return GameID{randomInt64(), randomInt64()}
}
