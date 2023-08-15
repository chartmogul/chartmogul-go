package integration

import (
	"log"
	"net/http"
	"reflect"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v3"
	"github.com/davecgh/go-spew/spew"
	"github.com/parnurzeal/gorequest"

	"github.com/dnaeon/go-vcr/recorder"
)

func TestPlansIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := recorder.New("./fixtures/plans")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() //nolint

	api := &cm.API{}
	api.SetClient(&http.Client{Transport: r})
	gorequest.DisableTransportSwap = true

	ds, err := api.CreateDataSource("1")
	if err != nil {
		t.Fatal(err)
	}
	defer api.DeleteDataSource(ds.UUID) //nolint

	plan1, err := api.CreatePlan(&cm.Plan{
		DataSourceUUID: ds.UUID,
		Name:           "Bronze Plan",
		IntervalCount:  1,
		IntervalUnit:   "month",
		ExternalID:     "plan_0001",
	})
	if err != nil {
		t.Fatal(err)
	}

	listPlans, err := api.ListPlans(&cm.ListPlansParams{
		DataSourceUUID:       ds.UUID,
		PaginationWithCursor: cm.PaginationWithCursor{PerPage: 10},
	})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(plan1, listPlans.Plans[0]) {
		spew.Dump(plan1)
		spew.Dump(listPlans.Plans[0])
		t.Fatal("Plan is not equal!")
	}

	if listPlans.HasMore != false {
		spew.Dump(listPlans.HasMore)
		t.Fatal("HasMore is invalid")
	}
}
