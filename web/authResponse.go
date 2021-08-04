package web

import (
	"encoding/json"
	"errors"
	"net/http"

	weberrors "github.com/thijzert/chesseract/internal/web-plumbing/errors"
)

var AuthResponseHandler authResponseHandler

type authResponseHandler struct{}

type AuthResponseRequest struct {
	Username string `json:"username"`
	Nonce    string `json:"nonce"`
	Response string `json:"response"`
}

// The AuthResponseResponse wraps a AuthResponseHandler API response
type AuthResponseResponse struct {
}

func (authResponseHandler) err401() error {
	err := errors.New("authorisation required")
	err = weberrors.WithStatus(err, 401)
	return weberrors.WithMessage(err, "Authorisation required", "Your authorisation request failed. By all means, keep trying.")
}

func (h authResponseHandler) handleAuthResponse(p Provider, r AuthResponseRequest) (AuthResponseResponse, error) {
	var rv AuthResponseResponse

	player, ok, err := p.LookupPlayer(r.Username)
	if err != nil {
		return rv, err
	}
	if !ok {
		return rv, h.err401()
	}

	ok, err = p.ValidateNonce(r.Username, r.Nonce)
	if err != nil {
		return rv, err
	}
	if !ok {
		return rv, h.err401()
	}

	err = p.SetPlayer(player)
	if err != nil {
		return rv, err
	}

	// TODO: some form of authentication.

	// TODO: this should invalidate all other sessions for this player

	return rv, nil
}

func (authResponseHandler) DecodeRequest(r *http.Request) (Request, error) {
	var rv AuthResponseRequest

	if r.Body == nil {
		return rv, errMethod("Method not allowed", "This is a POST resource")
	}
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&rv)

	return rv, err
}

// Below: boilerplate code

func (h authResponseHandler) HandleRequest(p Provider, r Request) (Response, error) {
	req, ok := r.(AuthResponseRequest)
	if !ok {
		return withError(errWrongRequestType{})
	}

	return h.handleAuthResponse(p, req)
}

func (AuthResponseRequest) FlaggedAsRequest() {}

func (AuthResponseResponse) FlaggedAsResponse() {}
