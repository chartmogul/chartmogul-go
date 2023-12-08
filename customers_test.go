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
	defer api.DeleteDataSource(ds.UUID) //nolint
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
				w.Write([]byte("{}")) //nolint

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
		ApiKey: "token",
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
		ApiKey: "token",
	}
	err := tested.DeleteCustomerInvoices("dataSourceUUID", "customerUUID")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestPurgeCustomerV2WithCustomerExternalID(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/data_sources/dataSourceUUID/customers/customerUUID/invoices?customer_external_id=externalID" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	err := tested.DeleteCustomerInvoicesV2("dataSourceUUID", "customerUUID", &DeleteCustomerInvoicesParams{
		CustomerExternalID: "externalID",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestPurgeCustomerV2WithoutCustomerExternalID(t *testing.T) {
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
		ApiKey: "token",
	}

	err := tested.DeleteCustomerInvoicesV2("dataSourceUUID", "customerUUID", nil)

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
				//nolint
				w.Write([]byte(`{"entries": [],"current_page": 1,"total_pages": 1,
					"has_more": false,"per_page": 200,"page": 1}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	customers, err := tested.ListCustomers(nil)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(customers.Entries) != 0 {
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
				//nolint
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
		ApiKey: "token",
	}
	customers, err := tested.ListCustomers(&ListCustomersParams{System: "whatnot"})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(customers.Entries) != 0 {
		spew.Dump(customers)
		t.Fatal("Unexpected result")
	}
}

func TestListCustomersContacts(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/customers/cus_00000000-0000-0000-0000-000000000000/contacts?per_page=3" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"entries": [{
						"uuid": "con_00000000-0000-0000-0000-000000000000",
						"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
						"customer_external_id": "123",
						"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000",
						"position": 1,
						"first_name": "Adam",
						"last_name": "Smith",
						"title": "CEO",
						"email": "adam@smith.com",
						"phone": "Lead",
						"linked_in": null,
						"twitter": null,
						"notes": null,
						"custom": {
							"Facebook": "https://www.facebook.com/adam.smith/",
							"date_of_birth": "1985-01-22"
						}
					}],
					"has_more": false,
					"cursor": "88abf99"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	uuid := "cus_00000000-0000-0000-0000-000000000000"
	params := &ListContactsParams{Cursor: Cursor{PerPage: 3}}
	contacts, err := tested.ListCustomersContacts(params, uuid)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(contacts.Entries) == 0 {
		spew.Dump(contacts)
		t.Fatal("Unexpected result")
	}
}

func TestCreateCustomersContact(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/customers/cus_00000000-0000-0000-0000-000000000000/contacts" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{
					"uuid": "con_00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"customer_external_id": "customer_001",
					"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000",
					"position": 9,
					"first_name": "Adam",
					"last_name": "Smith",
					"title": "CEO",
					"email": "adam@example.com",
					"phone": "+1234567890",
					"linked_in": "https://linkedin.com/linkedin",
					"twitter": "https://twitter.com/twitter",
					"notes": "Heading\nBody\nFooter",
					"custom": {
						"Facebook": "https://www.facebook.com/adam.smith",
						"date_of_birth": "1985-01-22"
					}
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	contact, err := tested.CreateCustomersContact(&NewContact{
		DataSourceUUID: "ds_00000000-0000-0000-0000-000000000000",
		FirstName:      "Adam",
		LastName:       "Smith",
		LinkedIn:       "https://linkedin.com/linkedin",
		Notes:          "Heading\nBody\nFooter",
		Phone:          "+1234567890",
		Position:       1,
		Title:          "CEO",
		Twitter:        "https://twitter.com/twitter",
		Custom: []Custom{
			{
				Key:   "Facebook",
				Value: "https://www.facebook.com/adam.smith",
			},
			{
				Key:   "date_of_birth",
				Value: "1985-01-22",
			},
		},
	}, "cus_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if contact.UUID != "con_00000000-0000-0000-0000-000000000000" {
		spew.Dump(contact)
		t.Fatal("Unexpected result")
	}
}

func TestListCustomerNotes(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/customer_notes?customer_uuid=cus_00000000-0000-0000-0000-000000000000&per_page=1" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"entries": [{
						"uuid": "note_00000000-0000-0000-0000-000000000000",
						"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
						"type": "note",
						"author": "John Doe (john@example.com)",
						"text": "This is a note",
						"call_duration": 0,
						"created_at": "2017-06-09T13:14:00-04:00",
						"updated_at": "2017-06-09T13:14:00-04:00"
					}],
					"has_more": false,
					"cursor": "88abf99"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	uuid := "cus_00000000-0000-0000-0000-000000000000"
	params := &ListNotesParams{PaginationWithCursor: PaginationWithCursor{PerPage: 1}}
	customer_notes, err := tested.ListCustomerNotes(params, uuid)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(customer_notes.Entries) == 0 {
		spew.Dump(customer_notes)
		t.Fatal("Unexpected result")
	}
}

func TestCreateCustomerNote(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/customer_notes" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{
					"uuid": "note_00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"type": "note",
					"author": "John Doe (john@example.com)",
					"text": "This is a note",
					"call_duration": 0,
					"created_at": "2017-06-09T13:14:00-04:00",
					"updated_at": "2017-06-09T13:14:00-04:00"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	customer_note, err := tested.CreateCustomerNote(&NewNote{
		Type:        "note",
		AuthorEmail: "john@example.com",
		Text:        "This is a note",
	}, "cus_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if customer_note.UUID != "note_00000000-0000-0000-0000-000000000000" {
		spew.Dump(customer_note)
		t.Fatal("Unexpected result")
	}
}
