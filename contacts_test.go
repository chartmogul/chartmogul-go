package chartmogul

import (
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
