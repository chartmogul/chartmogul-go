<p align="center">
<a href="https://chartmogul.com"><img width="200" src="https://chartmogul.com/assets/img/logo.png"></a>
</p>

<h3 align="center">Official ChartMogul API Go Client</h3>

<p align="center"><code>chartmogul-go</code> provides convenient Golang bindings for <a href="https://dev.chartmogul.com">ChartMogul's API</a>.</p>

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

## Warning
Library in testing stage and subject to possibly great changes!

## Installation

This library requires Go 1.7.3 or above.

```sh
go get github.com/chartmogul/chartmogul-go
```

## Configuration

First create the `API` struct by passing your account token and secret key, available from the administration section of your ChartMogul account.

```go
import cm "github.com/chartmogul/chartmogul-go"

api := cm.API{
    AccountToken: os.Getenv("CHARTMOGUL_ACCOUNT_TOKEN"),
    AccessKey:    os.Getenv("CHARTMOGUL_SECRET_KEY"),
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

### Import API

Available methods in Import API:

#### [Data Sources](https://dev.chartmogul.com/docs/data-sources)

```go
api.ImportCreateDataSource("name")
api.ImportListDataSources()
api.ImportDeleteDataSource("uuid")
```

#### [Customers](https://dev.chartmogul.com/docs/customers)

```go
api.ImportCreateCustomer(&cm.Customer{Name: "name", ExternalID: "external_id"}, "dataSourceUUID")
api.ImportListCustomers(&cm.ImportListCustomersParams{Cursor: cm.Cursor{Page: "1", PerPage: "10"}})
api.ImportDeleteCustomer("uuid")
```

#### [Plans](https://dev.chartmogul.com/docs/plans)

```go
api.ImportCreatePlan(&cm.Plan{Name: "name", ExternalID: "external_id"}, "dataSourceUUID")
api.ImportListPlans(&cm.ListPlansParams{Cursor: cm.Cursor{Page: "1", PerPage: "10"}})
```

#### [Invoices](https://dev.chartmogul.com/docs/invoices)

```go
api.ImportCreateInvoices([]*cm.Invoice{*cm.Invoice{}}, "customerUUID")
api.ImportListInvoices(&cm.Cursor{}, "customerUUID")
```

#### [Transactions](https://dev.chartmogul.com/docs/transactions)

```go
api.ImportCreateTransaction(&cm.Transaction{}, "invoiceUUID")
```

#### [Subscriptions](https://dev.chartmogul.com/docs/subscriptions)

```go
api.ImportCancelSubscription("subscriptionUUID", "2005-01-01T01:02:03.000Z")
api.ImportListSubscriptions(&cm.Cursor{}, "customerUUID")
```

### Enrichment API

Available methods in Enrichment API:


#### [Customers](https://dev.chartmogul.com/docs/retrieve-customer)

```go
api.EnrichmentCreateCustomer(&NewCustomer{})
api.EnrichmentRetrieveCustomer("customerUUID")
api.EnrichmentSearchCustomers(&cm.EnrichmentSearchCustomersParams{})
api.EnrichmentListCustomers(&cm.EnrichmentListCustomersParams{})
api.EnrichmentUpdateCustomer(&cm.EnrichmentCustomer{}, "customerUUID")
api.EnrichmentMergeCustomers(&cm.EnrichmentMergeCustomersParams{})
```

#### [Customer Attributes](https://dev.chartmogul.com/docs/customer-attributes)

```go
api.EnrichmentRetrieveCustomersAttributes("customerUUID")
```

#### [Tags](https://dev.chartmogul.com/docs/tags)

```go
api.EnrichmentAddTagsToCustomer("customerUUID", []string{})
api.EnrichmentAddTagsToCustomersWithEmail("email@customer.com", []string{})
```


#### [Custom Attributes](https://dev.chartmogul.com/docs/custom-attributes)

```go
api.EnrichmentAddCustomAttributesToCustomer("customerUUID", []*cm.CustomAttribute{})
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

api.MetricsListSubscriptions(&Cursor{}, "customerUUID")
api.MetricsListActivities(&Cursor{}, "customerUUID")
```


### Errors

The library returns parsed errors inside the structs as from the REST API,
which is handy eg. when you upload multiple invoices - to know, which one had issues.

```go
type Errors map[string]string
```

Non-2xx statuses and network problems will also be returned as a second parameter for standard:

```go
_, err := api.ImportListPlans(nil)
if err != nil {
    // ...
}
```

These errors can be either `HTTPError` or `RequestErrors` (network issues).
The `HTTPError.StatusCode == 404` is a coarse possibility to react on errors like
HTTP 404 Not Found.

For fine-grain reaction you can use the parsed errors from
API in the primary return structures, when there's an `Errors` field.
You can use hepler methods like `Errors.IsAlreadyExists()` or
`Errors.IsInvoiceAndTransactionAlreadyExist()` to spare code on checking map contents.

## Development

To work on the library:

* Fork it
* Create your feature branch (`git checkout -b my-new-feature`)
* Install dependencies: `go install`
* Fix bugs or add features. Make sure the changes pass the Go coding standards.
* Push to the branch (`git push origin my-new-feature`)
* Create a new Pull Request

## Contributing

Bug reports and pull requests are welcome on GitHub at https://github.com/chartmogul/chartmogul-go.

## License

The library is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).

### The MIT License (MIT)

*Copyright (c) 2017 ChartMogul Ltd.*

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
