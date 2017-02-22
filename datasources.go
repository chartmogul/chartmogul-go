package chartmogul

// DataSource represents API data source in ChartMogul.
// See https://dev.chartmogul.com/v1.0/reference#list-data-sources
type DataSource struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	System    string `json:"system"`
}

// DataSources is the result of listing data sources, but doesn't contain any paging.
type DataSources struct {
	DataSources []*DataSource `json:"data_sources"`
}

// createDataSourceCall represents arguments to be marshalled into JSON.
type createDataSourceCall struct {
	Name string `json:"name"`
}

const (
	dataSourcesEndpoint      = "data_sources"
	singleDataSourceEndpoint = "data_sources/:uuid"
)

// CreateDataSource creates an API Data Source in ChartMogul.
//
// See https://dev.chartmogul.com/v1.0/reference#data-sources
func (api API) CreateDataSource(name string) (*DataSource, error) {
	ds := &DataSource{}
	err := api.create(dataSourcesEndpoint, createDataSourceCall{Name: name}, ds)
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

// DeleteDataSource deletes the data source identified by its UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#data-sources
func (api API) DeleteDataSource(uuid string) error {
	return api.delete(singleDataSourceEndpoint, uuid)
}
