package chartmogul

import (
	"strings"

	"github.com/parnurzeal/gorequest"
)

// HTTPError is wrapper to easily handle HTTP states.
type HTTPError struct {
	StatusCode int
	Status     string
	Response   string
}

func (e HTTPError) Error() string {
	return strings.Join([]string{string(e.StatusCode), e.Status, e.Response}, ": ")
}

// RequestErrors wraps multiple request errors to normal 1 error struct.
type RequestErrors struct {
	Errors []error
}

func (e RequestErrors) Error() string {
	errs := make([]string, len(e.Errors))
	for i := range errs {
		errs[i] = e.Errors[i].Error()
	}
	return strings.Join(errs, "; ")
}

// LogErrors converts any unexpected HTTP statuses on API to errors wrapped in a handy struct.
// If there are multiple errors with request, it returns them as a wrapper struct RequestErrors.
//
// In case of no errors returns nil.
func wrapErrors(response gorequest.Response, body []byte, errs []error) error {
	if response != nil && (response.StatusCode < 200 || response.StatusCode >= 300) {
		return &HTTPError{
			StatusCode: response.StatusCode,
			Status:     response.Status,
			Response:   string(body)}
	}
	if len(errs) != 0 {
		return &RequestErrors{errs}
	}
	return nil
}
