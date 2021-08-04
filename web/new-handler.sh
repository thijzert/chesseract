#!/bin/bash

cd "$(dirname "$0")"

TYPE="$1"
TYPE="${TYPE%Handler}"
TYPE="${TYPE,}"

if [ -z "$TYPE" ]
then
	echo "Usage: $0  handlerTypeName" >&2
	exit 1
fi

if [[ "$TYPE" =~ ( |\.|/|\\|-) ]]
then
	echo "Illegal handler name." >&2
	echo "Usage: $0  handlerTypeName" >&2
	exit 1
fi

OUT_FILE="${TYPE}.go"
OUT_TEST="${TYPE}_test.go"

if [ -f "$OUT_FILE" -o -f "$OUT_TEST" ]
then
	echo "Warning: this script will overwrite $OUT_FILE and $OUT_TEST." >&2
	echo "Please make sure you don't lose what's currently in there." >&2
	exit 1
fi

cat > "$OUT_FILE" <<EOF
package web

import (
	"net/http"

	"github.com/thijzert/chesseract/internal/notimplemented"
)

var ${TYPE^}Handler ${TYPE}Handler

type ${TYPE}Handler struct{}

type ${TYPE}Request struct {
}

// The ${TYPE^}Response wraps a ${TYPE^}Handler API response
type ${TYPE^}Response struct{
}

func (${TYPE}Handler) handle${TYPE^}(p Provider, r ${TYPE}Request) (${TYPE^}Response, error) {
	return ${TYPE^}Response{}, notimplemented.Error()
}

func (${TYPE}Handler) DecodeRequest(r *http.Request) (Request, error) {
	var rv ${TYPE}Request

	// if r.Body == nil {
	// 	return rv, errMethod("Method not allowed", "This is a POST resource")
	// }
	// dec := json.NewDecoder(r.Body)
	// err := dec.Decode(&rv)

	return rv, nil
}

// Below: boilerplate code

func (h ${TYPE}Handler) HandleRequest(p Provider, r Request) (Response, error) {
	req, ok := r.(${TYPE}Request)
	if !ok {
		return withError(errWrongRequestType{})
	}

	return h.handle${TYPE^}(p, req)
}

func (${TYPE}Request) FlaggedAsRequest() {}

func (${TYPE^}Response) FlaggedAsResponse() {}

EOF

cat > "$OUT_TEST" <<EOF
package web

import (
	"net/http"
	"testing"
)

func TestDecode${TYPE^}Request(t *testing.T) {
	r, err := http.NewRequest("IMPLEMENT", "https://example.org/unittest/for/${TYPE}", nil)
	if err != nil {
		t.Errorf("error creating dummy request: %s", err)
	}
	req, err := ${TYPE^}Handler.DecodeRequest(r)

	t.Logf("req: %+v; error: %s", req, err)
	t.Logf("TODO: implement unit test for decoding ${TYPE^}Requests")
}

func TestHandle${TYPE^}(t *testing.T) {
	var p Provider = testProvider{}

	req := ${TYPE}Request{
	}

	resp, err := ${TYPE^}Handler.handle${TYPE^}(p, req)

	t.Logf("response: %+v; error: %s", resp, err)
	t.Logf("TODO: implement unit test for handling ${TYPE^}")
}

EOF
