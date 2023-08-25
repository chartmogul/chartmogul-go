package integration

import (
	"net/http"
	"os"
	"reflect"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v4"
	"github.com/parnurzeal/gorequest"
)

// TestListcustomersIntegration is an integration test that verifies the creation, listing,
// retrieval, and updating of customers using the ChartMogul API.
func TestListCustomersIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := NewRecorderWithAuthFilter("./fixtures/customers")
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

	// Create and verify test customers.
	customers := createTestCustomers(api, ds.UUID, t)

	validateListCustomers(api, ds.UUID, customers, t)
	validateSearchCustomers(api, ds.UUID, "abc@123.com", t)
	validateCustomerRetrievalAndUpdate(api, customers[1], t)
}

// createTestcustomers creates a set of test customers for the integration test.
// It returns the created customers.
func createTestCustomers(api *cm.API, dsUUID string, t *testing.T) []*cm.Customer {
	testCustomers := []*cm.NewCustomer{
		{
			DataSourceUUID: dsUUID,
			ExternalID:     "TestImportCustomers_01",
			Name:           "Test customer",
			Email:          "abc@123.com",
		},
		{
			DataSourceUUID: dsUUID,
			ExternalID:     "TestImportCustomers_02",
			Name:           "Test customer 2",
			Email:          "abc@123.com",
		},
		{
			DataSourceUUID: dsUUID,
			ExternalID:     "TestImportCustomers_03",
			Name:           "Test customer 3",
			Email:          "ghi@example.com",
		},
	}

	newCustomers := []*cm.Customer{}
	for _, cust := range testCustomers {
		respCust, err := api.CreateCustomer(cust)
		if err != nil {
			t.Fatalf("Failed to create customer '%s': %v", cust.Name, err)
		}
		newCustomers = append(newCustomers, respCust)
	}

	return newCustomers
}

// validateListcustomers validates that the created customers can be correctly listed using the API.
func validateListCustomers(api *cm.API, dsUUID string, customers []*cm.Customer, t *testing.T) {
	listCustomers, err := api.ListCustomers(&cm.ListCustomersParams{
		DataSourceUUID: dsUUID,
		Cursor:         cm.Cursor{PerPage: 2},
	})
	if err != nil {
		t.Fatalf("Failed to list customers: %v", err)
	}

	if listCustomers.HasMore != true {
		t.Fatalf("Expected HasMore to be true, got: %v", listCustomers.HasMore)
	}

	if len(listCustomers.Entries) != 2 {
		t.Fatalf("Expected 2 customers, got: %d", len(listCustomers.Entries))
	}
}

// validateSearchCustomers validates that the specified customers can be correctly searched using the API.
func validateSearchCustomers(api *cm.API, dsUUID string, email string, t *testing.T) {
	searchCustomers, err := api.SearchCustomers(&cm.SearchCustomersParams{
		Email:  email,
		Cursor: cm.Cursor{PerPage: 1},
	})
	if err != nil {
		t.Fatalf("Failed to find customers: %v", err)
	}

	if searchCustomers.HasMore != true {
		t.Fatalf("Expected HasMore to be true, got: %v", searchCustomers.HasMore)
	}
}

// validatecustomerRetrievalAndUpdate checks that a given customer can be retrieved and updated
// correctly through the API.
func validateCustomerRetrievalAndUpdate(api *cm.API, customerToUpdate *cm.Customer, t *testing.T) {
	retrievedCustomer, err := api.RetrieveCustomer(customerToUpdate.UUID)
	if err != nil {
		t.Fatalf("Failed to retrieve customer: %v", err)
	}
	if !reflect.DeepEqual(customerToUpdate, retrievedCustomer) {
		t.Fatalf("Expected customer %v, got %v", customerToUpdate, retrievedCustomer)
	}

	updatecustomerParams := &cm.Customer{
		Name: "Chocolate customer 2",
	}
	updatedcustomer, err := api.UpdateCustomer(updatecustomerParams, retrievedCustomer.UUID)
	if err != nil {
		t.Fatalf("Failed to update customer: %v", err)
	}

	updatedretrievedCustomer, err := api.RetrieveCustomer(retrievedCustomer.UUID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated customer: %v", err)
	}

	if !reflect.DeepEqual(updatedcustomer, updatedretrievedCustomer) {
		t.Fatalf("Expected updated customer %v, got %v", updatedcustomer, updatedretrievedCustomer)
	}
}
