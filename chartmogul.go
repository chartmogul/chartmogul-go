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
	logFormatting = "%v: %v (page %v of %v)"

	// ErrKeyExternalID is key in Errors map indicating there's a problem with External ID of the resource.
	ErrKeyExternalID = "external_id"
	// ErrKeyTransactionExternalID is key in Errors map indicating there's a problem with External ID of the transaction.
	ErrKeyTransactionExternalID = "transactions.external_id"
	// ErrKeyName - data source name
	ErrKeyName = "name"
	// ErrValCustomerExternalIDExists = can't import new customer with the same external ID
	ErrValCustomerExternalIDExists = "The external ID for this customer already exists in our system."
	// ErrValExternalIDExists = can't save Transaction, because it exists already.
	ErrValExternalIDExists = "has already been taken"
	// ErrValInvoiceExternalIDExists = invoice already exists
	ErrValInvoiceExternalIDExists = "The external ID for this invoice already exists in our system."
	// ErrValPlanExternalIDExists = plan already exists
	ErrValPlanExternalIDExists = "A plan with this identifier already exists in our system."
	// ErrValHasAlreadyBeenTaken = data source name taken
	ErrValHasAlreadyBeenTaken = "Has already been taken."
)

var (
	url     = "https://api.chartmogul.com/v1/%v"
	timeout = 30 * time.Second
)

// IApi defines the interface of the library.
// Necessary eg. for mocks in testing.
type IApi interface {
	Ping() (res bool, err error)
	// Data sources
	CreateDataSource(name string) (*DataSource, error)
	RetrieveDataSource(dataSourceUUID string) (*DataSource, error)
	ListDataSources() (*DataSources, error)
	DeleteDataSource(dataSourceUUID string) error
	// Invoices
	CreateInvoices(invoices []*Invoice, customerUUID string) (*Invoices, error)
	ListInvoices(cursor *Cursor, customerUUID string) (*Invoices, error)
	ListAllInvoices(listAllInvoicesParams *ListAllInvoicesParams) (*Invoices, error)
	DeleteInvoice(invoiceUUID string) error
	// Plans
	CreatePlan(plan *Plan) (result *Plan, err error)
	RetrievePlan(planUUID string) (*Plan, error)
	ListPlans(listPlansParams *ListPlansParams) (*Plans, error)
	UpdatePlan(plan *Plan, planUUID string) (*Plan, error)
	DeletePlan(planUUID string) error
	// Subscriptions
	CancelSubscription(subscriptionUUID string, cancelSubscriptionParams *CancelSubscriptionParams) (*Subscription, error)
	ListSubscriptions(cursor *Cursor, customerUUID string) (*Subscriptions, error)
	// Transactions
	CreateTransaction(transaction *Transaction, invoiceUUID string) (*Transaction, error)

	// Customers
	CreateCustomer(newCustomer *NewCustomer) (*Customer, error)
	RetrieveCustomer(customerUUID string) (*Customer, error)
	UpdateCustomer(Customer *Customer, customerUUID string) (*Customer, error)
	ListCustomers(ListCustomersParams *ListCustomersParams) (*Customers, error)
	SearchCustomers(SearchCustomersParams *SearchCustomersParams) (*Customers, error)
	MergeCustomers(MergeCustomersParams *MergeCustomersParams) error
	DeleteCustomer(customerUUID string) error

	//  - Cusomer Attributes
	RetrieveCustomersAttributes(customerUUID string) (*Attributes, error)

	//  Tags
	AddTagsToCustomer(customerUUID string, tags []string) (*TagsResult, error)
	AddTagsToCustomersWithEmail(email string, tags []string) (*Customers, error)
	RemoveTagsFromCustomer(customerUUID string, tags []string) (*TagsResult, error)

	// Custom Attributes
	AddCustomAttributesToCustomer(customerUUID string, customAttributes []*CustomAttribute) (*CustomAttributes, error)
	AddCustomAttributesWithEmail(email string, customAttributes []*CustomAttribute) (*Customers, error)
	UpdateCustomAttributesOfCustomer(customerUUID string, customAttributes map[string]interface{}) (*CustomAttributes, error)
	RemoveCustomAttributes(customerUUID string, customAttributes []string) (*CustomAttributes, error)

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
	Page    uint32 `json:"page,omitempty"`
	PerPage uint32 `json:"per_page,omitempty"`
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
	if !ok {
		msg, ok = e[ErrKeyName]
		return ok && msg == ErrValHasAlreadyBeenTaken
	}
	return msg == ErrValExternalIDExists ||
		msg == ErrValCustomerExternalIDExists ||
		msg == ErrValPlanExternalIDExists ||
		msg == ErrValInvoiceExternalIDExists
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

// Setup configures global timeout for the library.
func Setup(timeoutConf time.Duration) {
	timeout = timeoutConf
}

// SetURL changes target URL for the module globally.
func SetURL(specialURL string) {
	url = specialURL
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
