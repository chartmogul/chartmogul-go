package chartmogul

// MetricsFilter convenient object to hold all filtering parameters.
type MetricsFilter struct {
	StartDate string `json:"start-date,omitempty"`
	EndDate   string `json:"end-date,omitempty"`
	Interval  string `json:"interval,omitempty"`
	Geo       string `json:"geo,omitempty"`
	Plans     string `json:"plans,omitempty"`
}

// AllMetrics represents results of Metrics API.
type AllMetrics struct {
	Date                                string  `json:"date"`
	CustomerChurnRate                   float64 `json:"customer-churn-rate"`
	MrrChurnRate                        float64 `json:"mrr-churn-rate"`
	Ltv                                 float64 `json:"ltv"`
	Customers                           uint32  `json:"customers"`
	Asp                                 float64 `json:"asp"`
	Arpa                                float64 `json:"arpa"`
	Arr                                 float64 `json:"arr"`
	Mrr                                 float64 `json:"mrr"`
	CustomerChurnRatePercentageChange   float64 `json:"customer-churn-rate-percentage-change"`
	MrrChurnRatePercentageChange        float64 `json:"mrr-churn-rate-percentage-change"`
	LtvPercentageChange                 float64 `json:"ltv-percentage-change"`
	CustomersPercentageChange           float64 `json:"customers-percentage-change"`
	AspPercentageChange                 float64 `json:"asp-percentage-change"`
	ArpaPercentageChange                float64 `json:"arpa-percentage-change"`
	ArrPercentageChange                 float64 `json:"arr-percentage-change"`
	MrrPercentageChange                 float64 `json:"mrr-percentage-change"`
}

// MetricsResult represents results of Metrics API.
type MetricsResult struct {
	Entries []*AllMetrics `json:"entries,omitempty"`
}

// Summary represents results of Metrics API.
type Summary struct {
	Current          float64 `json:"current"`
	Previous         float64 `json:"previous"`
	PercentageChange float64 `json:"percentage-change"`
}

// Summary represents results of Metrics API.
type AllSummary struct {
	CurrentCustomerChurnRate            float64 `json:"current-customer-churn-rate"`
	PreviousCustomerChurnRate           float64 `json:"previous-customer-churn-rate"`
	CustomerChurnRatePercentageChange   float64 `json:"customer-churn-rate-percentage-change"`
	CurrentMrrChurnRate                 float64 `json:"current-mrr-churn-rate"`
    PreviousMrrChurnRate                float64 `json:"previous-mrr-churn-rate"`
    MrrChurnRatePercentageChange        float64 `json:"mrr-churn-rate-percentage-change"`
    CurrentLtv                          float64 `json:"current-ltv"`
    PreviousLtv                         float64 `json:"previous-ltv"`
    LtvPercentageChange                 float64 `json:"ltv-percentage-change"`
    CurrentCustomers                    float64 `json:"current-customers"`
    PreviousCustomers                   float64 `json:"previous-customers"`
    CustomersPercentageChange           float64 `json:"customers-percentage-change"`
    CurrentAsp                          float64 `json:"current-asp"`
    PreviousAsp                         float64 `json:"previous-asp"`
    AspPercentageChange                 float64 `json:"asp-percentage-change"`
    CurrentArpa                         float64 `json:"current-arpa"`
    PreviousArpa                        float64 `json:"previous-arpa"`
    ArpaPercentageChange                float64 `json:"arpa-percentage-change"`
    CurrentArr                          float64 `json:"current-arr"`
    PreviousArr                         float64 `json:"previous-arr"`
    ArrPercentageChange                 float64 `json:"arr-percentage-change"`
    CurrentMrr                          float64 `json:"current-mrr"`
    PreviousMrr                         float64 `json:"previous-mrr"`
    MrrPercentageChange                 float64 `json:"mrr-percentage-change"`
}

// MRRMetrics represents results of Metrics API.
type MRRMetrics struct {
	Date            string  `json:"date"`
	MRR             float64 `json:"mrr"`
	MRRNewBusiness  float64 `json:"mrr-new-business"`
	MRRExpansion    float64 `json:"mrr-expansion"`
	MRRContraction  float64 `json:"mrr-contraction"`
	MRRChurn        float64 `json:"mrr-churn"`
	MRRReactivation float64 `json:"mrr-reactivation"`
}

// MRRResult represents results of Metrics API.
type MRRResult struct {
	Entries []*MRRMetrics `json:"entries,omitempty"`
	Summary *Summary      `json:"summary"`
}

// ARRMetrics represents results of Metrics API.
type ARRMetrics struct {
	Date                string `json:"date"`
	ARR                 float64 `json:"arr"`
	PercentageChange    float64 `json:"percentage-change"`
}

// ARRResult represents results of Metrics API.
type ARRResult struct {
	Entries []*ARRMetrics `json:"entries,omitempty"`
	Summary *Summary      `json:"summary"`
}

// ARPAMetrics represents results of Metrics API.
type ARPAMetrics struct {
	Date                string `json:"date"`
	ARPA                float64 `json:"arpa"`
	PercentageChange    float64 `json:"percentage-change"`
}

// ARPAResult represents results of Metrics API.
type ARPAResult struct {
	Entries []*ARPAMetrics `json:"entries,omitempty"`
	Summary *Summary       `json:"summary"`
}

// ASPMetrics represents results of Metrics API.
type ASPMetrics struct {
	Date                string `json:"date"`
	ASP                 float64 `json:"asp"`
	PercentageChange    float64 `json:"percentage-change"`
}

// ASPResult represents results of Metrics API.
type ASPResult struct {
	Entries []*ASPMetrics `json:"entries,omitempty"`
	Summary *Summary      `json:"summary"`
}

// CustomerCountMetrics represents results of Metrics API.
type CustomerCountMetrics struct {
	Date                string `json:"date"`
	Customers           uint32 `json:"customers"`
	PercentageChange    float64 `json:"percentage-change"`
}

// CustomerCountResult represents results of Metrics API.
type CustomerCountResult struct {
	Entries []*CustomerCountMetrics `json:"entries,omitempty"`
	Summary *Summary                `json:"summary"`
}

// CustomerChurnRateMetrics represents results of Metrics API.
type CustomerChurnRateMetrics struct {
	Date                string  `json:"date"`
	CustomerChurnRate   float64 `json:"customer-churn-rate"`
	PercentageChange    float64 `json:"percentage-change"`
}

// CustomerChurnRateResult represents results of Metrics API.
type CustomerChurnRateResult struct {
	Entries []*CustomerChurnRateMetrics `json:"entries,omitempty"`
	Summary *Summary                    `json:"summary"`
}

// MRRChurnRateMetrics represents results of Metrics API.
type MRRChurnRateMetrics struct {
	Date                string  `json:"date"`
	MRRChurnRate        float64 `json:"mrr-churn-rate"`
	PercentageChange    float64 `json:"percentage-change"`
}

// MRRChurnRateResult represents results of Metrics API.
type MRRChurnRateResult struct {
	Entries []*MRRChurnRateMetrics `json:"entries,omitempty"`
	Summary *Summary               `json:"summary"`
}

// LTVMetrics represents results of Metrics API.
type LTVMetrics struct {
	Date                string  `json:"date"`
	LTV                 float64 `json:"ltv"`
	PercentageChange    float64 `json:"percentage-change"`
}

// LTVResult represents results of Metrics API.
type LTVResult struct {
	Entries []*LTVMetrics `json:"entries,omitempty"`
	Summary *Summary      `json:"summary"`
}

const (
	metricsEndpoint                  = "metrics/all"
	metricsMRREndpoint               = "metrics/mrr"
	metricsARREndpoint               = "metrics/arr"
	metricsARPAEndpoint              = "metrics/arpa"
	metricsASPEndpoint               = "metrics/asp"
	metricsCustomerCountEndpoint     = "metrics/customer-count"
	metricsCustomerChurnRateEndpoint = "metrics/customer-churn-rate"
	metricsMRRChurnRateEndpoint      = "metrics/mrr-churn-rate"
	metricsLTVEndpoint               = "metrics/ltv"
)

// MetricsRetrieveAll retrieves all key metrics, for the specified time period.
//
// See https://dev.chartmogul.com/v1.0/reference#retrieve-all-key-metrics
func (api API) MetricsRetrieveAll(metricsFilter *MetricsFilter) (*MetricsResult, error) {
	output := &MetricsResult{}
	err := api.list(metricsEndpoint, output, *metricsFilter)
	return output, err
}

// MetricsRetrieveMRR retrieves the MRR metrics, for the specified time period.
//
// See https://dev.chartmogul.com/v1.0/reference#retrieve-mrr
func (api API) MetricsRetrieveMRR(metricsFilter *MetricsFilter) (*MRRResult, error) {
	output := &MRRResult{}
	err := api.list(metricsMRREndpoint, output, *metricsFilter)
	return output, err
}

// MetricsRetrieveARR retrieves the ARR metrics, for the specified time period.
//
// See https://dev.chartmogul.com/v1.0/reference#retrieve-arr
func (api API) MetricsRetrieveARR(metricsFilter *MetricsFilter) (*ARRResult, error) {
	output := &ARRResult{}
	err := api.list(metricsARREndpoint, output, *metricsFilter)
	return output, err
}

// MetricsRetrieveARPA retrieves the ARPA metrics, for the specified time period.
//
// See https://dev.chartmogul.com/v1.0/reference#retrieve-arpa
func (api API) MetricsRetrieveARPA(metricsFilter *MetricsFilter) (*ARPAResult, error) {
	output := &ARPAResult{}
	err := api.list(metricsARPAEndpoint, output, *metricsFilter)
	return output, err
}

// MetricsRetrieveASP retrieves the ASP metrics, for the specified time period.
//
// See https://dev.chartmogul.com/v1.0/reference#retrieve-asp
func (api API) MetricsRetrieveASP(metricsFilter *MetricsFilter) (*ASPResult, error) {
	output := &ASPResult{}
	err := api.list(metricsASPEndpoint, output, *metricsFilter)
	return output, err
}

// MetricsRetrieveCustomerCount retrieves customer count, for the specified time period.
//
// See https://dev.chartmogul.com/v1.0/reference#retrieve-customer-count
func (api API) MetricsRetrieveCustomerCount(metricsFilter *MetricsFilter) (*CustomerCountResult, error) {
	output := &CustomerCountResult{}
	err := api.list(metricsCustomerCountEndpoint, output, *metricsFilter)
	return output, err
}

// MetricsRetrieveCustomerChurnRate retrieves customer churn rate, for the specified time period.
//
// See https://dev.chartmogul.com/v1.0/reference#retrieve-customer-churn-rate
func (api API) MetricsRetrieveCustomerChurnRate(metricsFilter *MetricsFilter) (*CustomerChurnRateResult, error) {
	output := &CustomerChurnRateResult{}
	err := api.list(metricsCustomerChurnRateEndpoint, output, *metricsFilter)
	return output, err
}

// MetricsRetrieveMRRChurnRate retrieves all key metrics, for the specified time period.
//
// See https://dev.chartmogul.com/v1.0/reference#retrieve-mrr-churn-rate
func (api API) MetricsRetrieveMRRChurnRate(metricsFilter *MetricsFilter) (*MRRChurnRateResult, error) {
	output := &MRRChurnRateResult{}
	err := api.list(metricsMRRChurnRateEndpoint, output, *metricsFilter)
	return output, err
}

// MetricsRetrieveLTV retrieves LTV metrics, for the specified time period.
//
// See https://dev.chartmogul.com/v1.0/reference#retrieve-ltv
func (api API) MetricsRetrieveLTV(metricsFilter *MetricsFilter) (*LTVResult, error) {
	output := &LTVResult{}
	err := api.list(metricsLTVEndpoint, output, *metricsFilter)
	return output, err
}
