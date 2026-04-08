package chartmogul

import "strings"

const (
	lineItemsEndpoint             = "import/invoices/:invoiceUUID/line_items"
	singleLineItemEndpoint        = "line_items/:uuid"
	lineItemDisabledStateEndpoint = "line_items/:uuid/disabled_state"
)

// LineItems is the result of creating line items.
type LineItems struct {
	LineItems []*LineItem `json:"line_items"`
}

// RetrieveLineItemParams optional parameters for RetrieveLineItem.
type RetrieveLineItemParams struct {
	ValidationType       string `json:"validation_type,omitempty"`
	IncludeEditHistories *bool  `json:"include_edit_histories,omitempty"`
	WithDisabled         *bool  `json:"with_disabled,omitempty"`
}

// UpdateLineItemParams holds the parameters for UpdateLineItem.
type UpdateLineItemParams struct {
	Type                      string `json:"type,omitempty"`
	AmountInCents             *int   `json:"amount_in_cents,omitempty"`
	Quantity                  *int   `json:"quantity,omitempty"`
	DiscountCode              string `json:"discount_code,omitempty"`
	DiscountAmountInCents     *int   `json:"discount_amount_in_cents,omitempty"`
	TaxAmountInCents          *int   `json:"tax_amount_in_cents,omitempty"`
	TransactionFeesInCents    *int   `json:"transaction_fees_in_cents,omitempty"`
	ExternalID                string `json:"external_id,omitempty"`
	AccountCode               string `json:"account_code,omitempty"`
	TransactionFeesCurrency   string `json:"transaction_fees_currency,omitempty"`
	DiscountDescription       string `json:"discount_description,omitempty"`
	EventOrder                *int   `json:"event_order,omitempty"`
	BalanceTransfer           *bool  `json:"balance_transfer,omitempty"`
	SubscriptionExternalID    string `json:"subscription_external_id,omitempty"`
	SubscriptionSetExternalID string `json:"subscription_set_external_id,omitempty"`
	PlanUUID                  string `json:"plan_uuid,omitempty"`
	ServicePeriodStart        string `json:"service_period_start,omitempty"`
	ServicePeriodEnd          string `json:"service_period_end,omitempty"`
	Prorated                  *bool  `json:"prorated,omitempty"`
	ProrationType             string `json:"proration_type,omitempty"`
	Description               string `json:"description,omitempty"`
	CancelledAt               string `json:"cancelled_at,omitempty"`
}

// ToggleLineItemDisabledParams holds the parameters for ToggleLineItemDisabled.
type ToggleLineItemDisabledParams struct {
	Disabled bool `json:"disabled"`
}

// CreateLineItems creates line items on an invoice.
func (api API) CreateLineItems(invoiceUUID string, lineItems []*LineItem) (*LineItems, error) {
	result := &LineItems{}
	path := strings.Replace(lineItemsEndpoint, ":invoiceUUID", invoiceUUID, 1)
	return result, api.create(path, LineItems{LineItems: lineItems}, result)
}

// RetrieveLineItem returns one line item by UUID.
func (api API) RetrieveLineItem(lineItemUUID string, params ...*RetrieveLineItemParams) (*LineItem, error) {
	result := &LineItem{}
	if len(params) > 0 && params[0] != nil {
		return result, api.retrieveWithParams(singleLineItemEndpoint, lineItemUUID, result, *params[0])
	}
	return result, api.retrieve(singleLineItemEndpoint, lineItemUUID, result)
}

// UpdateLineItem updates a line item by UUID.
func (api API) UpdateLineItem(lineItemUUID string, params *UpdateLineItemParams) (*LineItem, error) {
	result := &LineItem{}
	return result, api.update(singleLineItemEndpoint, lineItemUUID, params, result)
}

// ToggleLineItemDisabled toggles the disabled state of a line item.
func (api API) ToggleLineItemDisabled(lineItemUUID string, params *ToggleLineItemDisabledParams) (*LineItem, error) {
	result := &LineItem{}
	return result, api.update(lineItemDisabledStateEndpoint, lineItemUUID, params, result)
}

// DeleteLineItem deletes one line item by UUID.
func (api API) DeleteLineItem(lineItemUUID string) error {
	return api.delete(singleLineItemEndpoint, lineItemUUID)
}
