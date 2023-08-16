package integration

import (
	"log"
	"net/http"
	"os"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v3"
	"github.com/davecgh/go-spew/spew"
	"github.com/parnurzeal/gorequest"
)

func TestCreatePlanGroup(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := NewRecorderWithAuthFilter("./fixtures/create_plan_group")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() //nolint

	api := &cm.API{
		ApiKey: os.Getenv("CHARTMOGUL_API_KEY"),
		Client: &http.Client{Transport: r},
	}
	gorequest.DisableTransportSwap = true

	ds, err := api.CreateDataSource("Test Create Invoice")
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

	plg, err := api.CreatePlanGroup(&cm.PlanGroup{
		Name:  "My plan group",
		Plans: []*string{&plan.UUID},
	})

	if err != nil {
		t.Fatal(err)
	}

	if plg.PlansCount != 1 || plg.Name != "My plan group" {
		spew.Dump(plg)
		t.Fatal("the plan group was not created correctly")
	}

	plgUUUID := plg.UUID

	apiPlg, err := api.RetrievePlanGroup(plgUUUID)

	if err != nil {
		t.Fatal(err)
	}

	if apiPlg.UUID != plgUUUID {
		spew.Dump(apiPlg)
		spew.Dump(plgUUUID)
		t.Fatal("The uuids of the created Plan Group and the retreived Plan Group don't match")
	}

	plgPlans, err := api.ListPlanGroupPlans(&cm.PaginationWithCursor{PerPage: 200}, plgUUUID)

	if err != nil {
		t.Fatal(err)
	}

	if plgPlans.Plans[0].UUID != plan.UUID {
		spew.Dump(plgPlans)
		t.Fatal("The expected plan was not returned")
	}

	planGroups, err := api.ListPlanGroups(&cm.PaginationWithCursor{PerPage: 200})

	if err != nil {
		t.Fatal(err)
	}

	if planGroups.PlanGroups[0].UUID != apiPlg.UUID {
		spew.Dump(planGroups)
		t.Fatal("The expect plan groups were not returned")
	}
}
