package web

import (
	"net/http"

	"github.com/thijzert/chesseract/chesseract/game"
)

var GetGameHandler getGameHandler

type getGameHandler struct{}

type getGameRequest struct {
}

// The GetGameResponse wraps a GetGameHandler API response
type GetGameResponse struct {
	Game *game.Game
}

func (getGameHandler) handleGetGame(p Provider, r getGameRequest) (GetGameResponse, error) {
	var rv GetGameResponse

	g, err := p.Game()
	if err != nil {
		return rv, err
	}
	rv.Game = g

	return rv, nil
}

func (getGameHandler) DecodeRequest(r *http.Request) (Request, error) {
	var rv getGameRequest
	var err error

	// if r.Body == nil {
	// 	return rv, errMethod("Method not allowed", "This is a POST resource")
	// }
	// dec := json.NewDecoder(r.Body)
	// err = dec.Decode(&rv)

	return rv, err
}

// Below: boilerplate code

func (h getGameHandler) HandleRequest(p Provider, r Request) (Response, error) {
	req, ok := r.(getGameRequest)
	if !ok {
		return withError(errWrongRequestType{})
	}

	return h.handleGetGame(p, req)
}

func (getGameRequest) FlaggedAsRequest() {}

func (GetGameResponse) FlaggedAsResponse() {}
