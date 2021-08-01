package web

import (
	"net/http"
	"testing"
)

func TestDecodeAuthResponseRequest(t *testing.T) {
	r, err := http.NewRequest("IMPLEMENT", "https://example.org/unittest/for/authResponse", nil)
	if err != nil {
		t.Errorf("error creating dummy request: %s", err)
	}
	req, err := AuthResponseHandler.DecodeRequest(r)

	t.Logf("req: %+v; error: %s", req, err)
	t.Logf("TODO: implement unit test for decoding AuthResponseRequests")
}

func TestHandleAuthResponse(t *testing.T) {
	var p Provider = testProvider{}

	req := AuthResponseRequest{}

	resp, err := AuthResponseHandler.handleAuthResponse(p, req)

	t.Logf("response: %+v; error: %s", resp, err)
	t.Logf("TODO: implement unit test for handling AuthResponse")
}
