package integration

import (
	"net/http"
	"os"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v4"
	"github.com/parnurzeal/gorequest"
)

// TestListSubscriptionEvents is an integration test that verifies the listing,
// of subscription events using the ChartMogul API.
func TestListSubscriptionEvents(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := NewRecorderWithAuthFilter("./fixtures/subscription_events")
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

	createSubscriptionEvents(api, ds.UUID, cus.ExternalID, t)
	validateListSubscriptionEvents(api, cus.ExternalID, t)
}

// validateListSubscriptionEvents validates that subscriptions can be correctly listed using the API.
func validateListSubscriptionEvents(api *cm.API, customerExternalID string, t *testing.T) {
	listSubscriptionEvents, err := api.ListSubscriptionEvents(&cm.FilterSubscriptionEvents{CustomerExternalID: customerExternalID}, &cm.Cursor{PerPage: 1})
	if err != nil {
		t.Fatalf("Failed to list subscriptions: %v", err)
	}

	if listSubscriptionEvents.Meta.HasMore != true {
		t.Fatalf("Expected HasMore to be true, got: %v", listSubscriptionEvents.Meta.HasMore)
	}
}

func createSubscriptionEvents(api *cm.API, dsUUID, cusExternalID string, t *testing.T) {
	testSubEvents := []*cm.SubscriptionEvent{
		{
			DataSourceUUID:         dsUUID,
			CustomerExternalID:     cusExternalID,
			SubscriptionExternalID: "TestSubEventsSubscription01",
			EventDate:              "2022-02-16",
			EffectiveDate:          "2022-02-28",
			EventType:              "subscription_start_scheduled",
			ExternalID:             "TestSubEventsSubEvent02",
			Currency:               "USD",
			AmountInCents:          1000,
			Quantity:               1,
			PlanExternalID:         "ext_plan",
		},
		{
			DataSourceUUID:         dsUUID,
			CustomerExternalID:     cusExternalID,
			SubscriptionExternalID: "TestSubEventsSubscription02",
			EventDate:              "2022-02-16",
			EffectiveDate:          "2022-02-28",
			EventType:              "subscription_cancelled",
			ExternalID:             "TestSubEventsSubEvent03",
			Currency:               "USD",
			PlanExternalID:         "ext_plan",
		},
	}

	for i, subEvent := range testSubEvents {
		result, err := api.CreateSubscriptionEvent(subEvent)
		if err != nil {
			t.Fatalf("Failed to create subscription event '%s': %v", subEvent.SubscriptionExternalID, err)
		}
		testSubEvents[i].ID = result.ID
	}
}
