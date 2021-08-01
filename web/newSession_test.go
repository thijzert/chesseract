package web

import (
	"net/http"
	"testing"
)

func TestDecodeNewSessionRequest(t *testing.T) {
	r, err := http.NewRequest("IMPLEMENT", "https://example.org/unittest/for/newSession", nil)
	if err != nil {
		t.Errorf("error creating dummy request: %s", err)
	}

	req, err := NewSessionHandler.DecodeRequest(r)
	if err != nil {
		t.Errorf("req: %+v; error: %s", req, err)
	}
}

func TestHandleNewSession(t *testing.T) {
	var p Provider

	req := newSessionRequest{}

	resp, err := NewSessionHandler.handleNewSession(p, req)

	t.Logf("response: %+v; error: %s", resp, err)
	t.Logf("TODO: implement unit test for handling NewSession")
}
