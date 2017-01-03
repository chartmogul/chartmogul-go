// Package chartmogul is a simple Go API library for Chartmogul public API.
//
// HTTP 2
//
// ChartMogul's current stable version of nginx is incompatible with HTTP 2 implementation of Go.
// For this reason the application must run with the following (or otherwise prohibit HTTP 2):
//  export GODEBUG=http2client=0
//
// Uses the library gorequest, which allows simple struct->query, body->struct,
// struct->body.
package chartmogul

import (
	"fmt"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	url           = "https://api.chartmogul.com/v1/%v"
	logFormatting = "%v: %v (page %v of %v)"

	// ErrKeyExternalID is key in Errors map indicating there's a problem with External ID of the resource.
	ErrKeyExternalID = "external_id"
	// ErrKeyTransactionExternalID is key in Errors map indicating there's a problem with External ID of the transaction.
	ErrKeyTransactionExternalID = "transactions.external_id"
	// ErrValCustomerExternalIDExists = can't import new customer with the same external ID
	ErrValCustomerExternalIDExists = "The external ID for this customer already exists in our system."
	// ErrValExternalIDExists = can't save Transaction, because it exists already.
	ErrValExternalIDExists = "has already been taken"
	// ErrValInvoiceExternalIDExists = invoice already exists
	ErrValInvoiceExternalIDExists = "The external ID for this invoice already exists in our system."
)

var timeout = 30 * time.Second

// IApi defines the interface of the library.
// Necessary eg. for mocks in testing.
type IApi interface {
	Ping() (res bool, err error)
	// Data sources
	ImportCreateDataSource(name string) (*DataSource, error)
	ImportRetrieveDataSource(dataSourceUUID string) (*DataSource, error)
	ImportListDataSources() (*DataSources, error)
	ImportDeleteDataSource(dataSourceUUID string) error
	// Customers
	ImportCreateCustomer(customer *Customer, dataSourceUUID string) (*Customer, error)
	ImportListCustomers(importListCustomersParams *ImportListCustomersParams) (*Customers, error)
	ImportDeleteCustomer(uuid string) error
	// Invoices
	ImportCreateInvoices(invoices []*Invoice, customerUUID string) (*Invoices, error)
	ImportListInvoices(cursor *Cursor, customerUUID string) (*Invoices, error)
	// Plans
	ImportCreatePlan(plan *Plan, dataSourceUUID string) (result *Plan, err error)
	ImportRetrievePlan(planUUID string) (*Plan, error)
	ImportListPlans(listPlansParams *ListPlansParams) (*Plans, error)
	ImportUpdatePlan(plan *Plan, planUUID string) (*Plan, error)
	ImportDeletePlan(planUUID string) error
	// Subscriptions
	ImportCancelSubscription(subscriptionUUID string, cancelledAt string) (*Subscription, error)
	ImportListSubscriptions(cursor *Cursor, customerUUID string) (*Subscriptions, error)
	// Transactions
	ImportCreateTransaction(transaction *Transaction, invoiceUUID string) (*Transaction, error)

	// Enrichment - Customers
	EnrichmentRetrieveCustomer(customerUUID string) (*EnrichmentCustomer, error)
	EnrichmentUpdateCustomer(enrichmentCustomer *EnrichmentCustomer, customerUUID string) (*EnrichmentCustomer, error)
	EnrichmentListCustomers(enrichmentListCustomersParams *EnrichmentListCustomersParams) (*EnrichmentCustomers, error)
	EnrichmentSearchCustomers(enrichmentSearchCustomersParams *EnrichmentSearchCustomersParams) (*EnrichmentCustomers, error)
	EnrichmentMergeCustomers(enrichmentMergeCustomersParams *EnrichmentMergeCustomersParams) error

	// Enrichment - Cusomer Attributes
	EnrichmentRetrieveCustomersAttributes(customerUUID string) (*AttributesResult, error)

	// Enrichment - Tags
	EnrichmentAddTagsToCustomer(customerUUID string, tags []string) (*TagsResult, error)
	EnrichmentAddTagsToCustomersWithEmail(email string, tags []string) (*EnrichmentCustomers, error)
	EnrichmentRemoveTagsFromCustomer(customerUUID string, tags []string) (*TagsResult, error)

	// Enrichment - Custom Attributes
	EnrichmentAddCustomAttributesToCustomer(customerUUID string, customAttributes []*CustomAttribute) (*CustomAttributes, error)
	EnrichmentAddCustomAttributesWithEmail(email string, customAttributes []*CustomAttribute) (*EnrichmentCustomers, error)
	EnrichmentUpdateCustomAttributesOfCustomer(customerUUID string, customAttributes map[string]interface{}) (*CustomAttributes, error)
	EnrichmentRemoveCustomAttributes(customerUUID string, customAttributes []string) (*CustomAttributes, error)

	// Metrics
	MetricsRetrieveAll(metricsFilter *MetricsFilter) (*MetricsResult, error)
	MetricsRetrieveMRR(metricsFilter *MetricsFilter) (*MRRResult, error)
	MetricsRetrieveARR(metricsFilter *MetricsFilter) (*ARRResult, error)
	MetricsRetrieveARPA(metricsFilter *MetricsFilter) (*ARPAResult, error)
	MetricsRetrieveASP(metricsFilter *MetricsFilter) (*ASPResult, error)
	MetricsRetrieveCustomerCount(metricsFilter *MetricsFilter) (*CustomerCountResult, error)
	MetricsRetrieveCustomerChurnRate(metricsFilter *MetricsFilter) (*CustomerChurnRateResult, error)
	MetricsRetrieveMRRChurnRate(metricsFilter *MetricsFilter) (*MRRChurnRateResult, error)
	MetricsRetrieveLTV(metricsFilter *MetricsFilter) (*LTVResult, error)

	// Metrics - Subscriptions & Activities
	MetricsListSubscriptions(cursor *Cursor, customerUUID string) (*MetricsSubscriptions, error)
	MetricsListActivities(cursor *Cursor, customerUUID string) (*MetricsActivities, error)
}

// API is the handle for communicating with Chartmogul.
type API struct {
	AccountToken string
	AccessKey    string
}

// Cursor contains query parameters for paging in CM.
// Attributes for query must be string, because gorequest library cannot convert anything else.
type Cursor struct {
	Page    string `json:"page,omitempty"`
	PerPage string `json:"per_page,omitempty"`
}

// Errors contains error feedback from ChartMogul
type Errors map[string]string

func (e Errors) Error() string {
	return fmt.Sprintf("chartmogul: %v", map[string]string(e))
}

// IsAlreadyExists is helper that returns true, if there's only one error
// and it means the uploaded resource of the same external_id already exists.
func (e Errors) IsAlreadyExists() (is bool) {
	if e == nil {
		return
	}
	if len(e) != 1 {
		return
	}
	msg, ok := e[ErrKeyExternalID]
	if !ok {
		msg, ok = e[ErrKeyTransactionExternalID]
	}
	return ok && (msg == ErrValExternalIDExists ||
		msg == ErrValCustomerExternalIDExists ||
		msg == ErrValInvoiceExternalIDExists)
}

// IsInvoiceAndTransactionAlreadyExist occurs when both invoice and tx exist already.
func (e Errors) IsInvoiceAndTransactionAlreadyExist() (is bool) {
	if e == nil {
		return
	}
	if len(e) == 2 {
		return
	}
	msg1, ok1 := e[ErrKeyExternalID]
	msg2, ok2 := e[ErrKeyTransactionExternalID]
	return ok1 && ok2 &&
		msg1 == ErrValInvoiceExternalIDExists && msg2 == ErrValExternalIDExists
}

// Setup configures static globals for the library
func Setup(timeoutConf time.Duration) {
	timeout = timeoutConf
}

func prepareURL(path string) string {
	return fmt.Sprintf(url, path)
}

func (api API) req(req *gorequest.SuperAgent) *gorequest.SuperAgent {
	// defaults for client go here:
	return req.Timeout(timeout).
		SetBasicAuth(api.AccountToken, api.AccessKey).
		Set("Content-Type", "application/json")
}
