package chartmogul

import "strings"

// EditHistorySummary represents the edit history summary of an invoice
type EditHistorySummary struct {
	ValuesChanged         map[string]ValueChange `json:"values_changed,omitempty"`
	LatestEditAuthor      string                 `json:"latest_edit_author,omitempty"`
	LatestEditPerformedAt string                 `json:"latest_edit_performed_at,omitempty"`
}

// ValueChange represents a single value change in the edit history
type ValueChange struct {
	OriginalValue interface{} `json:"original_value,omitempty"`
	EditedValue   interface{} `json:"edited_value,omitempty"`
}

const (
	invoicesEndpoint             = "invoices"
	singleInvoiceEndpoint        = "invoices/:uuid"
	customersInvoicesEndpoint    = "import/customers/:customerUUID/invoices"
	invoiceUpdateStatusEndpoint  = "data_sources/:dataSourceUUID/invoices/:externalID/status"
	invoiceDisabledStateEndpoint = "invoices/:uuid/disabled_state"
)

// Invoices is wrapper for bulk importing invoices
// In case of /v1/invoices endpoint, the customer_uuid is on individual invoices and here it's empty.
type Invoices struct {
	CustomerUUID string     `json:"customer_uuid,omitempty"`
	Error        string     `json:"error,omitempty"`
	Invoices     []*Invoice `json:"invoices"`
	Pagination
}

// Invoice is the data for ChartMogul to auto-generate subscriptions.
type Invoice struct {
	UUID               string              `json:"uuid,omitempty"`
	CustomerUUID       string              `json:"customer_uuid,omitempty"`
	CustomerExternalID string              `json:"customer_external_id,omitempty"`
	Currency           string              `json:"currency"`
	DataSourceUUID     string              `json:"data_source_uuid,omitempty"`
	Date               string              `json:"date"`
	DueDate            string              `json:"due_date,omitempty"`
	ExternalID         string              `json:"external_id"`
	CollectionMethod   string              `json:"collection_method,omitempty"`
	Status             string              `json:"status,omitempty"`
	LineItems          []*LineItem         `json:"line_items"`
	Transactions       []*Transaction      `json:"transactions,omitempty"`
	Disabled           *bool               `json:"disabled,omitempty"`
	DisabledAt         string              `json:"disabled_at,omitempty"`
	DisabledBy         string              `json:"disabled_by,omitempty"`
	UserCreated        *bool               `json:"user_created,omitempty"`
	EditHistorySummary *EditHistorySummary `json:"edit_history_summary,omitempty"`
	Errors             *InvoiceErrors      `json:"errors,omitempty"`
}

// LineItem represents a singular items of the invoices
type LineItem struct {
	UUID                      string              `json:"uuid,omitempty"`
	AccountCode               string              `json:"account_code,omitempty"`
	AmountInCents             int                 `json:"amount_in_cents"`
	CancelledAt               string              `json:"cancelled_at,omitempty"`
	Description               string              `json:"description,omitempty"`
	DiscountAmountInCents     int                 `json:"discount_amount_in_cents,omitempty"`
	DiscountCode              string              `json:"discount_code,omitempty"`
	ExternalID                string              `json:"external_id,omitempty"`
	PlanUUID                  string              `json:"plan_uuid,omitempty"`
	PlanExternalID            string              `json:"plan_external_id,omitempty"`
	Prorated                  bool                `json:"prorated,omitempty"`
	ProrationType             string              `json:"proration_type,omitempty"`
	Quantity                  int                 `json:"quantity,omitempty"`
	ServicePeriodEnd          string              `json:"service_period_end,omitempty"`
	ServicePeriodStart        string              `json:"service_period_start,omitempty"`
	SubscriptionExternalID    string              `json:"subscription_external_id,omitempty"`
	SubscriptionSetExternalID string              `json:"subscription_set_external_id,omitempty"`
	SubscriptionUUID          string              `json:"subscription_uuid,omitempty"`
	TaxAmountInCents          int                 `json:"tax_amount_in_cents,omitempty"`
	TransactionFeesInCents    int                 `json:"transaction_fees_in_cents,omitempty"`
	TransactionFeesCurrency   string              `json:"transaction_fees_currency,omitempty"`
	DiscountDescription       string              `json:"discount_description,omitempty"`
	EventOrder                int                 `json:"event_order,omitempty"`
	BalanceTransfer           *bool               `json:"balance_transfer,omitempty"`
	Type                      string              `json:"type"`
	Disabled                  *bool               `json:"disabled,omitempty"`
	DisabledAt                string              `json:"disabled_at,omitempty"`
	DisabledBy                string              `json:"disabled_by,omitempty"`
	UserCreated               *bool               `json:"user_created,omitempty"`
	EditHistorySummary        *EditHistorySummary `json:"edit_history_summary,omitempty"`
	Errors                    interface{}         `json:"errors,omitempty"`
}

// ListAllInvoicesParams optional parameters for ListAllInvoices
type ListAllInvoicesParams struct {
	CustomerUUID         string `json:"customer_uuid,omitempty"`
	DataSourceUUID       string `json:"data_source_uuid,omitempty"`
	ExternalID           string `json:"external_id,omitempty"`
	ValidationType       string `json:"validation_type,omitempty"`        // Enum: "all", "valid" (default), "invalid"
	IncludeEditHistories *bool  `json:"include_edit_histories,omitempty"` // Include edit histories
	WithDisabled         *bool  `json:"with_disabled,omitempty"`          // Include disabled invoices
	Cursor
}

// RetrieveInvoiceParams optional parameters for RetrieveInvoice
type RetrieveInvoiceParams struct {
	ValidationType       string `json:"validation_type,omitempty"`        // Enum: "all", "valid" (default), "invalid"
	IncludeEditHistories *bool  `json:"include_edit_histories,omitempty"` // Include edit histories
	WithDisabled         *bool  `json:"with_disabled,omitempty"`          // Include disabled invoices
}

// CreateInvoices loads an invoice to a customer in Chartmogul.
// Customer must have a valid UUID! (use return value of API)
func (api API) CreateInvoices(invoices []*Invoice, customerUUID string) (*Invoices, error) {
	if len(invoices) == 0 {
		return nil, nil
	}
	input := Invoices{Invoices: invoices}
	result := &Invoices{}

	path := strings.Replace(customersInvoicesEndpoint, ":customerUUID", customerUUID, 1)
	return result, api.create(path, input, result)
}

// ListInvoices lists all imported invoices for a customer.
func (api API) ListInvoices(cursor *Cursor, customerUUID string) (*Invoices, error) {
	result := &Invoices{}
	path := strings.Replace(customersInvoicesEndpoint, ":customerUUID", customerUUID, 1)
	query := make([]interface{}, 0, 1)
	if cursor != nil {
		query = append(query, *cursor)
	}
	return result, api.list(path, result, query...)
}

// ListAllInvoices lists all imported invoices. Use parameters to narrow down the search/for paging.
// listAllInvoicesParams can be nil, in which case default values on API are used.
func (api API) ListAllInvoices(listAllInvoicesParams *ListAllInvoicesParams) (*Invoices, error) {
	result := &Invoices{}
	query := make([]interface{}, 0, 1)
	if listAllInvoicesParams != nil {
		query = append(query, *listAllInvoicesParams)
	}
	return result, api.list(invoicesEndpoint, result, query...)
}

// RetrieveInvoice returns one Invoice by UUID.
// Optionally accepts RetrieveInvoiceParams for additional query options.
func (api API) RetrieveInvoice(invoiceUUID string, params ...*RetrieveInvoiceParams) (*Invoice, error) {
	result := &Invoice{}

	if len(params) > 0 && params[0] != nil {
		// Params provided - use retrieveWithParams
		return result, api.retrieveWithParams(singleInvoiceEndpoint, invoiceUUID, result, *params[0])
	}

	// No params provided - use original retrieve method
	return result, api.retrieve(singleInvoiceEndpoint, invoiceUUID, result)
}

// DeleteInvoice deletes one invoice by UUID.
func (api API) DeleteInvoice(invoiceUUID string) error {
	return api.delete(singleInvoiceEndpoint, invoiceUUID)
}

// UpdateInvoiceStatusParams holds the parameters for UpdateInvoiceStatus.
type UpdateInvoiceStatusParams struct {
	Status string `json:"status"`
}

// UpdateInvoiceStatus updates the status of an invoice by data source UUID and invoice external ID.
func (api API) UpdateInvoiceStatus(dataSourceUUID, invoiceExternalID string, params *UpdateInvoiceStatusParams) error {
	path := strings.Replace(invoiceUpdateStatusEndpoint, ":dataSourceUUID", dataSourceUUID, 1)
	path = strings.Replace(path, ":externalID", invoiceExternalID, 1)
	result := &Invoice{}
	return api.update(path, "", params, result)
}

// ToggleInvoiceDisabledParams holds the parameters for ToggleInvoiceDisabled.
type ToggleInvoiceDisabledParams struct {
	Disabled bool `json:"disabled"`
}

// ToggleInvoiceDisabled toggles the disabled state of an invoice.
func (api API) ToggleInvoiceDisabled(invoiceUUID string, params *ToggleInvoiceDisabledParams) (*Invoice, error) {
	result := &Invoice{}
	return result, api.update(invoiceDisabledStateEndpoint, invoiceUUID, params, result)
}

// UpdateInvoiceParams holds the parameters for UpdateInvoice.
type UpdateInvoiceParams struct {
	Date             string `json:"date,omitempty"`
	DueDate          string `json:"due_date,omitempty"`
	Currency         string `json:"currency,omitempty"`
	CollectionMethod string `json:"collection_method,omitempty"`
}

// UpdateInvoice updates an invoice by UUID.
func (api API) UpdateInvoice(invoiceUUID string, params *UpdateInvoiceParams) (*Invoice, error) {
	result := &Invoice{}
	return result, api.update(singleInvoiceEndpoint, invoiceUUID, params, result)
}
