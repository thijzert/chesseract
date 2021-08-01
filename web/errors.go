package web

import (
	"errors"
	"fmt"

	weberrors "github.com/thijzert/chesseract/internal/web-plumbing/errors"
)

type errorResponse struct {
	Error error
}

func (errorResponse) FlaggedAsResponse() {}

func withError(e error) (Response, error) {
	return errorResponse{e}, e
}

type errWrongRequestType struct{}

func (errWrongRequestType) Error() string {
	return "wrong request type"
}

func (errWrongRequestType) HTTPCode() int {
	return 400
}

type errRedirect struct {
	URL string
}

func (errRedirect) Error() string {
	return "you are being redirected to another page"
}

func (errRedirect) Headline() string {
	return "Redirecting..."
}

func (e errRedirect) Message() string {
	return fmt.Sprintf("You are being redirected to the address '%s'", e.URL)
}

func (errRedirect) HTTPCode() int {
	return 302
}

func (e errRedirect) RedirectLocation() string {
	return e.URL
}

func errForbidden(headline, message string) error {
	rv := errors.New("access denied")
	rv = weberrors.WithStatus(rv, 403)

	if headline == "" {
		headline = "Access Denied"
	}
	if message == "" {
		message = "You don't have permission to access this resource"
	}

	rv = weberrors.WithMessage(rv, headline, message)
	return rv
}

func errNotFound(headline, message string) error {
	rv := errors.New("not found")
	rv = weberrors.WithStatus(rv, 404)

	if headline == "" {
		headline = "Not Found"
	}
	if message == "" {
		message = "The document or resource you requested could not be found"
	}

	rv = weberrors.WithMessage(rv, headline, message)
	return rv
}

func errBadRequest(headline, message string) error {
	rv := errors.New("bad request")
	rv = weberrors.WithStatus(rv, 400)

	if headline == "" {
		headline = "Bad request"
	}
	if message == "" {
		message = "The server cannot or will not process the request due to an apparent client error"
	}

	rv = weberrors.WithMessage(rv, headline, message)
	return rv
}

func errMethod(headline, message string) error {
	rv := errors.New("method not allowed")
	rv = weberrors.WithStatus(rv, 405)

	if headline == "" {
		headline = "Method not allowed"
	}
	if message == "" {
		message = "A request method is not supported for the requested resource"
	}

	rv = weberrors.WithMessage(rv, headline, message)
	return rv
}
