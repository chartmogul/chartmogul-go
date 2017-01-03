package chartmogul

const (
	importPlansEndpoint      = "plans"
	importSinglePlanEndpoint = "plans/:uuid"
)

// Plan represents ChartMogul categorization of subscriptions.
type Plan struct {
	UUID           string `json:"uuid,omitempty"`
	DataSourceUUID string `json:"data_source_uuid,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
	Name           string `json:"name,omitempty"`
	IntervalCount  uint32 `json:"interval_count,omitempty"`
	IntervalUnit   string `json:"interval_unit,omitempty"`
}

// Plans is result of listing: plans + paging.
type Plans struct {
	Plans       []*Plan `json:"plans"`
	TotalPages  uint32  `json:"total_pages"`
	CurrentPage uint32  `json:"current_page"`
}

// ListPlansParams = optional parameters for listing plans.
type ListPlansParams struct {
	DataSourceUUID string `json:"data_source_uuid"`
	ExternalID     string `json:"external_id,omitempty"`
	Cursor
}

// ImportCreatePlan creates plan under given Data Source.
//
// See https://dev.chartmogul.com/v1.0/reference#import-plan
func (api API) ImportCreatePlan(plan *Plan, dataSourceUUID string) (result *Plan, err error) {
	plan.DataSourceUUID = dataSourceUUID
	result = &Plan{}
	return result, api.create(importPlansEndpoint, plan, result)
}

// ImportRetrievePlan returns one plan by UUID.
func (api API) ImportRetrievePlan(planUUID string) (*Plan, error) {
	result := &Plan{}
	return result, api.retrieve(importSinglePlanEndpoint, planUUID, result)
}

// ImportListPlans returns list of plans.
//
// https://dev.chartmogul.com/reference#list-all-imported-plans
func (api API) ImportListPlans(listPlansParams *ListPlansParams) (*Plans, error) {
	result := &Plans{}
	return result, api.list(importPlansEndpoint, result, *listPlansParams)
}

// ImportUpdatePlan returns list of plans.
//
// https://dev.chartmogul.com/reference#list-all-imported-plans
func (api API) ImportUpdatePlan(plan *Plan, planUUID string) (*Plan, error) {
	result := &Plan{}
	return result, api.update(importSinglePlanEndpoint, planUUID, plan, result)
}

// ImportDeletePlan deletes one plan by UUID.
func (api API) ImportDeletePlan(planUUID string) error {
	return api.delete(importSinglePlanEndpoint, planUUID)
}
