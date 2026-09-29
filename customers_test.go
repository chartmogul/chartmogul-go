package chartmogul

import (
	"encoding/json"
	"io"
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

				body, err := io.ReadAll(r.Body)
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
	params := &ListNotesParams{Cursor: Cursor{PerPage: 1}}
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

func TestListCustomerOpportunities(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/opportunities?customer_uuid=cus_00000000-0000-0000-0000-000000000000&per_page=1" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"entries": [{
						"uuid": "00000000-0000-0000-0000-000000000000",
						"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
						"owner": "test1@example.org",
						"pipeline": "New business 1",
						"pipeline_stage": "Discovery",
						"estimated_close_date": "2023-12-22",
						"currency": "USD",
						"amount_in_cents": 100,
						"type": "recurring",
						"forecast_category": "pipeline",
						"win_likelihood": 3,
						"custom": {"from_campaign": true},
						"created_at": "2024-03-13T07:33:28.356Z",
						"updated_at": "2024-03-13T07:33:28.356Z"
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
	params := &ListOpportunitiesParams{Cursor: Cursor{PerPage: 1}}
	opportunities, err := tested.ListCustomerOpporunities(params, uuid)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(opportunities.Entries) == 0 {
		spew.Dump(opportunities)
		t.Fatal("Unexpected result")
	}
}

func TestCreateCustomerOpportunity(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/opportunities" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{
					"uuid": "00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"owner": "test1@example.org",
					"pipeline": "New business 1",
					"pipeline_stage": "Discovery",
					"estimated_close_date": "2023-12-22",
					"currency": "USD",
					"amount_in_cents": 100,
					"type": "recurring",
					"forecast_category": "pipeline",
					"win_likelihood": 3,
					"custom": {"from_campaign": true},
					"created_at": "2024-03-13T07:33:28.356Z",
					"updated_at": "2024-03-13T07:33:28.356Z"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	opportunity, err := tested.CreateCustomerOpportunity(&NewOpportunity{
		Owner:              "test1@example.org",
		Pipeline:           "New business 1",
		PipelineStage:      "Discovery",
		EstimatedCloseDate: "2023-12-22",
		Currency:           "USD",
		AmountInCents:      100,
		Type:               "recurring",
		ForecastCategory:   "pipeline",
		WinLikelihood:      3,
		Custom: []Custom{
			{
				Key:   "from_campaign",
				Value: true,
			},
		},
	}, "cus_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if opportunity.UUID != "00000000-0000-0000-0000-000000000000" {
		spew.Dump(opportunity)
		t.Fatal("Unexpected result")
	}
}

func TestListCustomerTasks(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/tasks?customer_uuid=cus_00000000-0000-0000-0000-000000000000&per_page=1" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"entries": [{
						"task_uuid": "00000000-0000-0000-0000-000000000000",
						"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
						"assignee": "keith+test1@chartmogul.com",
						"task_details": "This is some task details text.",
						"due_date": "2025-04-30T00:00:00Z",
						"completed_at": "2025-04-20T00:00:00Z",
						"created_at": "2025-04-01T12:00:00.000Z",
						"updated_at": "2025-04-01T12:00:00.000Z"
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
	params := &ListTasksParams{Cursor: Cursor{PerPage: 1}}
	tasks, err := tested.ListCustomerTasks(params, uuid)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(tasks.Entries) == 0 {
		spew.Dump(tasks)
		t.Fatal("Unexpected result")
	}
}

func TestCreateCustomerTask(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/tasks" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{
					"task_uuid": "00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"assignee": "keith+test1@chartmogul.com",
					"task_details": "This is some task details text.",
					"due_date": "2025-04-30T00:00:00Z",
					"completed_at": "2025-04-20T00:00:00Z",
					"created_at": "2025-04-01T12:00:00.000Z",
					"updated_at": "2025-04-01T12:00:00.000Z"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	task, err := tested.CreateCustomerTask(&NewTask{
		Assignee:    "keith+test1@chartmogul.com",
		TaskDetails: "This is some task details text.",
		DueDate:     "2025-04-30T00:00:00Z",
		CompletedAt: "2025-04-20T00:00:00Z",
	}, "cus_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if task.UUID != "00000000-0000-0000-0000-000000000000" {
		spew.Dump(task)
		t.Fatal("Unexpected result")
	}
}

func TestRetrieveCustomerWithOptions(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				query := r.URL.Query()
				if query.Get("with_overrides") != "true" {
					t.Errorf("Expected with_overrides=true, got: %v", query.Get("with_overrides"))
				}
				if query.Get("attributes_with_history") != "company,custom.channel" {
					t.Errorf("Expected attributes_with_history=company,custom.channel, got: %v", query.Get("attributes_with_history"))
				}
				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write([]byte(`{
					"uuid": "cus_00000000-0000-0000-0000-000000000000",
					"company": "Pinned Co",
					"overrides": {
						"company": true,
						"attributes": {"custom": {"channel": true}}
					},
					"historical_values": {
						"company": [
							{"value": "Old Co", "update_performed_at": null, "update_performed_by": null, "initial": true},
							{"value": "Pinned Co", "update_performed_at": "2026-09-25T10:00:00Z", "update_performed_by": "API", "initial": false}
						],
						"attributes": {"custom": {"channel": [
							{"value": "Facebook", "update_performed_at": "2026-09-25T10:00:00Z", "update_performed_by": "API", "initial": false}
						]}}
					}
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	trueBool := true
	customer, err := tested.RetrieveCustomerWithOptions("cus_00000000-0000-0000-0000-000000000000", &RetrieveCustomerParams{
		WithOverrides:         &trueBool,
		AttributesWithHistory: "company,custom.channel",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	expectedOverrides := map[string]interface{}{
		"company":    true,
		"attributes": map[string]interface{}{"custom": map[string]interface{}{"channel": true}},
	}
	if !reflect.DeepEqual(customer.Overrides, expectedOverrides) {
		spew.Dump(customer.Overrides)
		t.Fatal("Unexpected overrides")
	}
	if customer.HistoricalValues == nil {
		spew.Dump(customer)
		t.Fatal("Expected historical_values")
	}
}

func TestRetrieveCustomerWithOptionsNil(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.RawQuery != "" {
					t.Errorf("Expected no query, got: %v", r.URL.RawQuery)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"uuid": "cus_00000000-0000-0000-0000-000000000000"}`)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	customer, err := tested.RetrieveCustomerWithOptions("cus_00000000-0000-0000-0000-000000000000", nil)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if customer.Overrides != nil {
		spew.Dump(customer)
		t.Fatal("Expected no overrides")
	}
}

func TestCreateCustomerWithOverrides(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					return
				}
				var incoming map[string]interface{}
				if err := json.Unmarshal(body, &incoming); err != nil {
					t.Error(err)
					return
				}
				expectedOverrides := map[string]interface{}{
					"company":    true,
					"attributes": map[string]interface{}{"custom": map[string]interface{}{"channel": true}},
				}
				if !reflect.DeepEqual(incoming["overrides"], expectedOverrides) {
					spew.Dump(incoming["overrides"])
					t.Error("Request overrides don't equal expected value")
				}
				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write([]byte(`{
					"uuid": "cus_00000000-0000-0000-0000-000000000000",
					"company": "Pinata Technologies",
					"overrides": {
						"company": true,
						"attributes": {"custom": {"channel": true}}
					}
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	customer, err := tested.CreateCustomer(&NewCustomer{
		DataSourceUUID: "ds_00000000-0000-0000-0000-000000000000",
		ExternalID:     "cus_0001",
		Company:        "Pinata Technologies",
		Overrides: map[string]interface{}{
			"company":    true,
			"attributes": map[string]interface{}{"custom": map[string]interface{}{"channel": true}},
		},
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if customer.Overrides == nil {
		spew.Dump(customer)
		t.Fatal("Expected overrides in response")
	}
}

func TestUpdateCustomerDoesNotReplayOverrides(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Method == "GET" {
					//nolint
					w.Write([]byte(`{
						"uuid": "cus_00000000-0000-0000-0000-000000000000",
						"company": "Pinned Co",
						"attributes": {
							"custom": {"channel": "Facebook"},
							"overrides": {"custom": {"channel": true}},
							"historical_values": {"custom": {"channel": []}}
						},
						"overrides": {"company": true},
						"historical_values": {"company": []}
					}`))
					return
				}
				if r.Method != "PATCH" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					return
				}
				var incoming map[string]interface{}
				if err := json.Unmarshal(body, &incoming); err != nil {
					t.Error(err)
					return
				}
				for _, key := range []string{"overrides", "historical_values"} {
					if _, ok := incoming[key]; ok {
						t.Errorf("PATCH body must not replay %q from the retrieved customer", key)
					}
					if attrs, ok := incoming["attributes"].(map[string]interface{}); ok {
						if _, ok := attrs[key]; ok {
							t.Errorf("PATCH body attributes must not replay %q from the retrieved attributes", key)
						}
					}
				}
				w.Write([]byte(`{}`)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	trueBool := true
	retrieved, err := tested.RetrieveCustomerWithOptions("cus_00000000-0000-0000-0000-000000000000", &RetrieveCustomerParams{
		WithOverrides:         &trueBool,
		AttributesWithHistory: "company,custom.channel",
	})
	if err != nil {
		spew.Dump(err)
		t.Fatal("Retrieve not expected to fail")
	}
	if retrieved.Overrides == nil || retrieved.HistoricalValues == nil {
		spew.Dump(retrieved)
		t.Fatal("Expected retrieved customer to carry overrides and historical_values")
	}
	if retrieved.Attributes == nil || retrieved.Attributes.Overrides == nil {
		spew.Dump(retrieved)
		t.Fatal("Expected retrieved attributes to carry overrides")
	}

	if _, err := tested.UpdateCustomer(retrieved, retrieved.UUID); err != nil {
		spew.Dump(err)
		t.Fatal("Update not expected to fail")
	}
}

func TestUpdateCustomerV2WithOverrides(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PATCH" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					return
				}
				var incoming map[string]interface{}
				if err := json.Unmarshal(body, &incoming); err != nil {
					t.Error(err)
					return
				}
				expectedOverrides := map[string]interface{}{
					"company":    true,
					"attributes": map[string]interface{}{"custom": map[string]interface{}{"channel": true}},
				}
				if !reflect.DeepEqual(incoming["overrides"], expectedOverrides) {
					spew.Dump(incoming["overrides"])
					t.Error("Request overrides don't equal expected value")
				}
				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write([]byte(`{
					"uuid": "cus_00000000-0000-0000-0000-000000000000",
					"company": "Pinned Co",
					"overrides": {"company": true, "attributes": {"custom": {"channel": true}}}
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	company := "Pinned Co"
	customer, err := tested.UpdateCustomerV2(&UpdateCustomer{
		Company: &company,
		Overrides: map[string]interface{}{
			"company":    true,
			"attributes": map[string]interface{}{"custom": map[string]interface{}{"channel": true}},
		},
	}, "cus_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if customer.Overrides == nil {
		spew.Dump(customer)
		t.Fatal("Expected overrides in response")
	}
}
