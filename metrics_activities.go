package chartmogul

import "strings"

// MetricsActivity represents Metrics API activity in ChartMogul.
type MetricsActivity struct {
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

// MetricsActivities is the result of listing activities in Metrics API.
type MetricsActivities struct {
	Entries []*MetricsActivity `json:"entries"`
	HasMore bool               `json:"has_more"`
	PerPage uint32             `json:"per_page"`
	Page    uint32             `json:"page"`
}

const metricsActivitiesEndpoint = "customers/:uuid/activities"

// MetricsListActivities lists all activities for cutomer of a given UUID.
//
// See https://dev.chartmogul.com/reference#list-customer-activities
func (api API) MetricsListActivities(cursor *Cursor, customerUUID string) (*MetricsActivities, error) {
	result := &MetricsActivities{}
	path := strings.Replace(metricsActivitiesEndpoint, ":uuid", customerUUID, 1)
	query := make([]interface{}, 0, 1)
	if cursor != nil {
		query = append(query, *cursor)
	}
	return result, api.list(path, result, query...)
}
