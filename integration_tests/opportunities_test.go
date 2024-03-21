package integration

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v4"
	"github.com/davecgh/go-spew/spew"
	"github.com/parnurzeal/gorequest"
)

func TestOpportunitiesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := NewRecorderWithAuthFilter("./fixtures/opportunities")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() //nolint

	api := &cm.API{
		ApiKey: os.Getenv("CHARTMOGUL_API_KEY"),
		Client: &http.Client{Transport: r},
	}

	gorequest.DisableTransportSwap = true

	ds, err := api.CreateDataSource("Test Opportunities")
	if err != nil {
		t.Fatal(err)
	}
	defer api.DeleteDataSource(ds.UUID) //nolint

	cus1, err := api.CreateCustomer(&cm.NewCustomer{
		Name:           "Test Opportunities",
		ExternalID:     "ext_customer_1",
		DataSourceUUID: ds.UUID,
	})
	if err != nil {
		t.Fatal(err)
	}

	newOpportunityParams := &cm.NewOpportunity{
		CustomerUUID: cus1.UUID,
		Owner:              "kamil+pavlicko@chartmogul.com",
		Pipeline:           "New Business",
		PipelineStage:      "Discovery",
		EstimatedCloseDate: "2023-12-22",
		Currency:           "USD",
		AmountInCents:      100,
		Type:               "recurring",
		ForecastCategory:   "best_case",
		WinLikelihood:      3,
		Custom: []cm.Custom{
			{
				Key:   "from_campaign",
				Value:  true,
			},
		},
	}
	newOpportunity, err := api.CreateOpportunity(newOpportunityParams)
	if err != nil {
		t.Fatal(err)
	}
	allOpportunities, err := api.ListOpportunities(&cm.ListOpportunitiesParams{
		CustomerUUID: cus1.UUID,
		Cursor:       cm.Cursor{PerPage: 10},
	})
	if err != nil {
		t.Fatal(err)
	}

	var expectedAllOpportunities *cm.Opportunities = &cm.Opportunities{
		Entries: []*cm.Opportunity{newOpportunity},
	}
	expectedAllOpportunities.Cursor = allOpportunities.Cursor
	expectedAllOpportunities.HasMore = false

	if !reflect.DeepEqual(allOpportunities, expectedAllOpportunities) {
		spew.Dump(allOpportunities)
		spew.Dump(expectedAllOpportunities)
		t.Fatal("All opportunities are not equal!")
	}

	retrievedOpportunity, err := api.RetrieveOpportunity(newOpportunity.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(retrievedOpportunity, newOpportunity) {
		spew.Dump(retrievedOpportunity)
		t.Fatal("Created opportunity is not equal!")
	}

	updatedOpportunityParams := &cm.UpdateOpportunity{
		EstimatedCloseDate: "2024-12-22",
	}
	updatedOpportunity, err := api.UpdateOpportunity(updatedOpportunityParams, retrievedOpportunity.UUID)
	if err != nil {
		t.Fatal(err)
	}
	updatedRetrievedOpportunity, err := api.RetrieveOpportunity(updatedOpportunity.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(updatedOpportunity, updatedRetrievedOpportunity) {
		spew.Dump(updatedRetrievedOpportunity)
		t.Fatal("Updated opportunity is not equal!")
	}

	otherOpportunityParams := &cm.NewOpportunity{
		CustomerUUID: cus1.UUID,
		Owner:              "kamil+pavlicko@chartmogul.com",
		Pipeline:           "New Business",
		PipelineStage:      "Discovery",
		EstimatedCloseDate: "2023-12-22",
		Currency:           "EUR",
		AmountInCents:      1000,
		Type:               "recurring",
		ForecastCategory:   "best_case",
		WinLikelihood:      80,
		Custom: []cm.Custom{
			{
				Key:   "from_campaign",
				Value: false,
			},
		},
	}
	otherOpportunity, err := api.CreateOpportunity(otherOpportunityParams)
	if err != nil {
		t.Fatal(err)
	}

	deleteErr := api.DeleteOpportunity(otherOpportunity.UUID)
	if deleteErr != nil {
		t.Fatal(err)
	}
}
