package chartmogul

import (
	"fmt"
	"strings"
)

const invoicesEndpoint = "import/customers/:customerUUID/invoices"

// Invoices is wrapper for bulk importing invoices
type Invoices struct {
	CustomerUUID string     `json:"customer_uuid,omitempty"`
	CurrentPage  uint32     `json:"current_page,omitempty"`
	TotalPages   uint32     `json:"total_pages,omitempty"`
	Error        string     `json:"error,omitempty"`
	Invoices     []*Invoice `json:"invoices"`
}

func (i Invoices) String() string {
	if i.Error != "" {
		return fmt.Sprintf(logFormatting, i.CustomerUUID, i.Error, i.CurrentPage, i.TotalPages)
	}
	return fmt.Sprintf(logFormatting, i.CustomerUUID, i.Invoices, i.CurrentPage, i.TotalPages)
}

// Invoice is the data for ChartMogul to auto-generate subscriptions.
type Invoice struct {
	UUID         string         `json:"uuid,omitempty"`
	Currency     string         `json:"currency"`
	Date         string         `json:"date"`
	DueDate      string         `json:"due_date,omitempty"`
	ExternalID   string         `json:"external_id"`
	LineItems    []*LineItem    `json:"line_items"`
	Transactions []*Transaction `json:"transactions,omitempty"`
	Errors       *Errors        `json:"errors,omitempty"`
}

func (i Invoice) String() string {
	if i.Errors != nil && len(*i.Errors) != 0 {
		return fmt.Sprintf("Invoice{%v (%v)}: %v", i.UUID, i.ExternalID, *i.Errors)
	}
	return fmt.Sprintf("Invoice{%v (%v): items %v, tx %v}", i.UUID, i.ExternalID, len(i.LineItems), len(i.Transactions))
}

// LineItem represents a singular items of the invoices
type LineItem struct {
	AmountInCents          int    `json:"amount_in_cents,omitempty"`
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
	TaxAmountInCents       int    `json:"tax_amount_in_cents,omitempty"`
	Type                   string `json:"type"`
}

// CreateInvoices loads an invoice to a customer in Chartmogul.
// Customer must have a valid UUID! (use return value of API)
//
// See https://dev.chartmogul.com/reference#import-customers-invoices
func (api API) CreateInvoices(invoices []*Invoice, customerUUID string) (*Invoices, error) {
	if len(invoices) == 0 {
		return nil, nil
	}
	input := Invoices{Invoices: invoices}
	result := &Invoices{}

	path := strings.Replace(invoicesEndpoint, ":customerUUID", customerUUID, 1)
	return result, api.create(path, input, result)
}

// ListInvoices lists all imported invoices of customer with given UUID.
//
// See https://dev.chartmogul.com/reference#list-customers-invoices
func (api API) ListInvoices(cursor *Cursor, customerUUID string) (*Invoices, error) {
	result := &Invoices{}
	path := strings.Replace(invoicesEndpoint, ":customerUUID", customerUUID, 1)
	query := make([]interface{}, 0, 1)
	if cursor != nil {
		query = append(query, *cursor)
	}
	return result, api.list(path, result, query...)
}
