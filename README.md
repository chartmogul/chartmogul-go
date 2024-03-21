<p align="center">
<a href="https://chartmogul.com"><img width="200" src="https://user-images.githubusercontent.com/5329361/42206299-021e4184-7ea7-11e8-8160-8ecd5d9948b8.png"></a>
</p>

<h3 align="center">Official ChartMogul API Go Client</h3>

<p align="center"><code>chartmogul-go</code> provides convenient Golang bindings for <a href="https://dev.chartmogul.com">ChartMogul's API</a>.</p>
<p align="center">
  <a href="https://github.com/chartmogul/chartmogul-go/tree/v2"><img src="https://github.com/chartmogul/chartmogul-go/actions/workflows/test.yml/badge.svg?branch=v2" alt="Build Status"/></a>
</p>
<hr>

<p align="center">
<b><a href="#installation">Installation</a></b>
|
<b><a href="#configuration">Configuration</a></b>
|
<b><a href="#usage">Usage</a></b>
|
<b><a href="#development">Development</a></b>
|
<b><a href="#contributing">Contributing</a></b>
|
<b><a href="#license">License</a></b>
</p>
<hr>
<br>

[![Go Reference](https://pkg.go.dev/badge/github.com/chartmogul/chartmogul-go/v4.svg)](https://pkg.go.dev/github.com/chartmogul/chartmogul-go/v4)
[![Go Report Card](https://goreportcard.com/badge/github.com/chartmogul/chartmogul-go)](https://goreportcard.com/report/github.com/chartmogul/chartmogul-go)

## Installation

This library requires Go 1.11 or above.

```sh
go get github.com/chartmogul/chartmogul-go/v4
```

## Configuration
[Deprecation] - `account_token`/`secret_key` combo is deprecated. Please use API key for both fields.
Version 3.x will introduce a breaking change in authentication configuration. For more details, please visit: https://dev.chartmogul.com

Version 4.x will introduce a breaking change for pagination on List endpoints. The `cursor` object now expects a `per_page` and `cursor` parameter. `page` will no longer be accepted.

First create the `API` struct by passing your API key, available from the administration section of your ChartMogul account.

```go
import cm "github.com/chartmogul/chartmogul-go/v4"

api := cm.API{
    ApiKey: os.Getenv("CHARTMOGUL_API_KEY"),
}

// Try authentication
ok, err := api.Ping()
if err != nil {
    fmt.Printf("This didn't work out: %v", err)
}
```

This struct has all the methods you can use to interact with ChartMogul.

### HTTP 2
ChartMogul's current stable version of nginx is incompatible with HTTP 2
implementation of Go as of 1.7.3.
For this reason the application must run with the following
(or otherwise prohibit HTTP 2):
```bash
export GODEBUG=http2client=0
```

## Usage

### Rate Limits
The library retries on HTTP 429 (rate limit reached), so that you don't have to manually handle rate limiting.
See:
[ChartMogul: Rate Limits](https://dev.chartmogul.com/docs/rate-limits) &amp;
[BackOff constants](https://godoc.org/github.com/cenkalti/backoff#pkg-constants).
Exponential back-off algorithm is used.
It also automatically retries on network errors, eg. when the server can't be reached.

The API calls will retry automatically and block (several minutes), therefore it's still advisable
to only use reasonable parallelism. In case it keeps failing after maximum retry period, it will
return the HTTP 429 error.

Note: the `Ping` doesn't retry.

### Import API

Available methods in Import API:

#### [Data Sources](https://dev.chartmogul.com/docs/data-sources)

```go
api.CreateDataSource("name")
api.ListDataSources()
api.RetrieveDataSource("uuid")
api.DeleteDataSource("uuid")
```

#### [Customers](https://dev.chartmogul.com/docs/retrieve-customer)

```go
api.CreateCustomer(&cm.NewCustomer{})
api.RetrieveCustomer("customerUUID")
api.SearchCustomers(&cm.SearchCustomersParams{})
api.ListCustomers(&cm.ListCustomersParams{})
api.UpdateCustomer(&cm.NewCustomer{}, "customerUUID")
api.MergeCustomers(&cm.MergeCustomersParams{})
api.ConnectSubscriptions("customerUUID", []cm.Subscription{})
api.ListCustomersContact(&cm.ListContactsParams{}, "customerUUID")
api.CreateCustomersContact(&cm.NewContact{}, "customerUUID")
api.ListCustomerNotes(&cm.ListNotesParams{}, "customerUUID")
api.CreateCustomerNote(&cm.NewNote{}, "customerUUID")
api.ListCustomerOpportunities(&cm.ListOpportunitiesParams{}, "customerUUID")
api.CreateCustomerOpportunity(&cm.NewOpportunity{}, "customerUUID")
```

#### [Contacts](https://dev.chartmogul.com/reference/contacts)

```go
api.CreateContact(&cm.NewContact{})
api.RetrieveContact("customerUUID")
api.ListContacts(&cm.ListContactsParams{})
api.UpdateContact(&cm.UpdateContact{}, "contact")
api.DeleteContact("customerUUID")
api.MergeContacts("intoContactUUID", "fromContactUUID")
```

#### [Customer Notes](https://dev.chartmogul.com/reference/customer-notes)

```go
api.CreateNote(&cm.NewNote{})
api.RetrieveNote("noteUUID")
api.ListNote(&cm.ListNoteParams{})
api.UpdateNote(&cm.UpdateNote{}, "noteUUID")
api.DeleteNote("noteUUID")
```

#### [Opportunities](https://dev.chartmogul.com/reference/opportunities)

```go
api.CreateOpportunity(&cm.NewOpportunity{})
api.RetrieveOpportunity("opportunityUUID")
api.ListOpportunities(&cm.ListOpportunitiesParams{})
api.UpdateOpportunity(&cm.UpdateOpportunity{}, "opportunityUUID")
api.DeleteOpportunity("opportunityUUID")
```

#### [Plans](https://dev.chartmogul.com/reference#plans)

```go
api.CreatePlan(&cm.Plan{Name: "name", ExternalID: "external_id"}, "dataSourceUUID")
api.RetrievePlan("planUUID")
api.ListPlans(&cm.ListPlansParams{Cursor: cm.Cursor{PerPage: "10"}})
api.UpdatePlan(&cm.Plan{}, "planUUID")
api.DeletePlan("planUUID")
```

#### [Plan Groups](https://dev.chartmogul.com/reference#plan_groups)

```go
api.CreatePlanGroup(&cm.PlanGroup{Name: "name", Plans: []*string{&planOne.UUID, &planTwo.UUID}})
api.RetrievePlanGroup("planGroupUUID")
api.ListPlanGroups(&cm.ListPlansParams{Cursor: cm.Cursor{PerPage: "10"}})
api.UpdatePlanGroup(&cm.PlanGroup{}, "planGroupUUID")
api.DeletePlanGroup("planGroupUUID")
api.ListPlanGroupPlans(&cm.ListPlansParams{Cursor: cm.Cursor{PerPage: "10"}},  "planGroupUUID")
```

#### [Invoices](https://dev.chartmogul.com/docs/invoices)

```go
api.CreateInvoices([]*cm.Invoice{*cm.Invoice{}}, "customerUUID")
api.ListInvoices(&cm.Cursor{}, "customerUUID")
api.ListAllInvoices(&cm.ListAllInvoicesParams{})
api.RetrieveInvoice("invoiceUUID")
api.DeleteInvoice("invoiceUUID")
```

#### [Transactions](https://dev.chartmogul.com/docs/transactions)

```go
api.CreateTransaction(&cm.Transaction{}, "invoiceUUID")
```

#### [Subscriptions](https://dev.chartmogul.com/docs/subscriptions)

```go
api.CancelSubscription("subscriptionUUID", &cm.CancelSubscriptionParams{CancelledAt: "2005-01-01T01:02:03.000Z"})
api.CancelSubscription("subscriptionUUID", &cm.CancelSubscriptionParams{CancellationDates: []string{"2005-01-01T01:02:03.000Z", "2006-10-21T11:21:13.000Z"}})
api.ListSubscriptions(&cm.Cursor{}, "customerUUID")
```

#### [Customer Attributes](https://dev.chartmogul.com/docs/customer-attributes)

```go
api.RetrieveCustomersAttributes("customerUUID")
```

#### [Tags](https://dev.chartmogul.com/docs/tags)

```go
api.AddTagsToCustomer("customerUUID", []string{})
api.AddTagsToCustomersWithEmail("email@customer.com", []string{})
```


#### [Custom Attributes](https://dev.chartmogul.com/docs/custom-attributes)

```go
api.AddCustomAttributesToCustomer("customerUUID", []*cm.CustomAttribute{})
```

#### [Subscription Events](https://dev.chartmogul.com/reference/subscription-events)
```go
api.ListSubscriptionEvents(filters *FilterSubscriptionEvents, cursor *Cursor)
api.CreateSubscriptionEvent(newSubscriptionEvent *SubscriptionEvent)
api.UpdateSubscriptionEvent(subscriptionEvent *SubscriptionEvent)
api.DeleteSubscriptionEvent(deleteParams *DeleteSubscriptionEvent)
```

### [Metrics API](https://dev.chartmogul.com/docs/introduction-metrics-api)

Available methods in Metrics API:


```go
api.MetricsRetrieveAll(&MetricsFilter{})
api.MetricsRetrieveMRR(&MetricsFilter{})
api.MetricsRetrieveARR(&MetricsFilter{})
api.MetricsRetrieveARPA(&MetricsFilter{})
api.MetricsRetrieveASP(&MetricsFilter{})
api.MetricsRetrieveCustomerCount(&MetricsFilter{})
api.MetricsRetrieveCustomerChurnRate(&MetricsFilter{})
api.MetricsRetrieveMRRChurnRate(&MetricsFilter{})
api.MetricsRetrieveLTV(&MetricsFilter{})

api.MetricsListCustomerSubscriptions(&Cursor{}, "customerUUID")
api.MetricsListCustomerActivities(&Cursor{}, "customerUUID")

api.MetricsListActivities(&cm.MetricsListActivitiesParams{StartDate: "2016-09-16", AnchorCursor: cm.AnchorCursor{PerPage: 5, StartAfter: "b45b1d3f-3823-424f-ab47-5a1d0c00a7f6"}})

api.MetricsCreateActivitiesExport(&cm.CreateMetricsActivitiesExportParam{StartDate: "2016-09-16",Type: "contraction"})
api.MetricsRetrieveActivitiesExport("activitiesExportUUID")
```

### Account

Availiable methods:

```go
api.RetrieveAccount()
```


### Errors

The library returns parsed errors inside the structs as from the REST API,
which is handy eg. when you upload multiple invoices - to know, which one had issues.

```go
type Errors map[string]string
```

Non-2xx statuses and network problems will also be returned as a value `error`
for standard binary ok/problem handling:

```go
_, err := api.ImportListPlans(nil)
if err != nil {
    // ...
}
```

Such errors are `HTTPError` wrapped in [errors with stack](https://github.com/pkg/errors).
```go
type HTTPError interface {
	StatusCode() int
	Status() string
	Response() string
}
```

```go
// If you want to check specific HTTP status:
switch e := errors.Cause(err).(type) {
case cm.HTTPError:
    if e.StatusCode() == 422 {
        // special reaction
    }
}
```

If there are network/TLS issues it will be `RequestErrors interface`.
This has the method `Errors() []error`.

For fine-grain reaction you can use the parsed errors from
API in the primary return structures, when there's an `Errors` or `Error` field.
You can use hepler methods to spare code on checking map contents:
* `Errors.IsAlreadyExists()`
* `Errors.IsInvoiceAndTransactionAlreadyExist()`

## Development

To work on the library:

* Fork it
* Create your feature branch (`git checkout -b my-new-feature`)
* Install dependencies: `go install`
* Fix bugs or add features. Make sure the changes pass the Go coding standards.
* Push to the branch (`git push origin my-new-feature`)
* Create a new Pull Request

### Recommended

* [`github.com/davecgh/go-spew/spew`](https://github.com/davecgh/go-spew/spew) for debugging data (reference documentation output done using this library).
* add pre-commit hook `go test ./...` (in `.git/hooks/pre-commit`) to have a working state always.

### Testing
* Use `net/http/httptest` for mocking HTTP server directly, see file `generic_test.go` for examples.
* For `integration_tests` against real API use `github.com/dnaeon/go-vcr` library. Be careful to remove your API credentials from fixtures before committing! If Import API App changes, re-record the affected integration tests (by deleting fixtures).

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/chartmogul/chartmogul-go.

## License

The library is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).

### The MIT License (MIT)

*Copyright (c) 2017 ChartMogul Ltd.*

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
