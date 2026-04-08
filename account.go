package chartmogul

import (
	neturl "net/url"
	"strings"
)

const (
	accountEndpoint = "account"
)

// Account details in ChartMogul
type Account struct {
	ID                                string      `json:"id,omitempty"`
	Name                              string      `json:"name"`
	Currency                          string      `json:"currency"`
	TimeZone                          string      `json:"time_zone"`
	WeekStartOn                       string      `json:"week_start_on"`
	ChurnRecognition                  interface{} `json:"churn_recognition,omitempty"`
	ChurnWhenZeroMRR                  interface{} `json:"churn_when_zero_mrr,omitempty"`
	AutoChurnSubscription             interface{} `json:"auto_churn_subscription,omitempty"`
	RefundHandling                    interface{} `json:"refund_handling,omitempty"`
	ProximateMovementReclassification interface{} `json:"proximate_movement_reclassification,omitempty"`
}

// RetrieveAccountParams optional parameters for RetrieveAccount.
type RetrieveAccountParams struct {
	Include string `json:"include,omitempty"`
}

// RetrieveAccount returns details of current account.
// Optionally accepts RetrieveAccountParams with an Include field
// (comma-separated list: churn_recognition, churn_when_zero_mrr, etc.)
func (api API) RetrieveAccount(params ...*RetrieveAccountParams) (*Account, error) {
	result := &Account{}
	if len(params) > 0 && params[0] != nil && params[0].Include != "" {
		fields := params[0].Include
		// Normalize array-style input
		fields = strings.ReplaceAll(fields, " ", "")
		path := accountEndpoint + "?include=" + neturl.QueryEscape(fields)
		return result, api.retrieve(path, "", result)
	}
	return result, api.retrieve(accountEndpoint, "", result)
}
