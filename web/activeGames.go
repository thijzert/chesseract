package web

import (
	"net/http"
)

var ActiveGamesHandler activeGamesHandler

type activeGamesHandler struct{}

type activeGamesRequest struct {
}

// The ActiveGamesResponse wraps a ActiveGamesHandler API response
type ActiveGamesResponse struct {
	GameIDs []string `json:"gameid"`
}

func (activeGamesHandler) handleActiveGames(p Provider, r activeGamesRequest) (ActiveGamesResponse, error) {
	var rv ActiveGamesResponse
	var err error

	rv.GameIDs, err = p.ActiveGames()

	return rv, err
}

func (activeGamesHandler) DecodeRequest(r *http.Request) (Request, error) {
	var rv activeGamesRequest
	return rv, nil
}

// Below: boilerplate code

func (h activeGamesHandler) HandleRequest(p Provider, r Request) (Response, error) {
	req, ok := r.(activeGamesRequest)
	if !ok {
		return withError(errWrongRequestType{})
	}

	return h.handleActiveGames(p, req)
}

func (activeGamesRequest) FlaggedAsRequest() {}

func (ActiveGamesResponse) FlaggedAsResponse() {}
