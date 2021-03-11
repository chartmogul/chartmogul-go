package chartmogul

import (
	"errors"
	"strings"

	backoff "github.com/cenkalti/backoff/v3"
	"github.com/parnurzeal/gorequest"
)

// The methods here encompass common boilerplate for CRUD REST operations.
// They'd be generics if Go had generics...
// so these methods do *not* properly check types.
// Eg. nils cannot be easily checked (without reflection).

// static internal error helper
var errRetry = errors.New("Retrying")

// CREATE
func (api API) create(path string, input interface{}, output interface{}) error {
	var res gorequest.Response
	var body []byte
	var errs []error

	// Retry on HTTP 429 rate limit, or network error, see:
	// https://dev.chartmogul.com/docs/rate-limits
	// https://godoc.org/github.com/cenkalti/backoff#pkg-constants
	// nolint:errcheck
	backoff.Retry(func() error {
		res, body, errs = api.req(gorequest.New().
			Post(prepareURL(path))).
			SendStruct(input).
			EndStruct(output)

		if networkErrors(errs) || isHTTPStatusRetryable(res) {
			return errRetry
		}
		return nil
	}, backoff.NewExponentialBackOff())

	// wrapping []errors into compatible error & making HTTPError
	return wrapErrors(res, body, errs)
}

// READ
func (api API) list(path string, output interface{}, query ...interface{}) error {
	var res gorequest.Response
	var body []byte
	var errs []error

	// nolint:errcheck
	backoff.Retry(func() error {
		req := api.req(gorequest.New().Get(prepareURL(path)))
		for _, q := range query {
			req.Query(q)
		}

		res, body, errs = req.EndStruct(output)
		if networkErrors(errs) || isHTTPStatusRetryable(res) {
			return errRetry
		}
		return nil
	}, backoff.NewExponentialBackOff())

	return wrapErrors(res, body, errs)
}

// RETRIEVE
func (api API) retrieve(path string, uuid string, output interface{}) error {
	var res gorequest.Response
	var body []byte
	var errs []error
	if uuid != "" {
		path = strings.Replace(path, ":uuid", uuid, 1)
	}

	// nolint:errcheck
	backoff.Retry(func() error {
		res, body, errs = api.req(gorequest.New().Get(prepareURL(path))).
			EndStruct(output)

		if networkErrors(errs) || isHTTPStatusRetryable(res) {
			return errRetry
		}
		return nil
	}, backoff.NewExponentialBackOff())

	return wrapErrors(res, body, errs)
}

// UPDATE
func (api API) merge(path string, input interface{}) error {
	var res gorequest.Response
	var body string
	var errs []error

	// nolint:errcheck
	backoff.Retry(func() error {
		res, body, errs = api.req(gorequest.New().
			Post(prepareURL(path))).
			SendStruct(input).
			End()

		if networkErrors(errs) || isHTTPStatusRetryable(res) {
			return errRetry
		}
		return nil
	}, backoff.NewExponentialBackOff())

	return wrapErrors(res, []byte(body), errs)
}

// updateImpl adds another meta level, because this same pattern
// uses multiple HTTP methods in  API.
func (api API) updateImpl(path string, uuid string, input interface{}, output interface{}, method string) error {
	var res gorequest.Response
	var body []byte
	var errs []error

	path = strings.Replace(path, ":uuid", uuid, 1)
	path = prepareURL(path)

	// nolint:errcheck
	backoff.Retry(func() error {
		req := gorequest.New()
		switch method {
		case "update":
			req = req.Patch(path)
		case "add":
			req = req.Post(path)
		case "putTo":
			req = req.Put(path)
		}
		req = api.req(req).SendStruct(input)
		res, body, errs = req.EndStruct(output)

		if networkErrors(errs) || isHTTPStatusRetryable(res) {
			return errRetry
		}
		return nil
	}, backoff.NewExponentialBackOff())

	return wrapErrors(res, body, errs)
}

func (api API) update(path string, uuid string, input interface{}, output interface{}) error {
	return api.updateImpl(path, uuid, input, output, "update")
}

// add is like update, but POST
func (api API) add(path string, uuid string, input interface{}, output interface{}) error {
	return api.updateImpl(path, uuid, input, output, "add")
}

// putTo is like update, but PUT
func (api API) putTo(path string, uuid string, input interface{}, output interface{}) error {
	return api.updateImpl(path, uuid, input, output, "putTo")
}

// DELETE
func (api API) delete(path string, uuid string) error {
	var res gorequest.Response
	var body string
	var errs []error
	path = strings.Replace(path, ":uuid", uuid, 1)

	// nolint:errcheck
	backoff.Retry(func() error {
		res, body, errs = api.req(gorequest.New().Delete(prepareURL(path))).
			End()

		if networkErrors(errs) || isHTTPStatusRetryable(res) {
			return errRetry
		}
		return nil
	}, backoff.NewExponentialBackOff())

	return wrapErrors(res, []byte(body), errs)
}

func (api API) deleteWhat(path string, uuid string, input interface{}, output interface{}) error {
	var res gorequest.Response
	var body []byte
	var errs []error
	path = strings.Replace(path, ":uuid", uuid, 1)

	// nolint:errcheck
	backoff.Retry(func() error {
		res, body, errs = api.req(gorequest.New().Delete(prepareURL(path))).
			SendStruct(input).
			EndStruct(output)

		if networkErrors(errs) || isHTTPStatusRetryable(res) {
			return errRetry
		}
		return nil
	}, backoff.NewExponentialBackOff())

	return wrapErrors(res, body, errs)
}
