package web

import (
	"net/http"

	"github.com/thijzert/chesseract/internal/notimplemented"
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

func (authResponseHandler) handleAuthResponse(p Provider, r AuthResponseRequest) (AuthResponseResponse, error) {
	return AuthResponseResponse{}, notimplemented.Error()
}

func (authResponseHandler) DecodeRequest(r *http.Request) (Request, error) {
	return AuthResponseRequest{}, nil
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
