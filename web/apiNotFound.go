package web

import (
	"net/http"
)

var ApiNotFoundHandler apiNotFoundHandler

type apiNotFoundHandler struct{}

type apiNotFoundRequest struct {
}

// The ApiNotFoundResponse wraps a ApiNotFoundHandler API response
type ApiNotFoundResponse struct {
}

func (apiNotFoundHandler) handleApiNotFound(p Provider, r apiNotFoundRequest) (ApiNotFoundResponse, error) {
	return ApiNotFoundResponse{}, errNotFound("Not found", "This API target could not be found")
}

func (apiNotFoundHandler) DecodeRequest(r *http.Request) (Request, error) {
	return apiNotFoundRequest{}, nil
}

// Below: boilerplate code

func (h apiNotFoundHandler) HandleRequest(p Provider, r Request) (Response, error) {
	req, ok := r.(apiNotFoundRequest)
	if !ok {
		return withError(errWrongRequestType{})
	}

	return h.handleApiNotFound(p, req)
}

func (apiNotFoundRequest) FlaggedAsRequest() {}

func (ApiNotFoundResponse) FlaggedAsResponse() {}
