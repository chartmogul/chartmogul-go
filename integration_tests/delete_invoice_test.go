package integration

import (
	"log"
	"net/http"
	"testing"
	"time"

	cm "github.com/chartmogul/chartmogul-go/v2"
	"github.com/davecgh/go-spew/spew"
	"github.com/parnurzeal/gorequest"

	"github.com/dnaeon/go-vcr/recorder"
)

func TestDeleteInvoice(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := recorder.New("./fixtures/delete_invoice")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() //nolint

	api := &cm.API{Client: &http.Client{Transport: r}}
	gorequest.DisableTransportSwap = true

	ds, err := api.CreateDataSource("Test Delete Invoice 1")
	if err != nil {
		t.Fatal(err)
	}
	defer api.DeleteDataSource(ds.UUID) //nolint

	plan, err := api.CreatePlan(&cm.Plan{
		DataSourceUUID: ds.UUID,
		ExternalID:     "ext_plan",
		Name:           "test plan",
		IntervalCount:  1,
		IntervalUnit:   "month",
	})
	if err != nil {
		t.Fatal(err)
	}

	cus, err := api.CreateCustomer(&cm.NewCustomer{
		Name:           "Test Delete Invoice",
		Email:          "petr@chartmogul.com",
		ExternalID:     "ext_customer",
		DataSourceUUID: ds.UUID,
	})
	if err != nil {
		t.Fatal(err)
	}

	inv, err := api.CreateInvoices([]*cm.Invoice{
		{
			Date:       time.Now().Format(time.RFC3339),
			ExternalID: "INV_to_be_deleted",
			Currency:   "EUR",
			LineItems: []*cm.LineItem{
				{
					Type:          "one_time",
					AmountInCents: 4500,
					Description:   "fake_item",
					Quantity:      2,
				},
				{
					Type:                   "subscription",
					AmountInCents:          10000,
					ExternalID:             "ext_line_item",
					SubscriptionExternalID: "ext_subscription",
					PlanUUID:               plan.UUID,
					Quantity:               10,
					ServicePeriodStart:     "2017-05-01T00:00:00.000Z",
					ServicePeriodEnd:       "2017-05-31T00:00:00.000Z",
				},
			},
			Transactions: []*cm.Transaction{
				{
					Date:   time.Now().Add(3 * time.Hour).Format(time.RFC3339),
					Result: "successful",
					Type:   "payment",
				},
			},
		},
	}, cus.UUID)
	if err != nil {
		t.Fatal(err)
	}
	err = api.DeleteInvoice(inv.Invoices[0].UUID)
	if err != nil {
		t.Fatal(err)
	}
	invoices, err := api.ListInvoices(nil, cus.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if len(invoices.Invoices) != 0 {
		spew.Dump(invoices)
		t.Fatal("Should have been deleted.")
	}
}
