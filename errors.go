package chartmogul

import (
	"strings"

	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

// HTTPError is wrapper to easily handle HTTP states.
type HTTPError interface {
	StatusCode() int
	Status() string
	Response() string
}

// RequestErrors wraps multiple request errors to normal 1 error struct.
type RequestErrors interface {
	Errors() []error
}

type httpError struct {
	statusCode int
	status     string
	response   string
}

func (e httpError) StatusCode() int {
	return e.statusCode
}
func (e httpError) Status() string {
	return e.status
}
func (e httpError) Response() string {
	return e.response
}
func (e httpError) Error() string {
	return strings.Join([]string{string(e.statusCode), e.status, e.response}, ": ")
}

type requestErrors struct {
	errors []error
}

func (e requestErrors) Errors() []error {
	return e.errors
}

func (e requestErrors) Error() string {
	errs := make([]string, len(e.errors))
	for i := range errs {
		errs[i] = e.errors[i].Error()
	}
	return strings.Join(errs, "; ")
}

// LogErrors converts any unexpected HTTP statuses on API to errors wrapped in a handy struct.
// If there are multiple errors with request, it returns them as a wrapper struct RequestErrors.
//
// In case of no errors returns nil.
func wrapErrors(response gorequest.Response, body []byte, errs []error) error {
	if response != nil && (response.StatusCode < 200 || response.StatusCode >= 300) {
		return errors.Wrap(&httpError{
			statusCode: response.StatusCode,
			status:     response.Status,
			response:   string(body)},
			"API error")
	}
	if len(errs) != 0 {
		return errors.Wrap(&requestErrors{errs}, "Request error")
	}
	return nil
}
