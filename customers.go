package chartmogul

// Customer represents one customer from Import API.
// See https://dev.chartmogul.com/v1.0/reference#import-customer
type Customer struct {
	UUID               string  `json:"uuid,omitempty"`
	DataSourceUUID     string  `json:"data_source_uuid"`
	ExternalID         string  `json:"external_id"`
	Name               string  `json:"name"`
	Email              string  `json:"email,omitempty"`
	Company            string  `json:"company,omitempty"`
	Country            string  `json:"country,omitempty"`
	State              string  `json:"state,omitempty"`
	City               string  `json:"city,omitempty"`
	Zip                string  `json:"zip,omitempty"`
	LeadCreatedAt      string  `json:"lead_created_at,omitempty"`
	FreeTrialStartedAt string  `json:"free_trial_started_at,omitempty"`
	Errors             *Errors `json:"errors,omitempty"`
}

// Customers is result of listing: customers + paging.
type Customers struct {
	Customers   []Customer `json:"customers"`
	TotalPages  uint32     `json:"total_pages"`
	CurrentPage uint32     `json:"current_page"`
}

// ImportListCustomersParams = parameters for listing customers.
type ImportListCustomersParams struct {
	DataSourceUUID string `json:"data_source_uuid"`
	ExternalID     string `json:"external_id,omitempty"`
	Cursor
}

const (
	importCustomersEndpoint      = "import/customers"
	singleImportCustomerEndpoint = "import/customers/:uuid"
)

// ImportCreateCustomer loads the customer to Chartmogul.
//
// See https://dev.chartmogul.com/v1.0/reference#import-customer
func (api API) ImportCreateCustomer(customer *Customer, dataSourceUUID string) (*Customer, error) {
	customer.DataSourceUUID = dataSourceUUID
	result := &Customer{}
	return result, api.create(importCustomersEndpoint, customer, result)
}

// ImportListCustomers returns list of customers.
//
// https://dev.chartmogul.com/v1.0/reference#list-all-imported-customers
func (api API) ImportListCustomers(importListCustomersParams *ImportListCustomersParams) (*Customers, error) {
	result := &Customers{}
	return result, api.list(importCustomersEndpoint, result, *importListCustomersParams)
}

// ImportDeleteCustomer deletes the customer identified by their UUID.
//
// See https://dev.chartmogul.com/reference#delete-a-data-source
func (api API) ImportDeleteCustomer(uuid string) error {
	return api.delete(singleImportCustomerEndpoint, uuid)
}
