package chartmogul

import "strings"

// MetricsCustomerActivity represents Metrics API activity in ChartMogul.
type MetricsCustomerActivity struct {
	ID                  uint64  `json:"id"`
	Date                string  `json:"date"`
	ActivityArr         float64 `json:"activity-arr"`
	ActivityMrr         float64 `json:"activity-mrr"`
	ActivityMrrMovement float64 `json:"activity-mrr-movement"`
	Currency            string  `json:"currency"`
	CurrencySign        string  `json:"currency-sign"`
	Description         string  `json:"description"`
	Type                string  `json:"type"`
}

// MetricsCustomerActivities is the result of listing activities in Metrics API.
type MetricsCustomerActivities struct {
	Entries []*MetricsCustomerActivity `json:"entries"`
	HasMore bool                       `json:"has_more"`
	PerPage uint32                     `json:"per_page"`
	Page    uint32                     `json:"page"`
}

const metricsCustomerActivitiesEndpoint = "customers/:uuid/activities"

// MetricsListCustomerActivities lists all activities for cutomer of a given UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#list-customer-activities
func (api API) MetricsListCustomerActivities(cursor *Cursor, customerUUID string) (*MetricsCustomerActivities, error) {
	result := &MetricsCustomerActivities{}
	path := strings.Replace(metricsCustomerActivitiesEndpoint, ":uuid", customerUUID, 1)
	query := make([]interface{}, 0, 1)
	if cursor != nil {
		query = append(query, *cursor)
	}
	return result, api.list(path, result, query...)
}
