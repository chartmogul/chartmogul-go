package chartmogul

import "strings"

const (
	invoicesEndpoint          = "invoices"
	singleInvoiceEndpoint     = "invoices/:uuid"
	customersInvoicesEndpoint = "import/customers/:customerUUID/invoices"
)

// Invoices is wrapper for bulk importing invoices
// In case of /v1/invoices endpoint, the customer_uuid is on individual invoices and here it's empty.
type Invoices struct {
	CustomerUUID string     `json:"customer_uuid,omitempty"`
	CurrentPage  uint32     `json:"current_page,omitempty"`
	TotalPages   uint32     `json:"total_pages,omitempty"`
	Error        string     `json:"error,omitempty"`
	Invoices     []*Invoice `json:"invoices"`
}

// Invoice is the data for ChartMogul to auto-generate subscriptions.
type Invoice struct {
	UUID           string         `json:"uuid,omitempty"`
	CustomerUUID   string         `json:"customer_uuid,omitempty"`
	Currency       string         `json:"currency"`
	DataSourceUUID string         `json:"data_source_uuid,omitempty"`
	Date           string         `json:"date"`
	DueDate        string         `json:"due_date,omitempty"`
	ExternalID     string         `json:"external_id"`
	LineItems      []*LineItem    `json:"line_items"`
	Transactions   []*Transaction `json:"transactions,omitempty"`
	Errors         *Errors        `json:"errors,omitempty"`
}

// LineItem represents a singular items of the invoices
type LineItem struct {
	UUID                   string `json:"uuid,omitempty"`
	AccountCode            string `json:"account_code,omitempty"`
	AmountInCents          int    `json:"amount_in_cents"`
	CancelledAt            string `json:"cancelled_at,omitempty"`
	Description            string `json:"description,omitempty"`
	DiscountAmountInCents  int    `json:"discount_amount_in_cents,omitempty"`
	DiscountCode           string `json:"discount_code,omitempty"`
	ExternalID             string `json:"external_id,omitempty"`
	PlanUUID               string `json:"plan_uuid,omitempty"`
	Prorated               bool   `json:"prorated,omitempty"`
	Quantity               int    `json:"quantity,omitempty"`
	ServicePeriodEnd       string `json:"service_period_end,omitempty"`
	ServicePeriodStart     string `json:"service_period_start,omitempty"`
	SubscriptionExternalID string `json:"subscription_external_id,omitempty"`
	SubscriptionUUID       string `json:"subscription_uuid,omitempty"`
	TaxAmountInCents       int    `json:"tax_amount_in_cents,omitempty"`
	TransactionFeesInCents int    `json:"transaction_fees_in_cents,omitempty"`
	Type                   string `json:"type"`
}

// ListAllInvoicesParams optional parameters for ListAllInvoices
type ListAllInvoicesParams struct {
	CustomerUUID   string `json:"customer_uuid,omitempty"`
	DataSourceUUID string `json:"data_source_uuid,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
	Cursor
}

// CreateInvoices loads an invoice to a customer in Chartmogul.
// Customer must have a valid UUID! (use return value of API)
//
// See https://dev.chartmogul.com/v1.0/reference#invoices
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
//
// See https://dev.chartmogul.com/v1.0/reference#invoices
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
//
// See https://dev.chartmogul.com/v1.0/reference#invoices
func (api API) ListAllInvoices(listAllInvoicesParams *ListAllInvoicesParams) (*Invoices, error) {
	result := &Invoices{}
	query := make([]interface{}, 0, 1)
	if listAllInvoicesParams != nil {
		query = append(query, *listAllInvoicesParams)
	}
	return result, api.list(invoicesEndpoint, result, query...)
}

// RetrieveInvoice returns one Invoice by UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#invoices
func (api API) RetrieveInvoice(invoiceUUID string) (*Invoice, error) {
	result := &Invoice{}
	return result, api.retrieve(singleInvoiceEndpoint, invoiceUUID, result)
}

// DeleteInvoice deletes one invoice by UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#invoices
func (api API) DeleteInvoice(invoiceUUID string) error {
	return api.delete(singleInvoiceEndpoint, invoiceUUID)
}
