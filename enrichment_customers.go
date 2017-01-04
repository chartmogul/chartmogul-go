package chartmogul

// EnrichmentCustomer is the customer as represented in the Enrichment API.
type EnrichmentCustomer struct {
	ID uint32 `json:"id,omitempty"`
	// Basic info
	DataSourceUUID string `json:"data_source_uuid,omitempty"`
	UUID           string `json:"uuid,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
	Name           string `json:"name,omitempty"`
	Email          string `json:"email,omitempty"`
	Status         string `json:"status,omitempty"`
	CustomerSince  string `json:"customer-since,omitempty"`

	Attributes *Attributes `json:"attributes,omitempty"`
	Address    *Address    `json:"address,omitempty"`

	// Other info
	Mrr               string `json:"mrr,omitempty"`
	Arr               string `json:"arr,omitempty"`
	BillingSystemURL  string `json:"billing-system-url,omitempty"`
	ChartmogulURL     string `json:"chartmogul-url,omitempty"`
	BillingSystemType string `json:"billing-system-type,omitempty"`
	Currency          string `json:"currency,omitempty"`
	CurrencySign      string `json:"currency-sign,omitempty"`

	// For update
	Company            string `json:"company,omitempty"`
	Country            string `json:"country,omitempty"`
	State              string `json:"state,omitempty"`
	City               string `json:"city,omitempty"`
	LeadCreatedAt      string `json:"lead_created_at,omitempty"`
	FreeTrialStartedAt string `json:"free_trial_started_at,omitempty"`
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

// Attributes is subdocument of EnrichmentCustomer.
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

// Address is subdocument of EnrichmentCustomer.
type Address struct {
	AddressZIP string `json:"address_zip,omitempty"`
	City       string `json:"city,omitempty"`
	State      string `json:"state,omitempty"`
	Country    string `json:"country,omitempty"`
}

// Stripe is subdocument of EnrichmentCustomer.
type Stripe struct {
	UID    uint64 `json:"uid,omitempty"`
	Coupon bool   `json:"coupon,omitempty"`
}

// Clearbit is subdocument of EnrichmentCustomer.
type Clearbit struct {
	ID        string                 `json:"id,omitempty"`
	Name      string                 `json:"name,omitempty"`
	LegalName string                 `json:"legalName,omitempty"`
	Domain    string                 `json:"domain,omitempty"`
	URL       string                 `json:"url,omitempty"`
	Metrics   map[string]interface{} `json:"metrics,omitempty"`
	Category  *Category              `json:"category,omitempty"`
}

// Category is subdocument of EnrichmentCustomer.
type Category struct {
	Sector        string `json:"sector,omitempty"`
	IndustryGroup string `json:"industryGroup,omitempty"`
	Industry      string `json:"industry,omitempty"`
	SubIndustry   string `json:"subIndustry,omitempty"`
}

// EnrichmentListCustomersParams = parameters for listing customers in Enrichment API.
type EnrichmentListCustomersParams struct {
	DataSourceUUID string `json:"data_source_uuid,omitempty"`
	Status         string `json:"status,omitempty"`
	System         string `json:"system,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
	Cursor
}

// EnrichmentSearchCustomersParams - just email now.
type EnrichmentSearchCustomersParams struct {
	Email string `json:"email,omitempty"`
}

// EnrichmentCustomers is result of listing customers in Enrichment API.
type EnrichmentCustomers struct {
	Entries     []*EnrichmentCustomer `json:"entries,omitempty"`
	Page        uint32                `json:"page"`
	PerPage     uint32                `json:"per_page"`
	HasMore     bool                  `json:"has_more,omitempty"`
	CurrentPage int32                 `json:"current_page,omitempty"`
	TotalPages  int32                 `json:"total_pages,omitempty"`
}

// EnrichmentMergeCustomersParams - identify source and target for merging.
type EnrichmentMergeCustomersParams struct {
	From CustID `json:"from"`
	To   CustID `json:"to"`
}

// CustID - use either DataSourceUUID & ExternalID or CustomerUUID
type CustID struct {
	DataSourceUUID string `json:"data_source_uuid,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
	CustomerUUID   string `json:"customer_uuid,omitempty"`
}

const (
	singleEnrichmentCustomerEndpoint  = "customers/:uuid"
	enrichmentCustomersEndpoint       = "customers"
	enrichmentSearchCustomersEndpoint = "customers/search"
	enrichmentMergeCustomersEndpoint  = "customers/merges"
)

// EnrichmentCreateCustomer loads the customer to Chartmogul. New endpoint - with attributes.
//
// See https://dev.chartmogul.com/reference#customers-1
func (api API) EnrichmentCreateCustomer(newCustomer *NewCustomer) (*EnrichmentCustomer, error) {
	result := &EnrichmentCustomer{}
	return result, api.create(enrichmentCustomersEndpoint, newCustomer, result)
}

// EnrichmentRetrieveCustomer returns one customer as in Enrichment API.
//
// See https://dev.chartmogul.com/reference#retrieve-customer
func (api API) EnrichmentRetrieveCustomer(customerUUID string) (*EnrichmentCustomer, error) {
	result := &EnrichmentCustomer{}
	return result, api.retrieve(singleEnrichmentCustomerEndpoint, customerUUID, result)
}

// EnrichmentUpdateCustomer updates one customer in Enrichment API.
//
// See https://dev.chartmogul.com/reference#update-a-customer
func (api API) EnrichmentUpdateCustomer(enrichmentCustomer *EnrichmentCustomer, customerUUID string) (*EnrichmentCustomer, error) {
	result := &EnrichmentCustomer{}
	return result, api.update(singleEnrichmentCustomerEndpoint,
		customerUUID,
		enrichmentCustomer,
		result)
}

// EnrichmentListCustomers lists all EnrichmentCustomers for cutomer of given UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#list-a-customers-EnrichmentCustomers
func (api API) EnrichmentListCustomers(enrichmentListCustomersParams *EnrichmentListCustomersParams) (*EnrichmentCustomers, error) {
	result := &EnrichmentCustomers{}
	query := make([]interface{}, 0, 1)
	if enrichmentListCustomersParams != nil {
		query = append(query, *enrichmentListCustomersParams)
	}
	return result, api.list(enrichmentCustomersEndpoint, result, query...)
}

// EnrichmentSearchCustomers lists all EnrichmentCustomers for cutomer of given UUID.
//
// See https://dev.chartmogul.com/reference#search-for-customers
func (api API) EnrichmentSearchCustomers(enrichmentSearchCustomersParams *EnrichmentSearchCustomersParams) (*EnrichmentCustomers, error) {
	result := &EnrichmentCustomers{}
	return result, api.list(enrichmentSearchCustomersEndpoint, result, *enrichmentSearchCustomersParams)
}

// EnrichmentMergeCustomers merges two cutomers.
//
// See https://dev.chartmogul.com/reference#merge-customers
func (api API) EnrichmentMergeCustomers(enrichmentMergeCustomersParams *EnrichmentMergeCustomersParams) error {
	return api.list(enrichmentMergeCustomersEndpoint, *enrichmentMergeCustomersParams)
}
