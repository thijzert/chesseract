package web

import (
	"net/http"
	"testing"
)

func TestDecodeAuthChallengeRequest(t *testing.T) {
	r, err := http.NewRequest("IMPLEMENT", "https://example.org/unittest/for/authChallenge", nil)
	if err != nil {
		t.Errorf("error creating dummy request: %s", err)
	}
	req, err := AuthChallengeHandler.DecodeRequest(r)

	t.Logf("req: %+v; error: %s", req, err)
	t.Logf("TODO: implement unit test for decoding AuthChallengeRequests")
}

func TestHandleAuthChallenge(t *testing.T) {
	var p Provider

	req := AuthChallengeRequest{}

	resp, err := AuthChallengeHandler.handleAuthChallenge(p, req)

	t.Logf("response: %+v; error: %s", resp, err)
	t.Logf("TODO: implement unit test for handling AuthChallenge")
}
