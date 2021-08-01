package plumbing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	weberrors "github.com/thijzert/chesseract/internal/web-plumbing/errors"
	"github.com/thijzert/chesseract/web"
)

type jsonHandler struct {
	Server  *Server
	Handler web.Handler
}

// JSONFunc creates a HTTP handler that outputs JSON
func (s *Server) JSONFunc(handler web.Handler) http.Handler {
	return jsonHandler{
		Server:  s,
		Handler: handler,
	}
}

func (h jsonHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := h.Handler.DecodeRequest(r)
	if err != nil {
		h.Error(w, r, err)
		return
	}

	provider := h.Server.getProvider(r)
	resp, err := h.Handler.HandleRequest(provider, req)
	if err != nil {
		h.Error(w, r, err)
		return
	}

	// Alternative path: this response can write its own headers and response body
	if h, ok := resp.(http.Handler); ok {
		h.ServeHTTP(w, r)
		return
	}

	w.Header()["Content-Type"] = []string{"application/json"}
	w.Header()["X-Content-Type-Options"] = []string{"nosniff"}

	var b bytes.Buffer
	e := json.NewEncoder(&b)
	err = e.Encode(resp)
	if err != nil {
		h.Error(w, r, err)
		return
	}

	io.Copy(w, &b)
}

func (jsonHandler) Error(w http.ResponseWriter, r *http.Request, err error) {

	w.Header()["Content-Type"] = []string{"application/json"}
	w.Header()["X-Content-Type-Options"] = []string{"nosniff"}

	st, _ := weberrors.HTTPStatusCode(err)
	if st == 0 {
		st = 500
	}
	w.WriteHeader(st)

	errorResponse := struct {
		Code     int    `json:"error_code"`
		Headline string `json:"error"`
		Message  string `json:"message"`
	}{}

	errorResponse.Code, _ = weberrors.ErrorCode(err)
	errorResponse.Headline = weberrors.Headline(err)
	errorResponse.Message = weberrors.Message(err)

	var b bytes.Buffer
	e := json.NewEncoder(&b)
	err = e.Encode(errorResponse)
	if err != nil {
		fmt.Fprintf(w, "{errorCode: 500, errorMessage: \"I give up.\"}")
	} else {
		io.Copy(w, &b)
	}
}
