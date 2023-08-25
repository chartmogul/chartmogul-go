package integration

import (
	"net/http"
	"os"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v4"
	"github.com/parnurzeal/gorequest"
)

// TestListSubscriptions is an integration test that verifies the listing
// of subscriptions using the ChartMogul API.
func TestListSubscriptions(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := NewRecorderWithAuthFilter("./fixtures/subscriptions")
	if err != nil {
		t.Fatalf("Failed to initialize recorder: %v", err)
	}

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
	_ = createTestInvoicesForCustomer(api, *cus, plan.UUID, t)

	// Subscriptions are created through background jobs and we need to wait for them to be processed.
	// Only enable when hitting the real API
	// time.Sleep(120 * time.Second)

	validateListSubscriptions(api, cus.UUID, t)
}

// validateListSubscriptions validates that subscriptions can be correctly listed using the API.
func validateListSubscriptions(api *cm.API, cusUUID string, t *testing.T) {
	listSubscriptions, err := api.ListSubscriptions(&cm.Cursor{PerPage: 1}, cusUUID)
	if err != nil {
		t.Fatalf("Failed to list subscriptions: %v", err)
	}

	if listSubscriptions.HasMore != true {
		t.Fatalf("Expected HasMore to be true, got: %v", listSubscriptions.HasMore)
	}

	listSubscriptionsNextPage, err := api.ListSubscriptions(&cm.Cursor{Cursor: listSubscriptions.Cursor, PerPage: 1}, cusUUID)
	if err != nil {
		t.Fatalf("Failed to list subscriptions: %v", err)
	}

	if listSubscriptionsNextPage.HasMore != false {
		t.Fatalf("Expected HasMore to be false, got: %v", listSubscriptionsNextPage.HasMore)
	}

	entry1 := listSubscriptions.Subscriptions[0]
	entry2 := listSubscriptionsNextPage.Subscriptions[0]
	if entry1.UUID == entry2.UUID {
		t.Fatalf("Expected subscriptions to be different. Got UUIDs %s & %s", entry1.UUID, entry2.UUID)
	}
}
