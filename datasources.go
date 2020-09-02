package chartmogul

// DataSource represents API data source in ChartMogul.
// See https://dev.chartmogul.com/v1.0/reference#list-data-sources
type DataSource struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	System    string `json:"system"`
	Errors    Errors `json:"errors,omitempty"`
}

// DataSources is the result of listing data sources, but doesn't contain any paging.
type DataSources struct {
	DataSources []*DataSource `json:"data_sources"`
}

// ListDataSourcesParams are optional parameters for listing data sources.
type ListDataSourcesParams struct {
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
// See https://dev.chartmogul.com/v1.0/reference#data-sources
func (api API) CreateDataSource(name string) (*DataSource, error) {
	ds := &DataSource{}
	err := api.create(dataSourcesEndpoint, createDataSourceCall{Name: name}, ds)
	return ds, err
}

// CreateDataSourceWithSystem creates an API Data Source in ChartMogul.
// * Allows other parameters than just the name.
//
// See https://dev.chartmogul.com/v1.0/reference#data-sources
func (api API) CreateDataSourceWithSystem(dataSource *DataSource) (*DataSource, error) {
	ds := &DataSource{}
	err := api.create(dataSourcesEndpoint, dataSource, ds)
	return ds, err
}

// RetrieveDataSource returns one Data Source by UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#data-sources
func (api API) RetrieveDataSource(dataSourceUUID string) (*DataSource, error) {
	result := &DataSource{}
	return result, api.retrieve(singleDataSourceEndpoint, dataSourceUUID, result)
}

// ListDataSources lists all available Data Sources (no paging).
//
// See https://dev.chartmogul.com/v1.0/reference#data-sources
func (api API) ListDataSources() (*DataSources, error) {
	ds := &DataSources{}
	err := api.list(dataSourcesEndpoint, ds)
	return ds, err
}

// ListDataSourcesWithFilters lists all available Data Sources (no paging).
// * Allows filtering.
//
// See https://dev.chartmogul.com/v1.0/reference#data-sources
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
// See https://dev.chartmogul.com/v1.0/reference#data-sources
func (api API) DeleteDataSource(uuid string) error {
	return api.delete(singleDataSourceEndpoint, uuid)
}

// PurgeDataSource deletes all the data except the data source itself and the customers
//
// See https://dev.chartmogul.com/v1.0/reference#data-sources
func (api API) PurgeDataSource(dataSourceUUID string) error {
	return api.delete(purgeDataSourceEndpoint, dataSourceUUID)
}

// EmptyDataSource deletes all the data in the data source, but keeps the UUID.
func (api API) EmptyDataSource(dataSourceUUID string) error {
	return api.delete(emptyDataSourceEndpoint, dataSourceUUID)
}
