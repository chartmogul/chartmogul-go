package integration

import (
	"net/http"
	"os"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v4"
	"github.com/parnurzeal/gorequest"
)

// TestListCustomerSubscriptions is an integration test that verifies the listing,
// of subscriptions based on customer using the ChartMogul API.
func TestListCustomerSubscriptions(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := NewRecorderWithAuthFilter("./fixtures/customer_subscriptions")
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
	_ = createTestInvoicesForcustomer(api, ds.UUID, *cus, plan.UUID, t)

	// Subscriptions are created through background jobs and we need to wait for them to be processed.
	// Only enable when hitting the real API
	// time.Sleep(120 * time.Second)

	validateListCustomerSubscriptions(api, cus.UUID, t)
}

// validateListCustomerSubscriptions validates that subscriptions per customer can be correctly listed using the API.
func validateListCustomerSubscriptions(api *cm.API, cusUUID string, t *testing.T) {
	listCustomerSubscriptions, err := api.MetricsListCustomerSubscriptions(&cm.Cursor{PerPage: 1}, cusUUID)
	if err != nil {
		t.Fatalf("Failed to list customer subscriptions: %v", err)
	}

	if listCustomerSubscriptions.HasMore != true {
		t.Fatalf("Expected HasMore to be true, got: %v", listCustomerSubscriptions.HasMore)
	}

	listCustomerSubscriptionsNextPage, err := api.MetricsListCustomerSubscriptions(&cm.Cursor{Cursor: listCustomerSubscriptions.Cursor, PerPage: 1}, cusUUID)
	if err != nil {
		t.Fatalf("Failed to list customer subscriptions: %v", err)
	}

	if listCustomerSubscriptionsNextPage.HasMore != false {
		t.Fatalf("Expected HasMore to be false, got: %v", listCustomerSubscriptionsNextPage.HasMore)
	}

	entry1 := listCustomerSubscriptions.Entries[0]
	entry2 := listCustomerSubscriptionsNextPage.Entries[0]
	if entry1.ID == entry2.ID {
		t.Fatalf("Expected subscriptions to be different. Got UUIDs %d & %d", entry1.ID, entry2.ID)
	}
}
