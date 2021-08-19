package web

import (
	"encoding/json"
	"net/http"

	"github.com/thijzert/chesseract/chesseract/game"
)

var NewGameHandler newGameHandler

type newGameHandler struct{}

type NewGameRequest struct {
	RuleSet     string   `json:"ruleset"`
	PlayerNames []string `json:"players"`
}

// The NewGameResponse wraps a NewGameHandler API response
type NewGameResponse struct {
	GameID string     `json:"gameid"`
	Game   *game.Game `json:"game,omitempty"`
}

func (newGameHandler) handleNewGame(p Provider, r NewGameRequest) (NewGameResponse, error) {
	var rv NewGameResponse

	id, g, err := p.NewGame(r.RuleSet, r.PlayerNames)
	if err != nil {
		return rv, err
	}
	rv.GameID = id
	rv.Game = g

	return rv, nil
}

func (newGameHandler) DecodeRequest(r *http.Request) (Request, error) {
	var rv NewGameRequest

	if r.Body == nil {
		return rv, errMethod("Method not allowed", "This is a POST resource")
	}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&rv)

	return rv, err
}

// Below: boilerplate code

func (h newGameHandler) HandleRequest(p Provider, r Request) (Response, error) {
	req, ok := r.(NewGameRequest)
	if !ok {
		return withError(errWrongRequestType{})
	}

	return h.handleNewGame(p, req)
}

func (NewGameRequest) FlaggedAsRequest() {}

func (NewGameResponse) FlaggedAsResponse() {}
