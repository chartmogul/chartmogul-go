package chartmogul

import "strings"

// MetricsCustomerSubscription represents Metrics API subscription in ChartMogul.
type MetricsCustomerSubscription struct {
	ID                uint64  `json:"id"`
	ExternalID        string  `json:"external_id"`
	Plan              string  `json:"plan"`
	Quantity          uint32  `json:"quantity"`
	MRR               float64 `json:"mrr"`
	ARR               float64 `json:"arr"`
	Status            string  `json:"status"`
	BillingCycle      string  `json:"billing-cycle"`
	BillingCycleCount uint32  `json:"billing-cycle-count"`
	StartDate         string  `json:"start-date"`
	EndDate           string  `json:"end-date"`
	Currency          string  `json:"currency"`
	CurrencySign      string  `json:"currency-sign"`
}

// MetricsCustomerSubscriptions is the result of listing subscriptions in Metrics API.
type MetricsCustomerSubscriptions struct {
	Entries []*MetricsCustomerSubscription `json:"entries"`
	HasMore bool                   `json:"has_more"`
	PerPage uint32                 `json:"per_page"`
	Page    uint32                 `json:"page"`
}

const metricsCustomerSubscriptionsEndpoint = "customers/:uuid/subscriptions"

// MetricsListSubscriptions lists all subscriptions for cutomer of a given UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#list-customer-subscriptions
func (api API) MetricsListCustomerSubscriptions(cursor *Cursor, customerUUID string) (*MetricsCustomerSubscriptions, error) {
	result := &MetricsCustomerSubscriptions{}
	path := strings.Replace(metricsCustomerSubscriptionsEndpoint, ":uuid", customerUUID, 1)
	query := make([]interface{}, 0, 1)
	if cursor != nil {
		query = append(query, *cursor)
	}
	return result, api.list(path, result, query...)
}
