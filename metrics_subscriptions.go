package chartmogul

import "strings"

// MetricsSubscription represents Metrics API subscription in ChartMogul.
type MetricsSubscription struct {
	ID                uint64  `json:"id"`
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

// MetricsSubscriptions is the result of listing subscriptions in Metrics API.
type MetricsSubscriptions struct {
	Entries []*MetricsSubscription `json:"entries"`
	HasMore bool                   `json:"has_more"`
	PerPage uint32                 `json:"per_page"`
	Page    uint32                 `json:"page"`
}

const metricsSubscriptionsEndpoint = "customers/:uuid/subscriptions"

// MetricsListSubscriptions lists all subscriptions for cutomer of a given UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#list-customer-subscriptions
func (api API) MetricsListSubscriptions(cursor *Cursor, customerUUID string) (*MetricsSubscriptions, error) {
	result := &MetricsSubscriptions{}
	path := strings.Replace(metricsSubscriptionsEndpoint, ":uuid", customerUUID, 1)
	query := make([]interface{}, 0, 1)
	if cursor != nil {
		query = append(query, *cursor)
	}
	return result, api.list(path, result, query...)
}
