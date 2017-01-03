package chartmogul

// DataSource represents Import API data source in ChartMogul.
// See https://dev.chartmogul.com/v1.0/reference#list-data-sources
type DataSource struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
}

// DataSources is the result of listing data sources, but doesn't contain any paging.
type DataSources struct {
	DataSources []*DataSource `json:"data_sources"`
}

// importCreateDataSourceCall represents arguments to be marshalled into JSON.
type importCreateDataSourceCall struct {
	Name string `json:"name"`
}

const (
	dataSourcesEndpoint      = "data_sources"
	singleDataSourceEndpoint = "data_sources/:uuid"
)

// ImportCreateDataSource creates an Import API Data Source in ChartMogul.
//
// See https://dev.chartmogul.com/v1.0/reference#create-data-source
func (api API) ImportCreateDataSource(name string) (*DataSource, error) {
	ds := &DataSource{}
	err := api.create(dataSourcesEndpoint, importCreateDataSourceCall{Name: name}, ds)
	return ds, err
}

// ImportRetrieveDataSource returns one Data Source by UUID.
func (api API) ImportRetrieveDataSource(dataSourceUUID string) (*DataSource, error) {
	result := &DataSource{}
	return result, api.retrieve(singleDataSourceEndpoint, dataSourceUUID, result)
}

// ImportListDataSources lists all available Data Sources (no paging).
//
// See https://dev.chartmogul.com/v1.0/reference#list-data-sources
func (api API) ImportListDataSources() (*DataSources, error) {
	ds := &DataSources{}
	err := api.list(dataSourcesEndpoint, ds)
	return ds, err
}

// ImportDeleteDataSource deletes the data source identified by its UUID.
//
// See https://dev.chartmogul.com/reference#delete-a-data-source
func (api API) ImportDeleteDataSource(uuid string) error {
	return api.delete(singleDataSourceEndpoint, uuid)
}
