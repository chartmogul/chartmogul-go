package chartmogul

import (
	"testing"

	"github.com/Sirupsen/logrus"
)

// TestImportDataSource tests creation, listing & deletion of Data Sources.
func TestImportCustomers(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}

	ds, err := api.ImportCreateDataSource("customers test")
	if err != nil {
		t.Error(err)
		return
	}
	defer api.ImportDeleteDataSource(ds.UUID)
	logrus.Debug("Data source created.")

	createdCustomer, err := api.ImportCreateCustomer(&Customer{
		ExternalID: "TestImportCustomers_01",
		Name:       "Test customer",
	}, ds.UUID)
	if err != nil {
		t.Error(err)
		return
	}
	logrus.Debug("Customer created.")

	listRes, err := api.ImportListCustomers(&ImportListCustomersParams{DataSourceUUID: ds.UUID})
	if err != nil {
		t.Error(err)
		return
	}
	found := false
	for _, c := range listRes.Customers {
		if c.UUID == createdCustomer.UUID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Customer not found in listing! %+v", listRes)
	}
	logrus.Debug("Customer found.")

	createdCustomer, err = api.ImportCreateCustomer(&Customer{
		ExternalID: "TestImportCustomers_01",
		Name:       "Test duplicate customer",
	}, ds.UUID)
	if err == nil {
		t.Error("No error on duplicate customer!")
	} else if createdCustomer.Errors.IsAlreadyExists() {
		logrus.Debug("Correct AlreadyExists.")
	} else {
		t.Errorf("Incorrect error: %v", createdCustomer.Errors)
	}
}
