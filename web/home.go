package web

import "net/http"

var HomeHandler homeHandler

type homeHandler struct{}

func (homeHandler) handleHome(p Provider, r homeRequest) (homeResponse, error) {
	return homeResponse{}, nil
}

func (homeHandler) DecodeRequest(r *http.Request) (Request, error) {
	return homeRequest{}, nil
}

func (h homeHandler) HandleRequest(p Provider, r Request) (Response, error) {
	req, ok := r.(homeRequest)
	if !ok {
		return withError(errWrongRequestType{})
	}

	return h.handleHome(p, req)
}

type homeRequest struct {
}

func (homeRequest) FlaggedAsRequest() {}

type homeResponse struct{}

func (homeResponse) FlaggedAsResponse() {}
