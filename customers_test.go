package chartmogul

import (
	"testing"

	"github.com/Sirupsen/logrus"
)

// TestImportCustomers tests creation & listing of customers + duplicate error.
func TestImportCustomers(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}

	ds, err := api.CreateDataSource("customers test")
	if err != nil {
		t.Error(err)
		return
	}
	defer api.DeleteDataSource(ds.UUID)
	logrus.Debug("Data source created.")

	createdCustomer, err := api.CreateCustomer(&NewCustomer{
		DataSourceUUID: ds.UUID,
		ExternalID:     "TestImportCustomers_01",
		Name:           "Test customer",
	})
	toBeDeletedUUID := createdCustomer.UUID
	if err != nil {
		t.Error(err)
		return
	}
	logrus.Debug("Customer created.")

	listRes, err := api.ListCustomers(&ListCustomersParams{DataSourceUUID: ds.UUID})
	if err != nil {
		t.Error(err)
		return
	}
	found := false
	for _, c := range listRes.Entries {
		if c.UUID == createdCustomer.UUID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Customer not found in listing! %+v", listRes)
	}
	logrus.Debug("Customer found.")

	createdCustomer, err = api.CreateCustomer(&NewCustomer{
		DataSourceUUID: ds.UUID,
		ExternalID:     "TestImportCustomers_01",
		Name:           "Test duplicate customer",
	})
	if err == nil {
		t.Error("No error on duplicate customer!")
	} else if createdCustomer.Errors.IsAlreadyExists() {
		logrus.Debug("Correct AlreadyExists.")
	} else {
		t.Errorf("Incorrect error: %v", createdCustomer.Errors)
	}

	err = api.DeleteCustomer(toBeDeletedUUID)
	if err != nil {
		t.Errorf("Couldn't delete customer: %v", err)
	}
}
