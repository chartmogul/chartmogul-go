package integration

import (
	"net/http"
	"os"
	"testing"
	"time"

	cm "github.com/chartmogul/chartmogul-go/v3"
	"github.com/parnurzeal/gorequest"
)

// TestListInvoicesIntegration is an integration test that verifies the creation
// & listing of invoices using the ChartMogul API.
func TestListInvoicesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := NewRecorderWithAuthFilter("./fixtures/invoices")
	if err != nil {
		t.Fatalf("Failed to initialize recorder: %v", err)
	}
	t.Cleanup(func() { r.Stop() })

	api := &cm.API{
		ApiKey: os.Getenv("CHARTMOGUL_API_KEY"),
		Client: &http.Client{Transport: r},
	}
	gorequest.DisableTransportSwap = true

	// Create a new data source.
	ds, err := api.CreateDataSource("testing")
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

	// Create and verify test invoices.
	invoices := createTestInvoices(api, ds.UUID, cus.UUID, plan.UUID, t)

	validateListInvoices(api, ds.UUID, cus.UUID, invoices, t)
}

// createTestInvoices creates a set of test invoices for the integration test.
// It returns the created invoices.
func createTestInvoices(api *cm.API, dsUUID, cusUUID, planUUID string, t *testing.T) *cm.Invoices {
	testInvoices := []*cm.Invoice{
		{
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
					Type:                      "subscription",
					AmountInCents:             10000,
					ExternalID:                "ext_line_item",
					SubscriptionExternalID:    "ext_subscription",
					SubscriptionSetExternalID: "ext_subscription_set",
					PlanUUID:                  planUUID,
					Quantity:                  10,
					ServicePeriodStart:        "2017-05-01T00:00:00.000Z",
					ServicePeriodEnd:          "2017-05-31T00:00:00.000Z",
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
		{
			Date:       time.Now().Format(time.RFC3339),
			ExternalID: "INV_to_be_retrieved1",
			Currency:   "EUR",
			LineItems: []*cm.LineItem{
				{
					Type:          "one_time",
					AmountInCents: 4200,
					Description:   "fake_item1",
					Quantity:      2,
				},
				{
					Type:                      "subscription",
					AmountInCents:             11000,
					ExternalID:                "ext_line_item1",
					SubscriptionExternalID:    "ext_subscription1",
					SubscriptionSetExternalID: "ext_subscription_set1",
					PlanUUID:                  planUUID,
					Quantity:                  11,
					ServicePeriodStart:        "2017-05-01T00:00:00.000Z",
					ServicePeriodEnd:          "2017-05-31T00:00:00.000Z",
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
	}

	inv, err := api.CreateInvoices(testInvoices, cusUUID)
	if err != nil {
		t.Fatal(err)
	}

	return inv
}

// validateListinvoices validates that the created invoices can be correctly listed using the API.
func validateListInvoices(api *cm.API, dsUUID, cusUUID string, invoices *cm.Invoices, t *testing.T) {
	invoicesList, err := api.ListInvoices(&cm.PaginationWithCursor{PerPage: 1}, cusUUID)
	if err != nil {
		t.Fatalf("Failed to list invoices: %v", err)
	}

	if invoicesList.HasMore != true {
		t.Fatalf("Expected HasMore to be true, got: %v", invoices.HasMore)
	}

	invoicesListAll, err := api.ListInvoices(&cm.PaginationWithCursor{PerPage: 2}, cusUUID)
	if err != nil {
		t.Fatalf("Failed to list invoices: %v", err)
	}

	if invoicesListAll.HasMore != false {
		t.Fatalf("Expected HasMore to be false, got: %v", invoices.HasMore)
	}
}
