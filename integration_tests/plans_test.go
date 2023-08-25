package integration

import (
	"net/http"
	"os"
	"reflect"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v4"
	"github.com/parnurzeal/gorequest"
)

// TestListPlansIntegration is an integration test that verifies the creation, listing,
// retrieval, and updating of plans using the ChartMogul API.
func TestListPlansIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := NewRecorderWithAuthFilter("./fixtures/plans")
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

	// Create and verify test plans.
	plans := createTestPlans(api, ds.UUID, t)

	validateListPlans(api, ds.UUID, t)
	validatePlanRetrievalAndUpdate(api, plans[1], t)
}

// createTestPlans creates a set of test plans for the integration test.
// It returns the created plans.
func createTestPlans(api *cm.API, dsUUID string, t *testing.T) []*cm.Plan {
	testPlans := []*cm.Plan{
		{
			DataSourceUUID: dsUUID,
			Name:           "Bronze Plan",
			IntervalCount:  1,
			IntervalUnit:   "month",
			ExternalID:     "plan_0001",
		},
		{
			DataSourceUUID: dsUUID,
			Name:           "Chocolate Plan",
			IntervalCount:  1,
			IntervalUnit:   "day",
			ExternalID:     "plan_0002",
		},
		{
			DataSourceUUID: dsUUID,
			Name:           "Flash Plan",
			IntervalCount:  1,
			IntervalUnit:   "year",
			ExternalID:     "plan_0003",
		},
	}

	for i, plan := range testPlans {
		respPlan, err := api.CreatePlan(plan)
		if err != nil {
			t.Fatalf("Failed to create plan '%s': %v", plan.Name, err)
		}
		testPlans[i].UUID = respPlan.UUID
	}

	return testPlans
}

// validateListPlans validates that the created plans can be correctly listed using the API.
func validateListPlans(api *cm.API, dsUUID string, t *testing.T) {
	listPlans, err := api.ListPlans(&cm.ListPlansParams{
		DataSourceUUID: dsUUID,
		Cursor:         cm.Cursor{PerPage: 2},
	})
	if err != nil {
		t.Fatalf("Failed to list plans: %v", err)
	}

	if listPlans.HasMore != true {
		t.Fatalf("Expected HasMore to be true, got: %v", listPlans.HasMore)
	}

	if len(listPlans.Plans) != 2 {
		t.Fatalf("Expected 2 plans, got: %d", len(listPlans.Plans))
	}
}

// validatePlanRetrievalAndUpdate checks that a given plan can be retrieved and updated
// correctly through the API.
func validatePlanRetrievalAndUpdate(api *cm.API, chocolatePlan *cm.Plan, t *testing.T) {
	retrievedPlan, err := api.RetrievePlan(chocolatePlan.UUID)
	if err != nil {
		t.Fatalf("Failed to retrieve plan: %v", err)
	}
	if !reflect.DeepEqual(chocolatePlan, retrievedPlan) {
		t.Fatalf("Expected plan %v, got %v", chocolatePlan, retrievedPlan)
	}

	updatePlanParams := &cm.Plan{
		Name: "Chocolate Plan 2",
	}
	updatedPlan, err := api.UpdatePlan(updatePlanParams, retrievedPlan.UUID)
	if err != nil {
		t.Fatalf("Failed to update plan: %v", err)
	}

	updatedRetrievedPlan, err := api.RetrievePlan(retrievedPlan.UUID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated plan: %v", err)
	}

	if !reflect.DeepEqual(updatedPlan, updatedRetrievedPlan) {
		t.Fatalf("Expected updated plan %v, got %v", updatedPlan, updatedRetrievedPlan)
	}
}
