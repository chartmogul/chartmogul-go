package chartmogul

const (
	plansEndpoint      = "plans"
	singlePlanEndpoint = "plans/:uuid"
)

// Plan represents ChartMogul categorization of subscriptions.
type Plan struct {
	UUID           string `json:"uuid,omitempty"`
	DataSourceUUID string `json:"data_source_uuid,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
	Name           string `json:"name,omitempty"`
	IntervalCount  uint32 `json:"interval_count,omitempty"`
	IntervalUnit   string `json:"interval_unit,omitempty"`
	Errors         Errors `json:"errors,omitempty"`
}

// Plans is result of listing: plans + paging.
type Plans struct {
	Plans       []*Plan `json:"plans"`
	TotalPages  uint32  `json:"total_pages"`
	CurrentPage uint32  `json:"current_page"`
	Pagination
}

// ListPlansParams = optional parameters for listing plans.
type ListPlansParams struct {
	DataSourceUUID string `json:"data_source_uuid"`
	ExternalID     string `json:"external_id,omitempty"`
	System         string `json:"system,omitempty"`
	Cursor
	PaginationWithCursor
}

// CreatePlan creates plan under given Data Source.
//
// See https://dev.chartmogul.com/v1.0/reference#plans
func (api API) CreatePlan(plan *Plan) (result *Plan, err error) {
	result = &Plan{}
	return result, api.create(plansEndpoint, plan, result)
}

// RetrievePlan returns one plan by UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#plans
func (api API) RetrievePlan(planUUID string) (*Plan, error) {
	result := &Plan{}
	return result, api.retrieve(singlePlanEndpoint, planUUID, result)
}

// ListPlans returns list of plans.
//
// See https://dev.chartmogul.com/v1.0/reference#plans
func (api API) ListPlans(listPlansParams *ListPlansParams) (*Plans, error) {
	result := &Plans{}
	query := make([]interface{}, 0, 1)
	if listPlansParams != nil {
		query = append(query, *listPlansParams)
	}
	return result, api.list(plansEndpoint, result, query...)
}

// UpdatePlan returns list of plans.
//
// See https://dev.chartmogul.com/v1.0/reference#plans
func (api API) UpdatePlan(plan *Plan, planUUID string) (*Plan, error) {
	result := &Plan{}
	return result, api.update(singlePlanEndpoint, planUUID, plan, result)
}

// DeletePlan deletes one plan by UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#plans
func (api API) DeletePlan(planUUID string) error {
	return api.delete(singlePlanEndpoint, planUUID)
}
