// Package chartmogul is a simple Go API library for Chartmogul public API.
//
// HTTP 2
// ChartMogul's current stable version of nginx is incompatible with HTTP 2 implementation of Go.
// For this reason the application must run with the following (or otherwise prohibit HTTP 2):
// export GODEBUG=http2client=0
//
// Uses the library gorequest, which allows simple struct->query, body->struct,
// struct->body.
package chartmogul

import (
	"fmt"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	// ErrKeyExternalID is key in Errors map indicating there's a problem with External ID of the resource.
	ErrKeyExternalID = "external_id"
	// ErrKeyTransactionExternalID is key in Errors map indicating there's a problem with External ID of the transaction.
	ErrKeyTransactionExternalID = "transactions.external_id"
	// ErrKeyLineItemsExternalID indicates problem with one/any of line items' external IDs
	ErrKeyLineItemsExternalID = "line_items.external_id"
	// ErrKeyName - data source name
	ErrKeyName = "name"
	// ErrValCustomerExternalIDExists = can't import new customer with the same external ID
	ErrValCustomerExternalIDExists = "The external ID for this customer already exists in our system."
	// ErrValLineItemExternalIDExists = can't import invoice, b'c line item external ID exists
	ErrValLineItemExternalIDExists = "The external ID for this line item already exists in our system."
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
	CreateDataSourceWithSystem(dataSource *DataSource) (*DataSource, error)
	RetrieveDataSource(dataSourceUUID string) (*DataSource, error)
	ListDataSources() (*DataSources, error)
	ListDataSourcesWithFilters(listDataSourcesParams *ListDataSourcesParams) (*DataSources, error)
	PurgeDataSource(dataSourceUUID string) error
	EmptyDataSource(dataSourceUUID string) error
	DeleteDataSource(dataSourceUUID string) error
	// Invoices
	CreateInvoices(invoices []*Invoice, customerUUID string) (*Invoices, error)
	ListInvoices(cursor *Cursor, customerUUID string) (*Invoices, error)
	ListAllInvoices(listAllInvoicesParams *ListAllInvoicesParams) (*Invoices, error)
	RetrieveInvoice(invoiceUUID string) (*Invoice, error)
	DeleteInvoice(invoiceUUID string) error
	// Plans
	CreatePlan(plan *Plan) (result *Plan, err error)
	RetrievePlan(planUUID string) (*Plan, error)
	ListPlans(listPlansParams *ListPlansParams) (*Plans, error)
	UpdatePlan(plan *Plan, planUUID string) (*Plan, error)
	DeletePlan(planUUID string) error
	// Plan Groups
	CreatePlanGroup(planGroup *PlanGroup) (result *PlanGroup, err error)
	RetrievePlanGroup(planGroupUUID string) (*PlanGroup, error)
	ListPlanGroups(cursor *Cursor) (*PlanGroups, error)
	UpdatePlanGroup(plan *PlanGroup, planGroupUUID string) (*PlanGroup, error)
	DeletePlanGroup(planGroupUUID string) error
	ListPlanGroupPlans(cursor *Cursor, planGroupUUID string) (*PlanGroupPlans, error)
	// Subscriptions
	CancelSubscription(subscriptionUUID string, cancelSubscriptionParams *CancelSubscriptionParams) (*Subscription, error)
	ListSubscriptions(cursor *Cursor, customerUUID string) (*Subscriptions, error)
	// Transactions
	CreateTransaction(transaction *Transaction, invoiceUUID string) (*Transaction, error)

	// Customers
	CreateCustomer(newCustomer *NewCustomer) (*Customer, error)
	RetrieveCustomer(customerUUID string) (*Customer, error)
	UpdateCustomer(Customer *Customer, customerUUID string) (*Customer, error)
	UpdateCustomerV2(Customer *UpdateCustomer, customerUUID string) (*Customer, error)
	ListCustomers(ListCustomersParams *ListCustomersParams) (*Customers, error)
	SearchCustomers(SearchCustomersParams *SearchCustomersParams) (*Customers, error)
	MergeCustomers(MergeCustomersParams *MergeCustomersParams) error
	DeleteCustomer(customerUUID string) error
	DeleteCustomerInvoices(dataSourceUUID, customerUUID string) error
	DeleteCustomerInvoicesV2(dataSourceUUID, customerUUID string, DeleteCustomerInvoicesParams *DeleteCustomerInvoicesParams) error
	ListCustomersContacts(ListContactsParams *ListContactsParams, customerUUID string) (*Contacts, error)
	CreateCustomersContact(newContact *NewContact, customerUUID string) (*Contact, error)

	// Contacts
	CreateContact(newContact *NewContact) (*Contact, error)
	RetrieveContact(contactUUID string) (*Contact, error)
	UpdateContact(Contact *UpdateContact, contactUUID string) (*Contact, error)
	ListContacts(ListContactsParams *ListContactsParams) (*Contacts, error)
	DeleteContact(contactUUID string) error
	MergeContacts(intoContactUUID string, fromContactUUID string) (*Contact, error)

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
	MetricsListCustomerSubscriptions(cursor *Cursor, customerUUID string) (*MetricsCustomerSubscriptions, error)
	MetricsListCustomerActivities(cursor *Cursor, customerUUID string) (*MetricsCustomerActivities, error)
	MetricsListActivities(MetricsListActivitiesParams *MetricsListActivitiesParams) (*MetricsActivities, error)
	MetricsCreateActivitiesExport(CreateMetricsActivitiesExportParam *CreateMetricsActivitiesExportParam) (*MetricsActivitiesExport, error)
	MetricsRetrieveActivitiesExport(activitiesExportUUID string) (*MetricsActivitiesExport, error)

	// Account
	RetrieveAccount() (*Account, error)

	// Subscription Events
	ListSubscriptionEvents(filters *FilterSubscriptionEvents, cursor *Cursor) (*SubscriptionEvents, error)
	CreateSubscriptionEvent(newSubscriptionEvent *SubscriptionEvent) (*SubscriptionEvent, error)
	UpdateSubscriptionEvent(subscriptionEvent *SubscriptionEvent) (*SubscriptionEvent, error)
	DeleteSubscriptionEvent(deleteParams *DeleteSubscriptionEvent) error
}

// API is the handle for communicating with Chartmogul.
type API struct {
	ApiKey string
	Client *http.Client
}

// AnchorCursor contains query parameters for anchor based pagination used for some APIs in ChartMogul.
type AnchorCursor struct {
	PerPage uint32 `json:"per-page,omitempty"`
	//StartAfter is used to get the next set of Entries and its value should be the UUID of last Entry from previous response.
	StartAfter string `json:"start-after,omitempty"`
}

// Pagination is a struct to handle the new cursor based pagination pagination response.
type Pagination struct {
	// Cursor is a reference to get the next set of entries.
	Cursor  string `json:"cursor,omitempty"`
	HasMore bool   `json:"has_more,omitempty"`
}

// Cursor is the new standard for cursor with pagination.
type Cursor struct {
	PerPage uint32 `json:"per_page,omitempty"`
	// Cursor is a reference to get the next set of entries.
	Cursor string `json:"cursor,omitempty"`
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
		msg = e[ErrKeyName]
	}
	return msg == ErrValExternalIDExists ||
		msg == ErrValHasAlreadyBeenTaken ||
		msg == ErrValCustomerExternalIDExists ||
		msg == ErrValPlanExternalIDExists ||
		msg == ErrValInvoiceExternalIDExists
}

// IsInvoiceAndItsEntitiesAlreadyExist returns true if:
// * invoice already exists AND
// * ANY other entities (line items, transactions) already exist AND
// * no other error
//
// So, eg. if the invoice doesn't exist, but the transaction does,
// you have a different problem (duplicating txns) and this returns false.
func (e Errors) IsInvoiceAndItsEntitiesAlreadyExist() (is bool) {
	if e == nil {
		return
	}
	if msg := e[ErrKeyExternalID]; msg != ErrValInvoiceExternalIDExists {
		return
	}
	for key, val := range e {
		switch key {
		case ErrKeyTransactionExternalID:
			if val != ErrValExternalIDExists {
				return
			}
		case ErrKeyLineItemsExternalID:
			if val != ErrValLineItemExternalIDExists {
				return
			}
		case ErrKeyExternalID:
			// already checked
		default:
			return
		}
	}

	return true
}

// IsInvoiceAndTransactionAlreadyExist occurs when both invoice and tx exist already.
// Use `IsInvoiceAndItsEntitiesAlreadyExist` if you'd like to catch line items as well.
func (e Errors) IsInvoiceAndTransactionAlreadyExist() (is bool) {
	if e == nil {
		return
	}
	msg1, ok1 := e[ErrKeyExternalID]
	msg2, ok2 := e[ErrKeyTransactionExternalID]
	return ok1 && ok2 && len(e) == 2 &&
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

// SetClient changes the client - for VCR integration tests.
func (api *API) SetClient(newClient *http.Client) {
	api.Client = newClient
}

func prepareURL(path string) string {
	return fmt.Sprintf(url, path)
}

func (api API) req(req *gorequest.SuperAgent) *gorequest.SuperAgent {
	// defaults for client go here:
	if api.Client != nil {
		req.Client = api.Client
	}
	return req.Timeout(timeout).
		SetBasicAuth(api.ApiKey, "").
		Set("Content-Type", "application/json").
		Set("User-Agent", "chartmogul-go/"+Version)
}
