package chartmogul

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cenkalti/backoff"
	"github.com/parnurzeal/gorequest"
)

// The methods here encompass common boilerplate for CRUD REST operations.
// They'd be generics if Go had generics...
// so these methods do *not* properly check types.
// Eg. nils cannot be easily checked (without reflection).

// static internal error helper
var errRateLimit = errors.New("Rate limit reached")

// CREATE
func (api API) create(path string, input interface{}, output interface{}) error {
	var res gorequest.Response
	var body []byte
	var errs []error
	// This request object will be reused for consequent retries
	req := api.req(gorequest.New().
		Post(prepareURL(path))).
		SendStruct(input)

	// Retry on HTTP 429 rate limit, see:
	// https://dev.chartmogul.com/docs/rate-limits
	// https://godoc.org/github.com/cenkalti/backoff#pkg-constants
	backoff.Retry(func() error {
		res, body, errs = req.EndStruct(output)
		if res.StatusCode == http.StatusTooManyRequests {
			return errRateLimit
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
	req := api.req(gorequest.New().Get(prepareURL(path)))
	for _, q := range query {
		req.Query(q)
	}

	backoff.Retry(func() error {
		res, body, errs = req.EndStruct(output)
		if res.StatusCode == http.StatusTooManyRequests {
			return errRateLimit
		}
		return nil
	}, backoff.NewExponentialBackOff())

	return wrapErrors(res, body, errs)
}

func (api API) retrieve(path string, uuid string, output interface{}) error {
	var res gorequest.Response
	var body []byte
	var errs []error
	path = strings.Replace(path, ":uuid", uuid, 1)
	req := api.req(gorequest.New().Get(prepareURL(path)))

	backoff.Retry(func() error {
		res, body, errs = req.EndStruct(output)
		if res.StatusCode == http.StatusTooManyRequests {
			return errRateLimit
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
	req := api.req(gorequest.New().
		Post(prepareURL(path))).
		SendStruct(input)

	backoff.Retry(func() error {
		res, body, errs = req.End()
		if res.StatusCode == http.StatusTooManyRequests {
			return errRateLimit
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

	backoff.Retry(func() error {
		res, body, errs = req.EndStruct(output)
		if res.StatusCode == http.StatusTooManyRequests {
			return errRateLimit
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
	req := api.req(gorequest.New().Delete(prepareURL(path)))

	backoff.Retry(func() error {
		res, body, errs = req.End()
		if res.StatusCode == http.StatusTooManyRequests {
			return errRateLimit
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
	req := api.req(gorequest.New().Delete(prepareURL(path))).
		SendStruct(input)

	backoff.Retry(func() error {
		res, body, errs = req.EndStruct(output)
		if res.StatusCode == http.StatusTooManyRequests {
			return errRateLimit
		}
		return nil
	}, backoff.NewExponentialBackOff())

	return wrapErrors(res, body, errs)
}
