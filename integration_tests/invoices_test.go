package integration

import (
	"net/http"
	"os"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v4"
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
	invoices := createTestInvoicesForcustomer(api, ds.UUID, *cus, plan.UUID, t)

	validateListInvoices(api, ds.UUID, cus.UUID, invoices, t)
}

// validateListinvoices validates that the created invoices can be correctly listed using the API.
func validateListInvoices(api *cm.API, dsUUID, cusUUID string, invoices *cm.Invoices, t *testing.T) {
	invoicesList, err := api.ListInvoices(&cm.Cursor{PerPage: 1}, cusUUID)
	if err != nil {
		t.Fatalf("Failed to list invoices: %v", err)
	}

	if invoicesList.HasMore != true {
		t.Fatalf("Expected HasMore to be true, got: %v", invoices.HasMore)
	}

	invoicesListAll, err := api.ListInvoices(&cm.Cursor{PerPage: 2}, cusUUID)
	if err != nil {
		t.Fatalf("Failed to list invoices: %v", err)
	}

	if invoicesListAll.HasMore != false {
		t.Fatalf("Expected HasMore to be false, got: %v", invoices.HasMore)
	}
}
