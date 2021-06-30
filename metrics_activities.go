package chartmogul

// MetricsActivity represents Metrics API activity in ChartMogul.
type MetricsActivity struct {
	Date                   string  `json:"date"`
	ActivityArr            float64 `json:"activity-arr"`
	ActivityMrr            float64 `json:"activity-mrr"`
	ActivityMrrMovement    float64 `json:"activity-mrr-movement"`
	Currency               string  `json:"currency"`
	Description            string  `json:"description"`
	Type                   string  `json:"type"`
	SubscriptionExternalID string  `json:"subscription-external-id"`
	PlanExternalID         string  `json:"plan-external-id"`
	CustomerName           string  `json:"customer-name"`
	CustomerUUID           string  `json:"customer-uuid"`
	CustomerExternalID     string  `json:"customer-external-id"`
	BillingConnectorUUID   string  `json:"billing-connector-uuid"`
	UUID                   string  `json:"uuid"`
}

// MetricsActivities is the result of listing activities in Metrics API.
type MetricsActivities struct {
	Entries []*MetricsActivity `json:"entries"`
	HasMore bool               `json:"has_more"`
	PerPage uint32             `json:"per_page"`
}

type MetricsListActivitiesParams struct {
	Type      string `json:"type,omitempty"`
	StartDate string `json:"start-date,omitempty"`
	EndDate   string `json:"end-date,omitempty"`
	Order     string `json:"order,omitempty"`
	AnchorCursor
}

const metricsActivitiesEndpoint = "activities"

// MetricsListActivities lists all activities for an account
func (api API) MetricsListActivities(listActivitiesParams *MetricsListActivitiesParams) (*MetricsActivities, error) {
	result := &MetricsActivities{}
	query := make([]interface{}, 0, 1)
	if listActivitiesParams != nil {
		query = append(query, *listActivitiesParams)
	}
	return result, api.list(metricsActivitiesEndpoint, result, query...)
}
