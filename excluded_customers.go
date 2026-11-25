package chartmogul

// ExcludedCustomer is one excluded data source customer, as represented in the API.
// The same CustomerUUID appears once per data source it is excluded from.
type ExcludedCustomer struct {
	CustomerUUID   string `json:"customer_uuid"`
	DataSourceUUID string `json:"data_source_uuid"`
	ExternalID     string `json:"external_id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	ExcludedAt     string `json:"excluded_at"`
}

// ListExcludedCustomersParams = parameters for listing excluded customers in API.
type ListExcludedCustomersParams struct {
	Cursor
}

// ExcludedCustomers is result of listing excluded customers in API, newest first.
type ExcludedCustomers struct {
	Entries []*ExcludedCustomer `json:"entries"`
	Pagination
}

// ReincludedCustomer is the result of removing a customer from the exclude list.
type ReincludedCustomer struct {
	CustomerUUID string `json:"customer_uuid"`
	Status       string `json:"status"` // ReincludedCustomerStatusReimporting | ReincludedCustomerStatusManualReimportRequired
}

const (
	// ReincludedCustomerStatusReimporting means ChartMogul started re-importing the customer
	// from at least one of its data sources.
	ReincludedCustomerStatusReimporting = "reimporting"
	// ReincludedCustomerStatusManualReimportRequired means none of the customer's data sources
	// can re-import it automatically, so its data has to be sent again (e.g. via the Import API).
	ReincludedCustomerStatusManualReimportRequired = "manual_reimport_required"
)

const (
	excludedCustomersEndpoint      = "excluded_customers"
	singleExcludedCustomerEndpoint = "excluded_customers/:uuid"
)

// ListExcludedCustomers lists the account's excluded customers.
func (api API) ListExcludedCustomers(listExcludedCustomersParams *ListExcludedCustomersParams) (*ExcludedCustomers, error) {
	result := &ExcludedCustomers{}
	query := make([]interface{}, 0, 1)
	if listExcludedCustomersParams != nil {
		query = append(query, *listExcludedCustomersParams)
	}
	return result, api.list(excludedCustomersEndpoint, result, query...)
}

// DeleteExcludedCustomer removes the customer from the exclude list, re-including it.
func (api API) DeleteExcludedCustomer(customerUUID string) (*ReincludedCustomer, error) {
	result := &ReincludedCustomer{}
	return result, api.deleteWithOutput(singleExcludedCustomerEndpoint, customerUUID, result)
}
