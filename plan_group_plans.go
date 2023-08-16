package chartmogul

import (
	"strings"
)

const (
	planGroupPlansEndpoint = "plan_groups/:uuid/plans"
)

// Plan groups represents groups of plans in ChartMogul

// PlanGroupPlans is result of listing: plans + paging for a plan group.
type PlanGroupPlans struct {
	Plans []*Plan `json:"plans"`
	Pagination
}

// ListPlanGroupPlans returns list of plans in with a plan group given the plan group uuid.
//
// See https://dev.chartmogul.com/v1.0/reference#plan_groups
func (api API) ListPlanGroupPlans(cursor *PaginationWithCursor, planGroupUUID string) (*PlanGroupPlans, error) {
	result := &PlanGroupPlans{}
	query := make([]interface{}, 0, 1)
	if cursor != nil {
		query = append(query, *cursor)
	}
	path := strings.Replace(planGroupPlansEndpoint, ":uuid", planGroupUUID, 1)
	return result, api.list(path, result, query...)
}
