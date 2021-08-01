package weberrors

import (
	"errors"
)

type HTTPError interface {
	error
	HTTPStatus() int
}

type httpError struct {
	StatusCode int
	Cause      error
}

func WithStatus(e error, c int) HTTPError {
	if e == nil {
		return nil
	}

	return httpError{
		StatusCode: c,
		Cause:      e,
	}
}

func (e httpError) Error() string {
	return e.Cause.Error()
}

func (e httpError) Unwrap() error {
	return e.Cause
}

func (e httpError) HTTPStatus() int {
	return e.StatusCode
}

func HTTPStatusCode(e error) (statusCode int, cause error) {
	if e == nil {
		return 200, nil
	}

	var httpcode httpError
	if errors.As(e, &httpcode) {
		return httpcode.StatusCode, httpcode.Cause
	}

	var herr HTTPError
	if errors.As(e, &herr) {
		return herr.HTTPStatus(), herr
	}

	return 0, e
}

type CodeError interface {
	error
	ErrorCode() int
}

type codeError struct {
	Code  int
	Cause error
}

func WithCode(e error, c int) CodeError {
	if e == nil {
		return nil
	}

	return codeError{
		Code:  c,
		Cause: e,
	}
}

func (e codeError) Error() string {
	return e.Cause.Error()
}

func (e codeError) Unwrap() error {
	return e.Cause
}

func (e codeError) ErrorCode() int {
	return e.Code
}

func ErrorCode(e error) (errorCode int, cause error) {
	if e == nil {
		return 0, nil
	}

	var errcode codeError
	if errors.As(e, &errcode) {
		return errcode.Code, errcode.Cause
	}

	var err CodeError
	if errors.As(e, &err) {
		return err.ErrorCode(), err
	}

	lastResort, cause := HTTPStatusCode(e)
	if lastResort != 0 && lastResort != 200 {
		return lastResort, cause
	}

	return 1, e
}
