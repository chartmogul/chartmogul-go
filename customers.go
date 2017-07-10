package chartmogul

import "strings"

// Customer is the customer as represented in the API.
type Customer struct {
	ID uint32 `json:"id,omitempty"`
	// Basic info
	DataSourceUUID  string   `json:"data_source_uuid,omitempty"`
	DataSourceUUIDs []string `json:"data_source_uuids,omitempty"`
	UUID            string   `json:"uuid,omitempty"`
	ExternalID      string   `json:"external_id,omitempty"`
	ExternalIDs     []string `json:"external_ids,omitempty"`
	Name            string   `json:"name,omitempty"`
	Email           string   `json:"email,omitempty"`
	Status          string   `json:"status,omitempty"`
	CustomerSince   string   `json:"customer-since,omitempty"`

	Attributes *Attributes `json:"attributes,omitempty"`
	Address    *Address    `json:"address,omitempty"`

	// Other info
	Mrr               float64 `json:"mrr,omitempty"`
	Arr               float64 `json:"arr,omitempty"`
	BillingSystemURL  string  `json:"billing-system-url,omitempty"`
	ChartmogulURL     string  `json:"chartmogul-url,omitempty"`
	BillingSystemType string  `json:"billing-system-type,omitempty"`
	Currency          string  `json:"currency,omitempty"`
	CurrencySign      string  `json:"currency-sign,omitempty"`

	// For update
	Company            string `json:"company,omitempty"`
	Country            string `json:"country,omitempty"`
	State              string `json:"state,omitempty"`
	City               string `json:"city,omitempty"`
	Zip                string `json:"zip,omitempty"`
	LeadCreatedAt      string `json:"lead_created_at,omitempty"`
	FreeTrialStartedAt string `json:"free_trial_started_at,omitempty"`

	Errors Errors `json:"errors,omitempty"`
}

// NewCustomer allows creating customer on a new endpoint.
type NewCustomer struct {
	// Obligatory
	DataSourceUUID string `json:"data_source_uuid"`
	ExternalID     string `json:"external_id,omitempty"`
	Name           string `json:"name,omitempty"`

	//Optional
	Email      string         `json:"email,omitempty"`
	Attributes *NewAttributes `json:"attributes,omitempty"`
	// Address
	Company string `json:"company,omitempty"`
	Country string `json:"country,omitempty"`
	State   string `json:"state,omitempty"`
	City    string `json:"city,omitempty"`
	Zip     string `json:"zip,omitempty"`
	// Lead/Trial
	LeadCreatedAt      string `json:"lead_created_at,omitempty"`
	FreeTrialStartedAt string `json:"free_trial_started_at,omitempty"`
}

// Attributes is subdocument of Customer.
type Attributes struct {
	Tags     []string               `json:"tags,omitempty"`
	Stripe   *Stripe                `json:"stripe,omitempty"`
	Clearbit *Clearbit              `json:"clearbit,omitempty"`
	Custom   map[string]interface{} `json:"custom,omitempty"`
}

// NewAttributes is subdocument of NewCustomer.
type NewAttributes struct {
	Tags   []string           `json:"tags,omitempty"`
	Custom []*CustomAttribute `json:"custom,omitempty"`
}

// Address is subdocument of Customer.
type Address struct {
	AddressZIP string `json:"address_zip,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
}

// Stripe is subdocument of Customer.
type Stripe struct {
	UID    uint64 `json:"uid,omitempty"`
	Coupon bool   `json:"coupon,omitempty"`
}

// Clearbit is subdocument of Customer.
type Clearbit struct {
	ID        string                 `json:"id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	LegalName string                 `json:"legalName,omitempty"`
	Domain    string                 `json:"domain,omitempty"`
	URL       string                 `json:"url,omitempty"`
	Metrics   map[string]interface{} `json:"metrics,omitempty"`
	Category  *Category              `json:"category,omitempty"`
}

// Category is subdocument of Customer.
type Category struct {
	Sector        string `json:"sector,omitempty"`
	IndustryGroup string `json:"industryGroup,omitempty"`
	Industry      string `json:"industry,omitempty"`
	SubIndustry   string `json:"subIndustry,omitempty"`
}

// ListCustomersParams = parameters for listing customers in API.
type ListCustomersParams struct {
	DataSourceUUID string `json:"data_source_uuid,omitempty"`
	Status         string `json:"status,omitempty"`
	System         string `json:"system,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
	Cursor
}

// SearchCustomersParams - just email now.
type SearchCustomersParams struct {
	Email string `json:"email,omitempty"`
	Cursor
}

// Customers is result of listing customers in API.
type Customers struct {
	Entries     []*Customer `json:"entries,omitempty"`
	Page        uint32      `json:"page"`
	PerPage     uint32      `json:"per_page"`
	HasMore     bool        `json:"has_more,omitempty"`
	CurrentPage int32       `json:"current_page,omitempty"`
	TotalPages  int32       `json:"total_pages,omitempty"`
}

// MergeCustomersParams - identify source and target for merging.
type MergeCustomersParams struct {
	From CustID `json:"from"`
	Into CustID `json:"into"`
}

// CustID - use either DataSourceUUID & ExternalID or CustomerUUID
type CustID struct {
	DataSourceUUID string `json:"data_source_uuid,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
	CustomerUUID   string `json:"customer_uuid,omitempty"`
}

const (
	singleCustomerEndpoint  = "customers/:uuid"
	purgeCustomerEndpoint   = "data_sources/:data_source_uuid/customers/:uuid/invoices"
	customersEndpoint       = "customers"
	searchCustomersEndpoint = "customers/search"
	mergeCustomersEndpoint  = "customers/merges"
)

// CreateCustomer loads the customer to Chartmogul. New endpoint - with attributes.
//
// See https://dev.chartmogul.com/v1.0/reference#customers
func (api API) CreateCustomer(newCustomer *NewCustomer) (*Customer, error) {
	result := &Customer{}
	return result, api.create(customersEndpoint, newCustomer, result)
}

// RetrieveCustomer returns one customer as in API.
//
// See https://dev.chartmogul.com/v1.0/reference#customers
func (api API) RetrieveCustomer(customerUUID string) (*Customer, error) {
	result := &Customer{}
	return result, api.retrieve(singleCustomerEndpoint, customerUUID, result)
}

// UpdateCustomer updates one customer in API.
//
// See https://dev.chartmogul.com/v1.0/reference#customers
func (api API) UpdateCustomer(customer *Customer, customerUUID string) (*Customer, error) {
	result := &Customer{}
	return result, api.update(singleCustomerEndpoint,
		customerUUID,
		customer,
		result)
}

// ListCustomers lists all Customers for cutomer of given UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#customers
func (api API) ListCustomers(listCustomersParams *ListCustomersParams) (*Customers, error) {
	result := &Customers{}
	query := make([]interface{}, 0, 1)
	if listCustomersParams != nil {
		query = append(query, *listCustomersParams)
	}
	return result, api.list(customersEndpoint, result, query...)
}

// SearchCustomers lists all Customers for cutomer of given UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#customers
func (api API) SearchCustomers(searchCustomersParams *SearchCustomersParams) (*Customers, error) {
	result := &Customers{}
	return result, api.list(searchCustomersEndpoint, result, *searchCustomersParams)
}

// MergeCustomers merges two cutomers.
//
// See https://dev.chartmogul.com/v1.0/reference#customers
func (api API) MergeCustomers(mergeCustomersParams *MergeCustomersParams) error {
	return api.merge(mergeCustomersEndpoint, *mergeCustomersParams)
}

// DeleteCustomer deletes one customer by UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#customers
func (api API) DeleteCustomer(customerUUID string) error {
	return api.delete(singleCustomerEndpoint, customerUUID)
}

// DeleteCustomerInvoices deletes all customer's invoices by UUID for given data source UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#customers
func (api API) DeleteCustomerInvoices(dataSourceUUID, customerUUID string) error {
	path := strings.Replace(purgeCustomerEndpoint, ":data_source_uuid", dataSourceUUID, 1)
	return api.delete(path, customerUUID)
}
