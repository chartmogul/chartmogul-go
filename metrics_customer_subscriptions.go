package chartmogul

import "strings"

// MetricsCustomerSubscription represents Metrics API subscription in ChartMogul.
type MetricsCustomerSubscription struct {
	ID                uint64  `json:"id"`
	ExternalID        string  `json:"external_id"`
	Plan              string  `json:"plan"`
	Quantity          uint32  `json:"quantity"`
	BillingCycleCount uint32  `json:"billing-cycle-count"`
	MRR               float64 `json:"mrr"`
	ARR               float64 `json:"arr"`
	Status            string  `json:"status"`
	BillingCycle      string  `json:"billing-cycle"`
	StartDate         string  `json:"start-date"`
	EndDate           string  `json:"end-date"`
	Currency          string  `json:"currency"`
	CurrencySign      string  `json:"currency-sign"`
	UUID              string  `json:"uuid,omitempty"` // UUID field for connect/disconnect operations
}

// MetricsCustomerSubscriptions is the result of listing subscriptions in Metrics API.
type MetricsCustomerSubscriptions struct {
	Entries []*MetricsCustomerSubscription `json:"entries"`
	Pagination
}

// MetricsConnectSubscriptionsParams represents subscriptions data for connect/disconnect operations.
type MetricsConnectSubscriptionsParams struct {
	Subscriptions []MetricsSubscriptionReference `json:"subscriptions"`
}

// MetricsSubscriptionReference represents a minimal subscription reference with UUID and data source.
type MetricsSubscriptionReference struct {
	UUID           string `json:"uuid"`
	DataSourceUUID string `json:"data_source_uuid"`
}

const metricsCustomerSubscriptionsEndpoint = "customers/:uuid/subscriptions"

// MetricsListCustomerSubscriptions lists all subscriptions for customer of a given UUID.
//
func (api API) MetricsListCustomerSubscriptions(cursor *Cursor, customerUUID string) (*MetricsCustomerSubscriptions, error) {
	result := &MetricsCustomerSubscriptions{}
	path := strings.Replace(metricsCustomerSubscriptionsEndpoint, ":uuid", customerUUID, 1)
	query := make([]interface{}, 0, 1)
	if cursor != nil {
		query = append(query, *cursor)
	}
	return result, api.list(path, result, query...)
}

// MetricsConnectSubscriptions connects subscription objects for a customer.
//
func (api API) MetricsConnectSubscriptions(dataSourceUUID string, customerUUID string, subscriptions []*MetricsCustomerSubscription) error {
	path := strings.Replace(connectSubscriptionEndpoint, ":uuid", customerUUID, 1)

	refs := make([]MetricsSubscriptionReference, len(subscriptions))
	for i, sub := range subscriptions {
		refs[i] = MetricsSubscriptionReference{
			UUID:           sub.UUID,
			DataSourceUUID: dataSourceUUID,
		}
	}

	return api.merge(path, MetricsConnectSubscriptionsParams{
		Subscriptions: refs,
	})
}

// MetricsDisconnectSubscriptions disconnects subscription objects for a customer.
//
func (api API) MetricsDisconnectSubscriptions(dataSourceUUID string, customerUUID string, subscriptions []*MetricsCustomerSubscription) error {
	path := strings.Replace(disconnectSubscriptionEndpoint, ":uuid", customerUUID, 1)

	refs := make([]MetricsSubscriptionReference, len(subscriptions))
	for i, sub := range subscriptions {
		refs[i] = MetricsSubscriptionReference{
			UUID:           sub.UUID,
			DataSourceUUID: dataSourceUUID,
		}
	}

	return api.merge(path, MetricsConnectSubscriptionsParams{
		Subscriptions: refs,
	})
}
