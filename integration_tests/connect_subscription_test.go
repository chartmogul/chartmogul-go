package integration

import (
	"log"
	"net/http"
	"testing"
	"time"

	cm "github.com/chartmogul/chartmogul-go"
	"github.com/parnurzeal/gorequest"

	"github.com/dnaeon/go-vcr/recorder"
)

func TestConnectSubscriptions(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := recorder.New("./fixtures/connect_subscriptions")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() // Make sure recorder is stopped once done with it

	api := &cm.API{}
	api.SetClient(&http.Client{Transport: r})
	gorequest.DisableTransportSwap = true

	ds, err := api.CreateDataSource("Test Connect Subscriptions")
	if err != nil {
		t.Fatal(err)
	}
	defer api.DeleteDataSource(ds.UUID)

	plan1, err := api.CreatePlan(&cm.Plan{
		DataSourceUUID: ds.UUID,
		ExternalID:     "ext_plan_1",
		Name:           "test plan 1",
		IntervalCount:  1,
		IntervalUnit:   "month",
	})
	if err != nil {
		t.Fatal(err)
	}
	plan2, err := api.CreatePlan(&cm.Plan{
		DataSourceUUID: ds.UUID,
		ExternalID:     "ext_plan_2",
		Name:           "test plan 2",
		IntervalCount:  3,
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
	invoice1 := &cm.Invoice{
		Date:       time.Now().Format(time.RFC3339),
		ExternalID: "INV_to_be_retrieved",
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
				PlanUUID:               plan1.UUID,
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
	}
	invoice2 := &cm.Invoice{
		Date:       time.Now().Format(time.RFC3339),
		ExternalID: "INV_to_be_retrieved_2",
		Currency:   "EUR",
		LineItems: []*cm.LineItem{
			{
				Type:                   "subscription",
				AmountInCents:          10000,
				ExternalID:             "ext_line_item_2",
				SubscriptionExternalID: "ext_subscription_2",
				PlanUUID:               plan2.UUID,
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
	}
	_, err = api.CreateInvoices([]*cm.Invoice{invoice1, invoice2}, cus.UUID)
	if err != nil {
		t.Fatal(err)
	}
	subs, err := api.ListSubscriptions(&cm.Cursor{PerPage: 200}, cus.UUID)
	if err != nil {
		t.Fatal(err)
	}
	//time.Sleep(time.Minute)
	_, err = api.CancelSubscription(subs.Subscriptions[0].UUID, &cm.CancelSubscriptionParams{CancelledAt: "2017-05-15T00:00:00Z"})
	if err != nil {
		t.Fatal(err)
	}
	//time.Sleep(time.Minute)
	err = api.ConnectSubscriptions(cus.UUID, subs.Subscriptions)
	if err != nil {
		t.Fatal(err)
	}
	//time.Sleep(time.Minute)
	msubs, err := api.MetricsListSubscriptions(&cm.Cursor{PerPage: 200}, cus.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if len(msubs.Entries) != 1 {
		t.Fatal("expected one subscription")
	}
	if msubs.Entries[0].ExternalID != subs.Subscriptions[1].ExternalID {
		t.Fatal("subscription external IDs do not match")
	}
}
