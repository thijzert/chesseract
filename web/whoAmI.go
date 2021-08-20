package web

import (
	"net/http"

	"github.com/thijzert/chesseract/chesseract/game"
)

var WhoAmIHandler whoAmIHandler

type whoAmIHandler struct{}

type whoAmIRequest struct {
}

// The WhoAmIResponse wraps a WhoAmIHandler API response
type WhoAmIResponse struct {
	Profile game.Player `json:"profile"`
}

func (whoAmIHandler) handleWhoAmI(p Provider, r whoAmIRequest) (WhoAmIResponse, error) {
	var rv WhoAmIResponse

	player, err := p.Player()
	if err != nil {
		return rv, err
	}

	rv.Profile = player
	return rv, nil
}

func (whoAmIHandler) DecodeRequest(r *http.Request) (Request, error) {
	return whoAmIRequest{}, nil
}

// Below: boilerplate code

func (h whoAmIHandler) HandleRequest(p Provider, r Request) (Response, error) {
	req, ok := r.(whoAmIRequest)
	if !ok {
		return withError(errWrongRequestType{})
	}

	return h.handleWhoAmI(p, req)
}

func (whoAmIRequest) FlaggedAsRequest() {}

func (WhoAmIResponse) FlaggedAsResponse() {}
