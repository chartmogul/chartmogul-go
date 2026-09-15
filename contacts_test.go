package chartmogul

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestListContacts(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/contacts?per_page=3" {
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
	params := &ListContactsParams{Cursor: Cursor{PerPage: 3}}
	contacts, err := tested.ListContacts(params)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(contacts.Entries) == 0 {
		spew.Dump(contacts)
		t.Fatal("Unexpected result")
	}
}

func TestRetrieveContact(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/contacts/con_00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
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
					"last_seen": "2026-01-01T16:58:58.000Z",
					"custom": {
						"Facebook": "https://www.facebook.com/adam.smith/",
						"date_of_birth": "1985-01-22"
					}
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	contact, err := tested.RetrieveContact("con_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if contact == nil {
		spew.Dump(contact)
		t.Fatal("Unexpected result")
	}
	if contact.Email != "adam@smith.com" || contact.LastSeen != "2026-01-01T16:58:58.000Z" {
		spew.Dump(contact)
		t.Fatal("Expected email and last_seen to be decoded")
	}
}

func TestCreateContact(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/contacts" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{
					"uuid": "con_00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"customer_external_id": "customer_001",
					"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000",
					"external_id": "contact_external_id_001",
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

	externalID := "contact_external_id_001"
	contact, err := tested.CreateContact(&NewContact{
		CustomerUUID:   "cus_00000000-0000-0000-0000-000000000000",
		DataSourceUUID: "ds_00000000-0000-0000-0000-000000000000",
		ExternalID:     &externalID,
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
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if contact.UUID != "con_00000000-0000-0000-0000-000000000000" {
		spew.Dump(contact)
		t.Fatal("Unexpected result")
	}
	if contact.ExternalID == nil || *contact.ExternalID != "contact_external_id_001" {
		spew.Dump(contact)
		t.Fatal("Unexpected ExternalID")
	}
}

func TestCreateContactWithNullExternalID(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{
					"uuid": "con_00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000",
					"external_id": null
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	contact, err := tested.CreateContact(&NewContact{
		CustomerUUID:   "cus_00000000-0000-0000-0000-000000000000",
		DataSourceUUID: "ds_00000000-0000-0000-0000-000000000000",
		ExternalID:     nil,
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if contact.ExternalID != nil {
		spew.Dump(contact)
		t.Fatal("Expected ExternalID to be nil")
	}
}

// Note: Since ExternalID *string with omitempty serializes nil as an omitted field, both the
// TestCreateContactWithNullExternalID test case above and this test case are structurally identical
// from a serialization standpoint. They exist to document intent: (a) field can be absent and
// (b) field can be explicitly null. Both verify the response parses null correctly.
func TestCreateContactWithoutExternalID(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{
					"uuid": "con_00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000",
					"external_id": null
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	contact, err := tested.CreateContact(&NewContact{
		CustomerUUID:   "cus_00000000-0000-0000-0000-000000000000",
		DataSourceUUID: "ds_00000000-0000-0000-0000-000000000000",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if contact.ExternalID != nil {
		spew.Dump(contact)
		t.Fatal("Expected ExternalID to be nil")
	}
}

func TestUpdateContact(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PATCH" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/contacts/con_00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"uuid": "con_00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"customer_external_id": "customer_001",
					"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000",
					"external_id": "contact_external_id_002",
					"position": 10,
					"first_name": "Bill",
					"last_name": "Thompson",
					"title": "CTO",
					"email": "bill@example.com",
					"phone": "+987654321",
					"linked_in": "https://linkedin.com/bill-linkedin",
					"twitter": "https://twitter.com/bill-twitter",
					"notes": "New Heading\nBody\nFooter",
					"custom": {
						"Facebook": "https://www.facebook.com/bill.thompson",
						"date_of_birth": "1990-01-01"
					}
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	externalID := "contact_external_id_002"
	contact, err := tested.UpdateContact(&UpdateContact{
		ExternalID: &externalID,
		FirstName:  "Bill",
		LastName:   "Thompson",
		LinkedIn:   "https://linkedin.com/bill-linkedin",
		Notes:      "New Heading\nNew Body\nNew Footer",
		Phone:      "+987654321",
		Position:   10,
		Title:      "CTO",
		Twitter:    "https://twitter.com/bill-twitter",
		Custom: []Custom{
			{
				Key:   "Facebook",
				Value: "https://www.facebook.com/bill.thompson",
			},
			{
				Key:   "date_of_birth",
				Value: "1990-01-01",
			},
		},
	}, "con_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if contact.FirstName != "Bill" {
		spew.Dump(contact)
		t.Fatal("Unexpected result")
	}
	if contact.ExternalID == nil || *contact.ExternalID != "contact_external_id_002" {
		spew.Dump(contact)
		t.Fatal("Unexpected ExternalID")
	}
}

func TestUpdateContactWithNullExternalID(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PATCH" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"uuid": "con_00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000",
					"external_id": null
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	contact, err := tested.UpdateContact(&UpdateContact{
		ExternalID: nil,
	}, "con_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if contact.ExternalID != nil {
		spew.Dump(contact)
		t.Fatal("Expected ExternalID to be nil")
	}
}

func TestDeleteContact(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/contacts/con_00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	err := tested.DeleteContact("con_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestMergeContactParams(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/contacts/con_00000000-0000-0000-0000-000000000000/merge/con_00000000-0000-0000-0000-000000000001" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)

				w.Write([]byte(`{
					"uuid": "con_00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"customer_external_id": "customer_001",
					"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000",
					"position": 10,
					"first_name": "Bill",
					"last_name": "Thompson",
					"title": "CTO",
					"email": "bill@example.com",
					"phone": "+987654321",
					"linked_in": "https://linkedin.com/bill-linkedin",
					"twitter": "https://twitter.com/bill-twitter",
					"notes": "New Heading\nBody\nFooter",
					"custom": {
						"Facebook": "https://www.facebook.com/bill.thompson",
						"date_of_birth": "1990-01-01"
					}
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	contact, err := tested.MergeContacts("con_00000000-0000-0000-0000-000000000000", "con_00000000-0000-0000-0000-000000000001")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if contact == nil {
		spew.Dump(contact)
		t.Fatal("Unexpected result")
	}
}

func TestListContactsWithFilters(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				query := r.URL.Query()
				if query.Get("email") != "adam@smith.com" ||
					query.Get("customer_external_id") != "customer_001" ||
					query.Get("external_id") != "cont_001" ||
					query.Get("per_page") != "1" {
					t.Errorf("Unexpected query %v", r.URL.RawQuery)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{"entries": [], "has_more": false, "cursor": ""}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	params := &ListContactsParams{
		Cursor:             Cursor{PerPage: 1},
		Email:              "adam@smith.com",
		CustomerExternalID: "customer_001",
		ExternalID:         "cont_001",
	}
	_, err := tested.ListContacts(params)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestCreateContactWithoutCustomer(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				raw, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				var body map[string]interface{}
				if err := json.Unmarshal(raw, &body); err != nil {
					t.Fatal(err)
				}
				if _, ok := body["customer_uuid"]; ok {
					t.Errorf("Unexpected customer_uuid in body %s", raw)
				}
				if _, ok := body["data_source_uuid"]; ok {
					t.Errorf("Unexpected data_source_uuid in body %s", raw)
				}
				if body["email"] != "adam@smith.com" || body["last_active_at"] != "2026-01-01T16:58:58Z" {
					t.Errorf("Unexpected body %s", raw)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{
					"uuid": "con_00000000-0000-0000-0000-000000000000",
					"customer_uuid": null,
					"customer_external_id": null,
					"data_source_uuid": null,
					"first_name": "Adam",
					"last_name": "Smith",
					"email": "adam@smith.com",
					"last_seen": "2026-01-01T16:58:58.000Z"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	contact, err := tested.CreateContact(&NewContact{
		FirstName:    "Adam",
		LastName:     "Smith",
		Email:        "adam@smith.com",
		LastActiveAt: "2026-01-01T16:58:58Z",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if contact.CustomerUUID != "" || contact.Email != "adam@smith.com" {
		spew.Dump(contact)
		t.Fatal("Unexpected result")
	}
}

func TestListContactTasks(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/tasks?contact_uuid=con_00000000-0000-0000-0000-000000000000&per_page=1" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"entries": [{
						"task_uuid": "00000000-0000-0000-0000-000000000000",
						"customer_uuid": null,
						"associated_object": "contact",
						"associated_object_uuid": "con_00000000-0000-0000-0000-000000000000",
						"assignee": "keith+test1@chartmogul.com",
						"task_details": "This is some task details text.",
						"due_date": "2025-04-30T00:00:00Z",
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
	params := &ListTasksParams{Cursor: Cursor{PerPage: 1}}
	tasks, err := tested.ListContactTasks(params, "con_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(tasks.Entries) != 1 {
		spew.Dump(tasks)
		t.Fatal("Unexpected result")
	}
}

func TestListContactTasksWithNilParams(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.RequestURI != "/v/tasks?contact_uuid=con_00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{"entries": [], "has_more": false, "cursor": ""}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	_, err := tested.ListContactTasks(nil, "con_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestCreateContactTask(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/tasks" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				body := decodeNewTask(t, r)
				expected := AssociatedObjectIdentifier{
					AssociatedObject: "contact",
					Method:           "uuid",
					Value:            "con_00000000-0000-0000-0000-000000000000",
				}
				if body.CustomerUUID != "" || body.AssociatedObjectIdentifier == nil || *body.AssociatedObjectIdentifier != expected {
					t.Errorf("Unexpected body %+v", body)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{
					"task_uuid": "00000000-0000-0000-0000-000000000000",
					"customer_uuid": null,
					"associated_object": "contact",
					"associated_object_uuid": "con_00000000-0000-0000-0000-000000000000",
					"assignee": "keith+test1@chartmogul.com",
					"task_details": "This is some task details text.",
					"due_date": "2025-04-30T00:00:00Z",
					"created_at": "2025-04-01T12:00:00.000Z",
					"updated_at": "2025-04-01T12:00:00.000Z"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	task, err := tested.CreateContactTask(&NewTask{
		Assignee:    "keith+test1@chartmogul.com",
		TaskDetails: "This is some task details text.",
		DueDate:     "2025-04-30T00:00:00Z",
	}, "con_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if task.UUID != "00000000-0000-0000-0000-000000000000" {
		spew.Dump(task)
		t.Fatal("Unexpected result")
	}
}

func TestCreateContactTaskKeepsExplicitCustomer(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				body := decodeNewTask(t, r)
				if body.CustomerUUID != "cus_00000000-0000-0000-0000-000000000000" || body.AssociatedObjectIdentifier != nil {
					t.Errorf("Unexpected body %+v", body)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{"task_uuid": "00000000-0000-0000-0000-000000000000"}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	_, err := tested.CreateContactTask(&NewTask{
		CustomerUUID: "cus_00000000-0000-0000-0000-000000000000",
		Assignee:     "keith+test1@chartmogul.com",
		TaskDetails:  "This is some task details text.",
		DueDate:      "2025-04-30T00:00:00Z",
	}, "con_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestListContactEntityNotes(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/notes?contact_uuid=con_00000000-0000-0000-0000-000000000000&per_page=1" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"entries": [` + contactEntityNoteExample + `],
					"has_more": false,
					"cursor": "88abf99"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	params := &ListEntityNotesParams{Cursor: Cursor{PerPage: 1}}
	notes, err := tested.ListContactEntityNotes(params, "con_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(notes.Entries) != 1 {
		spew.Dump(notes)
		t.Fatal("Unexpected result")
	}
}

func TestCreateContactEntityNote(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/notes" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				body := decodeNewEntityNote(t, r)
				expected := AssociatedObjectIdentifier{
					AssociatedObject: "contact",
					Method:           "uuid",
					Value:            "con_00000000-0000-0000-0000-000000000000",
				}
				if body.CustomerUUID != "" || body.AssociatedObjectIdentifier == nil || *body.AssociatedObjectIdentifier != expected {
					t.Errorf("Unexpected body %+v", body)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(contactEntityNoteExample))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	note, err := tested.CreateContactEntityNote(&NewEntityNote{
		Type:         "call",
		Text:         "Call with the contact",
		CallDuration: 60,
	}, "con_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if note.UUID != "note_11111111-1111-1111-1111-111111111111" {
		spew.Dump(note)
		t.Fatal("Unexpected result")
	}
}
