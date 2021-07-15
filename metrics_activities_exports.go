package chartmogul

// MetricsActivitiesExport represents Metrics API activity export in ChartMogul.
type MetricsActivitiesExport struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	FileURL   string `json:"file_url"`
	Params    Params `json:"params"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}

// Params provides information on the requested export.
type Params struct {
	Kind   string       `json:"kind"`
	Params NestedParams `json:"params,omitempty"`
}

// NestedParams represents the params of the requested type of export.
type NestedParams struct {
	ActivityType string `json:"activity_type,omitempty"`
	StartDate    string `json:"start_date,omitempty"`
	EndDate      string `json:"end_date,omitempty"`
}

// CreateMetricsActivitiesExportParam to create a MetricsActivitiesExport.
type CreateMetricsActivitiesExportParam struct {
	Type      string `json:"type,omitempty"`
	StartDate string `json:"start-date,omitempty"`
	EndDate   string `json:"end-date,omitempty"`
}

const (
	metricsActivitiesExportEndpoint       = "activities_export"
	singleMetricsActivitiesExportEndpoint = "activities_export/:uuid"
)

// MetricsCreateActivitiesExport requests creation of an activities export in Chartmogul.
//
// See https://dev.chartmogul.com/v1.0/reference#activities_export
func (api API) MetricsCreateActivitiesExport(CreateMetricsActivitiesExportParam *CreateMetricsActivitiesExportParam) (*MetricsActivitiesExport, error) {
	result := &MetricsActivitiesExport{}
	return result, api.create(metricsActivitiesExportEndpoint, CreateMetricsActivitiesExportParam, result)
}

// MetricsRetrieveActivitiesExport returns one activities export as in API.
//
// See https://dev.chartmogul.com/v1.0/reference#activities_export
func (api API) MetricsRetrieveActivitiesExport(activitiesExportUUID string) (*MetricsActivitiesExport, error) {
	result := &MetricsActivitiesExport{}
	return result, api.retrieve(singleMetricsActivitiesExportEndpoint, activitiesExportUUID, result)
}
