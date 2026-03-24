package chartmogul

const (
	accountEndpoint = "account"
)

// Account details in ChartMogul
type Account struct {
	ID                                string                 `json:"id,omitempty"`
	Name                              string                 `json:"name"`
	Currency                          string                 `json:"currency"`
	TimeZone                          string                 `json:"time_zone"`
	WeekStartOn                       string                 `json:"week_start_on"`
	ChurnRecognition                  map[string]interface{} `json:"churn_recognition,omitempty"`
	ChurnWhenZeroMRR                  map[string]interface{} `json:"churn_when_zero_mrr,omitempty"`
	AutoChurnSubscription             map[string]interface{} `json:"auto_churn_subscription,omitempty"`
	RefundHandling                    map[string]interface{} `json:"refund_handling,omitempty"`
	ProximateMovementReclassification map[string]interface{} `json:"proximate_movement_reclassification,omitempty"`
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
		return result, api.retrieveWithParams(accountEndpoint, "", result, *params[0])
	}
	return result, api.retrieve(accountEndpoint, "", result)
}
