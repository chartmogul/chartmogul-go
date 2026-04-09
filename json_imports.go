package chartmogul

import "strings"

const (
	jsonImportsEndpoint      = "data_sources/:dataSourceUUID/json_imports"
	singleJsonImportEndpoint = "data_sources/:dataSourceUUID/json_imports/:id"
)

// JsonImport represents the response from a bulk import operation.
type JsonImport struct {
	ID             string      `json:"id,omitempty"`
	DataSourceUUID string      `json:"data_source_uuid,omitempty"`
	Status         string      `json:"status,omitempty"`
	ExternalID     string      `json:"external_id,omitempty"`
	StatusDetails  interface{} `json:"status_details,omitempty"`
	CreatedAt      string      `json:"created_at,omitempty"`
	UpdatedAt      string      `json:"updated_at,omitempty"`
}

// JsonImportData represents the request body for creating a bulk import.
type JsonImportData struct {
	ExternalID         string                         `json:"external_id"`
	Customers          []*JsonImportCustomer          `json:"customers,omitempty"`
	Plans              []*JsonImportPlan              `json:"plans,omitempty"`
	Invoices           []*JsonImportInvoice           `json:"invoices,omitempty"`
	LineItems          []*JsonImportLineItem          `json:"line_items,omitempty"`
	Transactions       []*JsonImportTransaction       `json:"transactions,omitempty"`
	SubscriptionEvents []*JsonImportSubscriptionEvent `json:"subscription_events,omitempty"`
}

// JsonImportCustomer represents a customer in a bulk import.
type JsonImportCustomer struct {
	ExternalID         string `json:"external_id"`
	Name               string `json:"name,omitempty"`
	Email              string `json:"email,omitempty"`
	Country            string `json:"country,omitempty"`
	State              string `json:"state,omitempty"`
	City               string `json:"city,omitempty"`
	Zip                string `json:"zip,omitempty"`
	Company            string `json:"company,omitempty"`
	LeadCreatedAt      string `json:"lead_created_at,omitempty"`
	FreeTrialStartedAt string `json:"free_trial_started_at,omitempty"`
}

// JsonImportPlan represents a plan in a bulk import.
type JsonImportPlan struct {
	Name          string `json:"name"`
	IntervalCount int    `json:"interval_count"`
	IntervalUnit  string `json:"interval_unit"`
	ExternalID    string `json:"external_id,omitempty"`
}

// JsonImportInvoice represents an invoice in a bulk import.
type JsonImportInvoice struct {
	CustomerExternalID string `json:"customer_external_id"`
	ExternalID         string `json:"external_id"`
	Date               string `json:"date"`
	DueDate            string `json:"due_date,omitempty"`
	Currency           string `json:"currency,omitempty"`
	CollectionMethod   string `json:"collection_method,omitempty"`
}

// JsonImportLineItem represents a line item in a bulk import.
type JsonImportLineItem struct {
	InvoiceExternalID         string `json:"invoice_external_id"`
	Type                      string `json:"type"`
	AmountInCents             int    `json:"amount_in_cents"`
	ExternalID                string `json:"external_id,omitempty"`
	SubscriptionExternalID    string `json:"subscription_external_id,omitempty"`
	SubscriptionSetExternalID string `json:"subscription_set_external_id,omitempty"`
	PlanExternalID            string `json:"plan_external_id,omitempty"`
	ServicePeriodStart        string `json:"service_period_start,omitempty"`
	ServicePeriodEnd          string `json:"service_period_end,omitempty"`
	Prorated                  *bool  `json:"prorated,omitempty"`
	ProrationType             string `json:"proration_type,omitempty"`
	Description               string `json:"description,omitempty"`
	Quantity                  *int   `json:"quantity,omitempty"`
	DiscountAmountInCents     *int   `json:"discount_amount_in_cents,omitempty"`
	DiscountCode              string `json:"discount_code,omitempty"`
	TaxAmountInCents          *int   `json:"tax_amount_in_cents,omitempty"`
	TransactionFeesInCents    *int   `json:"transaction_fees_in_cents,omitempty"`
	TransactionFeesCurrency   string `json:"transaction_fees_currency,omitempty"`
	DiscountDescription       string `json:"discount_description,omitempty"`
	EventOrder                *int   `json:"event_order,omitempty"`
}

// JsonImportTransaction represents a transaction in a bulk import.
type JsonImportTransaction struct {
	InvoiceExternalID string `json:"invoice_external_id"`
	Type              string `json:"type"`
	Result            string `json:"result"`
	Date              string `json:"date"`
	ExternalID        string `json:"external_id,omitempty"`
	AmountInCents     *int   `json:"amount_in_cents,omitempty"`
}

// JsonImportSubscriptionEvent represents a subscription event in a bulk import.
type JsonImportSubscriptionEvent struct {
	CustomerExternalID        string `json:"customer_external_id"`
	EventType                 string `json:"event_type"`
	EventDate                 string `json:"event_date"`
	EffectiveDate             string `json:"effective_date"`
	SubscriptionExternalID    string `json:"subscription_external_id,omitempty"`
	SubscriptionSetExternalID string `json:"subscription_set_external_id,omitempty"`
	ExternalID                string `json:"external_id,omitempty"`
	PlanExternalID            string `json:"plan_external_id,omitempty"`
	Currency                  string `json:"currency,omitempty"`
	AmountInCents             *int   `json:"amount_in_cents,omitempty"`
	Quantity                  *int   `json:"quantity,omitempty"`
	TaxAmountInCents          *int   `json:"tax_amount_in_cents,omitempty"`
	RetractedEventId          string `json:"retracted_event_id,omitempty"`
}

// CreateJsonImport creates a bulk import for a data source.
func (api API) CreateJsonImport(dataSourceUUID string, data *JsonImportData) (*JsonImport, error) {
	result := &JsonImport{}
	path := strings.Replace(jsonImportsEndpoint, ":dataSourceUUID", dataSourceUUID, 1)
	return result, api.create(path, data, result)
}

// RetrieveJsonImport retrieves the status of a bulk import.
func (api API) RetrieveJsonImport(dataSourceUUID string, importID string) (*JsonImport, error) {
	result := &JsonImport{}
	path := strings.Replace(singleJsonImportEndpoint, ":dataSourceUUID", dataSourceUUID, 1)
	path = strings.Replace(path, ":id", importID, 1)
	return result, api.retrieve(path, "", result)
}
