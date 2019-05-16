package chartmogul

import (
	"net"
	"net/http"
	"strconv"
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
	return strings.Join([]string{strconv.Itoa(e.statusCode), e.status, e.response}, ": ")
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

// List of HTTP status codes which are retryable.
var retryableHTTPStatusCodes = map[int]struct{}{
	http.StatusTooManyRequests:     {},
	http.StatusInternalServerError: {},
	http.StatusBadGateway:          {},
	http.StatusServiceUnavailable:  {},
	// Add more HTTP status codes here.
}

//isHTTPStatusRetryable return true if error message contains the HTTP statuses that needs to be retried
func isHTTPStatusRetryable(res gorequest.Response) (ok bool) {
	if res == nil {
		return false
	}
	_, ok = retryableHTTPStatusCodes[res.StatusCode]
	return ok
}

// timeoutError - when reading response body takes too long
// It's inside an internal package, so only way to assert this error is to match with following error string
// https://github.com/golang/go/blob/master/src/internal/poll/fd.go#L40-L45
const timeoutError = "i/o timeout"

// networkError returns true if the underlying error is caused by net.OpError
// or if it's an i/o timeout error
func networkError(err error) bool {
	if strings.Contains(err.Error(), timeoutError) {
		return true
	}

	if _, ok := (err).(*net.OpError); ok {
		return true
	}
	return false
}

// networkErrors checks if any of the errors is a Network related error
func networkErrors(errs []error) bool {
	for _, err := range errs {
		if networkError(err) {
			return true
		}
	}
	return false
}
