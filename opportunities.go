package chartmogul

type Opportunity struct {
	UUID string `json:"uuid"`
	// Basic info
	CustomerUUID       string `json:"customer_uuid"`
	Owner              string `json:"owner"`
	Pipeline           string `json:"pipeline"`
	PipelineStage      string `json:"pipeline_stage"`
	EstimatedCloseDate string `json:"estimated_close_date"`
	Currency           string `json:"currency"`
	AmountInCents      int    `json:"amount_in_cents"`
	Type               string `json:"type"`
	ForecastCategory   string `json:"forecast_category"`
	WinLikelihood      int    `json:"win_likelihood"`
	Custom             map[string]interface{}  `json:"custom,omitempty"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type UpdateOpportunity struct {
	Owner              string `json:"owner,omitempty"`
	Pipeline           string `json:"pipeline,omitempty"`
	PipelineStage      string `json:"pipeline_stage,omitempty"`
	EstimatedCloseDate string `json:"estimated_close_date,omitempty"`
	Currency           string `json:"currency,omitempty"`
	AmountInCents      int    `json:"amount_in_cents,omitempty"`
	Type               string `json:"type,omitempty"`
	ForecastCategory   string `json:"forecast_category,omitempty"`
	WinLikelihood      int    `json:"win_likelihood,omitempty"`
	Custom             []Custom `json:"custom,omitempty"`
}

type NewOpportunity struct {
	// Obligatory
	CustomerUUID       string `json:"customer_uuid"`
	Owner              string `json:"owner"`
	Pipeline           string `json:"pipeline"`
	PipelineStage      string `json:"pipeline_stage"`
	EstimatedCloseDate string `json:"estimated_close_date"`
	Currency           string `json:"currency"`
	AmountInCents      int    `json:"amount_in_cents"`

	//Optional
	Type               string `json:"type,omitempty"`
	ForecastCategory   string `json:"forecast_category,omitempty"`
	WinLikelihood      int    `json:"win_likelihood,omitempty"`
	Custom             []Custom `json:"custom,omitempty"`
}

// ListOpportunitiesParams = parameters for listing customer opportunities in API.
type ListOpportunitiesParams struct {
	CustomerUUID string `json:"customer_uuid"`
	Cursor
}

// Opportunities is result of listing opportunities in API.
type Opportunities struct {
	Entries []*Opportunity `json:"entries"`
	Pagination
}

const (
	singleOpportunityEndpoint = "opportunities/:uuid"
	opportunitiesEndpoint      = "opportunities"
)

// CreateOpportunity create the opportunity to Chartmogul
//
// See https://dev.chartmogul.com/reference/create-an-opportunity
func (api API) CreateOpportunity(input *NewOpportunity) (*Opportunity, error) {
	result := &Opportunity{}
	return result, api.create(opportunitiesEndpoint, input, result)
}

// RetrieveOpportunity returns one opportunity as in API.
//
// See https://dev.chartmogul.com/reference/retrieve-an-opportunity
func (api API) RetrieveOpportunity(opportunityUUID string) (*Opportunity, error) {
	result := &Opportunity{}
	return result, api.retrieve(singleOpportunityEndpoint, opportunityUUID, result)
}

// UpdateOpportunity updates one opportunity in API.
//
// See https://dev.chartmogul.com/reference/update-an-opportunity
func (api API) UpdateOpportunity(input *UpdateOpportunity, opportunityUUID string) (*Opportunity, error) {
	output := &Opportunity{}
	return output, api.update(singleOpportunityEndpoint, opportunityUUID, input, output)
}

// ListOpportunities lists all opportunities.
//
// See https://dev.chartmogul.com/reference/list-opportunities
func (api API) ListOpportunities(listOpportunitiesParams *ListOpportunitiesParams) (*Opportunities, error) {
	result := &Opportunities{}
	query := make([]interface{}, 0, 1)
	if listOpportunitiesParams != nil {
		query = append(query, *listOpportunitiesParams)
	}
	return result, api.list(opportunitiesEndpoint, result, query...)
}

// DeleteOpportunity deletes one opportunity by UUID.
//
// See https://dev.chartmogul.com/reference/delete-an-opportunity
func (api API) DeleteOpportunity(opportunityUUID string) error {
	return api.delete(singleOpportunityEndpoint, opportunityUUID)
}
