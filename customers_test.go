package chartmogul

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
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
	log.Println("Data source created.")

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
	log.Println("Customer created.")

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
	log.Println("Customer found.")

	createdCustomer, err = api.CreateCustomer(&NewCustomer{
		DataSourceUUID: ds.UUID,
		ExternalID:     "TestImportCustomers_01",
		Name:           "Test duplicate customer",
	})
	if err == nil {
		t.Error("No error on duplicate customer!")
	} else if createdCustomer.Errors.IsAlreadyExists() {
		log.Println("Correct AlreadyExists.")
	} else {
		t.Errorf("Incorrect error: %v", createdCustomer.Errors)
	}

	err = api.DeleteCustomer(toBeDeletedUUID)
	if err != nil {
		t.Errorf("Couldn't delete customer: %v", err)
	}
}

func TestFormattingOfSourceInCustomAttributeUpdate(t *testing.T) {

	expected := map[string]interface{}{
		"attributes": map[string]interface{}{
			"custom": map[string]interface{}{
				"some key": map[string]interface{}{
					"value":  "some value",
					"source": "some awesome integration",
				},
			},
			"stripe": map[string]interface{}{
				"some key": map[string]interface{}{
					"value":  "some value",
					"source": "some awesome integration",
				},
			},
		},
	}

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{}"))

				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					return
				}
				var incoming interface{}
				err = json.Unmarshal(body, &incoming)
				if err != nil {
					t.Error(err)
					return
				}
				if !reflect.DeepEqual(expected, incoming) {
					spew.Dump(expected, incoming)
					t.Error("Doesn't equal expected value")
					return
				}
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		AccountToken: "token",
		AccessKey:    "key",
	}
	_, err := tested.UpdateCustomer(&Customer{
		Attributes: &Attributes{
			Custom: map[string]interface{}{
				"some key": AttributeWithSource{
					Value:  "some value",
					Source: "some awesome integration",
				}},
			Stripe: map[string]interface{}{
				"some key": AttributeWithSource{
					Value:  "some value",
					Source: "some awesome integration",
				}},
		},
	}, "customerUUID")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestPurgeCustomer(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/data_sources/dataSourceUUID/customers/customerUUID/invoices" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		AccountToken: "token",
		AccessKey:    "key",
	}
	err := tested.DeleteCustomerInvoices("dataSourceUUID", "customerUUID")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestUpdateCustomerSerialization(t *testing.T) {
	empty := ""
	cus := &UpdateCustomer{
		Name: &empty,
	}
	output, err := json.Marshal(cus)
	if err != nil {
		t.Fatal("Not expected to fail")
	}

	result := string(output)
	if result != `{"name":""}` {
		t.Fatal("Not expected to fail")
	}
}

func TestNilListCustomers(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/customers" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"entries": [],"current_page": 1,"total_pages": 1,
					"has_more": false,"per_page": 200,"page": 1}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		AccountToken: "token",
		AccessKey:    "key",
	}
	customers, err := tested.ListCustomers(nil)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(customers.Entries) != 0 || customers.PerPage != 200 {
		spew.Dump(customers)
		t.Fatal("Unexpected result")
	}
}

func TestSystemListCustomers(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/customers?system=whatnot" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"entries": [],
												 "current_page": 1,
												 "total_pages": 1,
												 "has_more": false,
												 "per_page": 200,
												 "page": 1}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		AccountToken: "token",
		AccessKey:    "key",
	}
	customers, err := tested.ListCustomers(&ListCustomersParams{System: "whatnot"})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(customers.Entries) != 0 || customers.PerPage != 200 {
		spew.Dump(customers)
		t.Fatal("Unexpected result")
	}
}
