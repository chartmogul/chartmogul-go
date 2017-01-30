package chartmogul

import (
	"github.com/chartmogul/gocardless-integration/util"
	"github.com/parnurzeal/gorequest"
)

// Ping is simple struct for the authentication test endpoint.
type Ping struct {
	Data string
}

const pingEndpoint = "ping"

// Ping is the authentication test endpoint.
//
// See https://dev.chartmogul.com/v1.0/docs/authentication
func (api API) Ping() (bool, error) {
	ping := &Ping{}
	res, body, errs := api.req(gorequest.New().Get(prepareURL(pingEndpoint))).EndStruct(ping)
	return ping.Data == "pong!", util.LogErrors(res, body, errs)
}
