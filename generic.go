package chartmogul

import (
	"strings"

	"github.com/parnurzeal/gorequest"
)

// The methods here encompass common boilerplate for CRUD REST operations.
// They'd be generics if Go had generics...
// so these methods do *not* properly check types.
// Eg. nils cannot be easily checked (without reflection).

// CREATE
func (api API) create(path string, input interface{}, output interface{}) error {
	res, body, errs := api.req(gorequest.New().Post(prepareURL(path))).
		SendStruct(input).
		EndStruct(output)

	return wrapErrors(res, body, errs)
}

// READ
func (api API) list(path string, output interface{}, query ...interface{}) error {
	req := api.req(gorequest.New().Get(prepareURL(path)))
	for _, q := range query {
		req.Query(q)
	}
	res, body, errs := req.EndStruct(output)

	return wrapErrors(res, body, errs)
}

func (api API) retrieve(path string, uuid string, output interface{}) error {
	path = strings.Replace(path, ":uuid", uuid, 1)
	res, body, errs := api.req(gorequest.New().Get(prepareURL(path))).
		EndStruct(output)
	return wrapErrors(res, body, errs)
}

// UPDATE
func (api API) merge(path string, input interface{}) error {
	res, body, errs := api.req(gorequest.New().Post(prepareURL(path))).
		SendStruct(input).
		End()
	return wrapErrors(res, []byte(body), errs)
}

// updateImpl adds another meta level, because this same pattern
// uses multiple HTTP methods in  API.
func (api API) updateImpl(path string, uuid string, input interface{}, output interface{}, method string) error {
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

	res, body, errs := api.req(req).
		SendStruct(input).
		EndStruct(output)

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
	path = strings.Replace(path, ":uuid", uuid, 1)
	res, body, errs := api.req(gorequest.New().Delete(prepareURL(path))).
		End()
	return wrapErrors(res, []byte(body), errs)
}

func (api API) deleteWhat(path string, uuid string, input interface{}, output interface{}) error {
	path = strings.Replace(path, ":uuid", uuid, 1)
	res, body, errs := api.req(gorequest.New().Delete(prepareURL(path))).
		SendStruct(input).
		EndStruct(output)
	return wrapErrors(res, body, errs)
}
