package chartmogul

import (
	"strings"
)

const subscriptionsEndpoint = "import/customers/:customerUUID/subscriptions"
const cancelSubscriptionEndpoint = "import/subscriptions/:uuid"
const connectSubscriptionEndpoint = "customers/:uuid/connect_subscriptions"

// Subscription represents Import API subscription in ChartMogul.
type Subscription struct {
	UUID              string   `json:"uuid,omitempty"`
	ExternalID        string   `json:"external_id"`
	PlanUUID          string   `json:"plan_uuid,omitempty"`
	CustomerUUID      string   `json:"customer_uuid,omitempty"`
	DataSourceUUID    string   `json:"data_source_uuid"`
	CancellationDates []string `json:"cancellation_dates,omitempty"`
}

// Subscriptions is the result of listing subscriptions with paging.
type Subscriptions struct {
	Subscriptions []Subscription `json:"subscriptions"`
	CustomerUUID  string         `json:"customer_uuid,omitempty"`
	TotalPages    uint32         `json:"total_pages,omitempty"`
	CurrentPage   uint32         `json:"current_page,omitempty"`
}

// CancelSubscriptionParams represents arguments to be marshalled into JSON.
type CancelSubscriptionParams struct {
	CancelledAt       string    `json:"cancelled_at,omitempty"`
	CancellationDates *[]string `json:"cancellation_dates,omitempty"`
}

// CancelSubscription creates an Import API Data Source in ChartMogul.
//
// See https://dev.chartmogul.com/v1.0/reference#subscriptions
func (api API) CancelSubscription(subscriptionUUID string, cancelSubscriptionParams *CancelSubscriptionParams) (*Subscription, error) {
	result := &Subscription{}
	return result, api.update(cancelSubscriptionEndpoint,
		subscriptionUUID,
		*cancelSubscriptionParams,
		result)
}

// ListSubscriptions lists all subscriptions for cutomer of given UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#subscriptions
func (api API) ListSubscriptions(cursor *Cursor, customerUUID string) (*Subscriptions, error) {
	result := &Subscriptions{}
	path := strings.Replace(subscriptionsEndpoint, ":customerUUID", customerUUID, 1)
	query := make([]interface{}, 0, 1)
	if cursor != nil {
		query = append(query, *cursor)
	}
	return result, api.list(path, result, query...)
}

// ConnectSubscriptions connects two subscription objects
//
// See https://dev.chartmogul.com/reference#connect-subscriptions
func (api API) ConnectSubscriptions(customerUUID string, subscriptions []Subscription) error {
	path := strings.Replace(connectSubscriptionEndpoint, ":uuid", customerUUID, 1)
	return api.merge(path, Subscriptions{
		Subscriptions: subscriptions,
	})
}
