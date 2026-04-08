package chartmogul

import "strings"

const (
	transactionsEndpoint             = "import/invoices/:invoiceUUID/transactions"
	singleTransactionEndpoint        = "transactions/:uuid"
	transactionDisabledStateEndpoint = "transactions/:uuid/disabled_state"
)

// Transaction is either payment/refund on an invoice, for its full value.
type Transaction struct {
	UUID                   string              `json:"uuid,omitempty"`
	Date                   string              `json:"date"`
	ExternalID             string              `json:"external_id,omitempty"`
	Result                 string              `json:"result"`
	Type                   string              `json:"type"`
	AmountInCents          *int                `json:"amount_in_cents,omitempty"`
	TransactionFeesInCents *int                `json:"transaction_fees_in_cents,omitempty"`
	TransactionFeesCurrency string             `json:"transaction_fees_currency,omitempty"`
	Disabled               *bool               `json:"disabled,omitempty"`
	DisabledAt             string              `json:"disabled_at,omitempty"`
	DisabledBy             string              `json:"disabled_by,omitempty"`
	UserCreated            *bool               `json:"user_created,omitempty"`
	EditHistorySummary     *EditHistorySummary `json:"edit_history_summary,omitempty"`
	Errors                 Errors              `json:"errors,omitempty"`
}

// RetrieveTransactionParams optional parameters for RetrieveTransaction.
type RetrieveTransactionParams struct {
	ValidationType       string `json:"validation_type,omitempty"`
	IncludeEditHistories *bool  `json:"include_edit_histories,omitempty"`
	WithDisabled         *bool  `json:"with_disabled,omitempty"`
}

// UpdateTransactionParams holds the parameters for UpdateTransaction.
type UpdateTransactionParams struct {
	Type                    string `json:"type,omitempty"`
	Date                    string `json:"date,omitempty"`
	Result                  string `json:"result,omitempty"`
	AmountInCents           *int   `json:"amount_in_cents,omitempty"`
	TransactionFeesInCents  *int   `json:"transaction_fees_in_cents,omitempty"`
	TransactionFeesCurrency string `json:"transaction_fees_currency,omitempty"`
}

// ToggleTransactionDisabledParams holds the parameters for ToggleTransactionDisabled.
type ToggleTransactionDisabledParams struct {
	Disabled bool `json:"disabled"`
}

// CreateTransaction loads a transaction to an invoice in Chartmogul.
func (api API) CreateTransaction(transaction *Transaction, invoiceUUID string) (*Transaction, error) {
	result := &Transaction{}
	path := strings.Replace(transactionsEndpoint, ":invoiceUUID", invoiceUUID, 1)
	return result, api.create(path, transaction, result)
}

// RetrieveTransaction returns one transaction by UUID.
func (api API) RetrieveTransaction(transactionUUID string, params ...*RetrieveTransactionParams) (*Transaction, error) {
	result := &Transaction{}
	if len(params) > 0 && params[0] != nil {
		return result, api.retrieveWithParams(singleTransactionEndpoint, transactionUUID, result, *params[0])
	}
	return result, api.retrieve(singleTransactionEndpoint, transactionUUID, result)
}

// UpdateTransaction updates a transaction by UUID.
func (api API) UpdateTransaction(transactionUUID string, params *UpdateTransactionParams) (*Transaction, error) {
	result := &Transaction{}
	return result, api.update(singleTransactionEndpoint, transactionUUID, params, result)
}

// ToggleTransactionDisabled toggles the disabled state of a transaction.
func (api API) ToggleTransactionDisabled(transactionUUID string, params *ToggleTransactionDisabledParams) (*Transaction, error) {
	result := &Transaction{}
	return result, api.update(transactionDisabledStateEndpoint, transactionUUID, params, result)
}

// DeleteTransaction deletes one transaction by UUID.
func (api API) DeleteTransaction(transactionUUID string) error {
	return api.delete(singleTransactionEndpoint, transactionUUID)
}
