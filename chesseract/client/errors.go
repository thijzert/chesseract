package client

import (
	"encoding/json"
	"fmt"
)

var (
	ErrShenanigans     error = clientError(2)
	ErrUnknownPlayer   error = clientError(4)
	ErrUnknownGame     error = clientError(5)
	ErrIllegalMove     error = clientError(6)
	ErrNotYourTurn     error = clientError(7)
	ErrGameHasFinished error = clientError(8)
	ErrInvalidResult   error = clientError(9)
)

type clientError int

func (c clientError) Error() string {
	if c == 2 {
		return "not all parties are in agreement about the rules of this game"
	} else if c == 4 {
		return "unknown player"
	} else if c == 5 {
		return "unknown game"
	} else if c == 6 {
		return "illegal move"
	} else if c == 7 {
		return "not your turn"
	} else if c == 8 {
		return "game has finished"
	} else if c == 9 {
		return "invalid result value"
	}

	return fmt.Sprintf("unknown error %x", int(c))
}

type jsonClientError struct {
	ErrorCode    int
	ErrorMessage string
}

func (c clientError) MarshalJSON() ([]byte, error) {
	rv := jsonClientError{
		ErrorCode:    int(c),
		ErrorMessage: c.Error(),
	}
	return json.Marshal(rv)
}

func (c *clientError) UnmarshalJSON(data []byte) error {
	var rv jsonClientError
	err := json.Unmarshal(data, &rv)
	if err != nil {
		return err
	}
	*c = clientError(rv.ErrorCode)
	return nil
}
