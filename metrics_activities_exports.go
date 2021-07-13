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

type Params struct {
	Kind   string       `json:"kind"`
	Params NestedParams `json:"params,omitempty"`
}

type NestedParams struct {
	ActivityType string `json:"activity_type,omitempty"`
	StartDate    string `json:"start-date,omitempty"`
	EndDate      string `json:"end-date,omitempty"`
}

// NewMetricsActivitiesExport is the POST-ed to create a MetricsActivitiesExport .
type NewMetricsActivitiesExport struct {
	Type      string `json:"type,omitempty"`
	StartDate string `json:"start-date,omitempty"`
	EndDate   string `json:"end-date,omitempty"`
}

const (
	metricsActivitiesExportEndpoint       = "activities_export"
	singleMetricsActivitiesExportEndpoint = "activities_export/:activities_export_uuid"
)

// MetricsCreateActivitiesExport requests creation of an activities export in Chartmogul.
//
// See https://dev.chartmogul.com/v1.0/reference#activities_export
func (api API) MetricsCreateActivitiesExport(NewMetricsActivitiesExport *NewMetricsActivitiesExport) (*MetricsActivitiesExport, error) {
	result := &MetricsActivitiesExport{}
	return result, api.create(metricsActivitiesExportEndpoint, NewMetricsActivitiesExport, result)
}

// MetricsRetrieveActivitiesExport returns one activities export as in API.
//
// See https://dev.chartmogul.com/v1.0/reference#activities_export
func (api API) MetricsRetrieveActivitiesExport(activitiesExportUUID string) (*MetricsActivitiesExport, error) {
	result := &MetricsActivitiesExport{}
	return result, api.retrieve(singleMetricsActivitiesExportEndpoint, activitiesExportUUID, result)
}
