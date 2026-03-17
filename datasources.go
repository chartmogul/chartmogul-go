package chartmogul

// ProcessingStatus represents the processing status of a data source.
type ProcessingStatus struct {
	Processed *int `json:"processed,omitempty"`
	Pending   *int `json:"pending,omitempty"`
	Failed    *int `json:"failed,omitempty"`
}

// AutoChurnSubscriptionSetting represents the auto churn subscription setting for a data source.
type AutoChurnSubscriptionSetting struct {
	Enabled  bool `json:"enabled"`
	Interval *int `json:"interval"`
}

// DataSource represents API data source in ChartMogul.
type DataSource struct {
	UUID                         string                        `json:"uuid"`
	Name                         string                        `json:"name"`
	CreatedAt                    string                        `json:"created_at"`
	Status                       string                        `json:"status"`
	System                       string                        `json:"system"`
	ProcessingStatus             *ProcessingStatus             `json:"processing_status,omitempty"`
	AutoChurnSubscriptionSetting *AutoChurnSubscriptionSetting `json:"auto_churn_subscription_setting,omitempty"`
	InvoiceHandlingSetting       map[string]interface{}        `json:"invoice_handling_setting,omitempty"`
	Errors                       Errors                        `json:"errors,omitempty"`
}

// DataSources is the result of listing data sources, but doesn't contain any paging.
type DataSources struct {
	DataSources []*DataSource `json:"data_sources"`
}

// ExtraDataSourceParams are optional parameters for additional data source information.
type ExtraDataSourceParams struct {
	WithProcessingStatus             *bool `json:"with_processing_status,omitempty"`
	WithAutoChurnSubscriptionSetting *bool `json:"with_auto_churn_subscription_setting,omitempty"`
	WithInvoiceHandlingSetting       *bool `json:"with_invoice_handling_setting,omitempty"`
}

// ListDataSourcesParams are optional parameters for listing data sources.
type ListDataSourcesParams struct {
	ExtraDataSourceParams
	Name   string `json:"name,omitempty"`
	System string `json:"system,omitempty"`
}

// createDataSourceCall represents arguments to be marshalled into JSON.
type createDataSourceCall struct {
	Name string `json:"name"`
}

const (
	dataSourcesEndpoint      = "data_sources"
	singleDataSourceEndpoint = "data_sources/:uuid"
	purgeDataSourceEndpoint  = "data_sources/:uuid/dependent"
	emptyDataSourceEndpoint  = "data_sources/:uuid/all"
)

// CreateDataSource creates an API Data Source in ChartMogul.
//
func (api API) CreateDataSource(name string) (*DataSource, error) {
	ds := &DataSource{}
	err := api.create(dataSourcesEndpoint, createDataSourceCall{Name: name}, ds)
	return ds, err
}

// CreateDataSourceWithSystem creates an API Data Source in ChartMogul.
// * Allows other parameters than just the name.
//
func (api API) CreateDataSourceWithSystem(dataSource *DataSource) (*DataSource, error) {
	ds := &DataSource{}
	err := api.create(dataSourcesEndpoint, dataSource, ds)
	return ds, err
}

// RetrieveDataSource returns one Data Source by UUID.
// Optionally accepts ExtraDataSourceParams for additional query options with named parameters:
// - WithProcessingStatus: include processing status information
// - WithAutoChurnSubscriptionSetting: include auto-churn subscription settings
// - WithInvoiceHandlingSetting: include invoice handling settings
//
func (api API) RetrieveDataSource(dataSourceUUID string, params ...*ExtraDataSourceParams) (*DataSource, error) {
	result := &DataSource{}

	if len(params) == 0 || params[0] == nil {
		// No params provided - use original retrieve method
		return result, api.retrieve(singleDataSourceEndpoint, dataSourceUUID, result)
	}

	return result, api.retrieveWithParams(singleDataSourceEndpoint, dataSourceUUID, result, params[0])
}

// ListDataSources lists all available Data Sources (no paging).
// Optionally accepts ExtraDataSourceParams for additional query options with named parameters:
// - WithProcessingStatus: include processing status information
// - WithAutoChurnSubscriptionSetting: include auto-churn subscription settings
// - WithInvoiceHandlingSetting: include invoice handling settings
//
func (api API) ListDataSources(params ...*ExtraDataSourceParams) (*DataSources, error) {
	ds := &DataSources{}
	if len(params) == 0 || params[0] == nil {
		// No params provided - use original list method
		err := api.list(dataSourcesEndpoint, ds)
		return ds, err
	}

	// Use list with query parameters
	query := make([]interface{}, 0, 1)
	query = append(query, *params[0])
	err := api.list(dataSourcesEndpoint, ds, query...)
	return ds, err
}

// ListDataSourcesWithFilters lists all available Data Sources (no paging).
// Accepts filtering parameters and extra data parameters in a single struct.
//
func (api API) ListDataSourcesWithFilters(listDataSourcesParams *ListDataSourcesParams) (*DataSources, error) {
	ds := &DataSources{}
	query := make([]interface{}, 0, 1)
	if listDataSourcesParams != nil {
		query = append(query, *listDataSourcesParams)
	}
	err := api.list(dataSourcesEndpoint, ds, query...)
	return ds, err
}

// DeleteDataSource deletes the data source identified by its UUID.
//
func (api API) DeleteDataSource(uuid string) error {
	return api.delete(singleDataSourceEndpoint, uuid)
}

// PurgeDataSource deletes all the data except the data source itself and the customers
//
func (api API) PurgeDataSource(dataSourceUUID string) error {
	return api.delete(purgeDataSourceEndpoint, dataSourceUUID)
}

// EmptyDataSource deletes all the data in the data source, but keeps the UUID.
func (api API) EmptyDataSource(dataSourceUUID string) error {
	return api.delete(emptyDataSourceEndpoint, dataSourceUUID)
}
