package chartmogul

const (
	planGroupsEndpoint      = "plan_groups"
	singlePlanGroupEndpoint = "plan_groups/:uuid"
)

// PlanGroup represents groups of plans in ChartMogul
type PlanGroup struct {
	UUID       string    `json:"uuid,omitempty"`
	Name       string    `json:"name,omitempty"`
	Plans      []*string `json:"plans,omitempty"`
	PlansCount int       `json:"plans_count,omitempty"`
	Errors     Errors    `json:"errors,omitempty"`
}

// PlanGroups is result of listing: plan_groups + page.
type PlanGroups struct {
	PlanGroups []*PlanGroup `json:"plan_groups"`
	Pagination
}

// CreatePlanGroup creates plan group with given name and plans.
//
// See https://dev.chartmogul.com/v1.0/reference#plan_groups
func (api API) CreatePlanGroup(planGroup *PlanGroup) (result *PlanGroup, err error) {
	result = &PlanGroup{}
	return result, api.create(planGroupsEndpoint, planGroup, result)
}

// RetrievePlanGroup returns one plan group by UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#plan_groups
func (api API) RetrievePlanGroup(planGroupUUID string) (*PlanGroup, error) {
	result := &PlanGroup{}
	return result, api.retrieve(singlePlanGroupEndpoint, planGroupUUID, result)
}

// ListPlanGroups returns list of plan groups.
//
// See https://dev.chartmogul.com/v1.0/reference#plan_groups
func (api API) ListPlanGroups(cursor *Cursor) (*PlanGroups, error) {
	result := &PlanGroups{}
	query := make([]interface{}, 0, 1)
	if cursor != nil {
		query = append(query, *cursor)
	}
	return result, api.list(planGroupsEndpoint, result, query...)
}

// UpdatePlanGroup updates a name or plans.
//
// See https://dev.chartmogul.com/v1.0/reference#plan_groups
func (api API) UpdatePlanGroup(planGroup *PlanGroup, planGroupUUID string) (*PlanGroup, error) {
	result := &PlanGroup{}
	return result, api.update(singlePlanGroupEndpoint, planGroupUUID, planGroup, result)
}

// DeletePlanGroup deletes one plan group by UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#plan_groups
func (api API) DeletePlanGroup(planGroupUUID string) error {
	return api.delete(singlePlanGroupEndpoint, planGroupUUID)
}
