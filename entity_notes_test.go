package chartmogul

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

const customerEntityNoteExample = `{
	"uuid": "note_00000000-0000-0000-0000-000000000000",
	"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
	"associated_object": "customer",
	"associated_object_uuid": "cus_00000000-0000-0000-0000-000000000000",
	"type": "note",
	"author": "John Doe (john@example.com)",
	"text": "This is a note",
	"call_duration": 0,
	"created_at": "2015-06-09T19:20:30Z",
	"updated_at": "2015-06-09T19:20:30Z"
}`

const contactEntityNoteExample = `{
	"uuid": "note_11111111-1111-1111-1111-111111111111",
	"customer_uuid": null,
	"associated_object": "contact",
	"associated_object_uuid": "con_00000000-0000-0000-0000-000000000000",
	"type": "call",
	"author": "John Doe (john@example.com)",
	"text": "Call with the contact",
	"call_duration": 60,
	"created_at": "2015-06-09T19:20:30Z",
	"updated_at": "2015-06-09T19:20:30Z"
}`

func TestListEntityNotes(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/notes?per_page=1" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"entries": [` + customerEntityNoteExample + `],
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
	notes, err := tested.ListEntityNotes(params)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(notes.Entries) != 1 {
		spew.Dump(notes)
		t.Fatal("Unexpected result")
	}
	if notes.Entries[0].AssociatedObject != "customer" || notes.Entries[0].AssociatedObjectUUID != "cus_00000000-0000-0000-0000-000000000000" {
		spew.Dump(notes)
		t.Fatal("Expected associated object to be decoded")
	}
}

func TestListEntityNotesWithContactUuid(t *testing.T) {
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
	params := &ListEntityNotesParams{Cursor: Cursor{PerPage: 1}, ContactUUID: "con_00000000-0000-0000-0000-000000000000"}
	notes, err := tested.ListEntityNotes(params)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(notes.Entries) != 1 {
		spew.Dump(notes)
		t.Fatal("Unexpected result")
	}
	if notes.Entries[0].CustomerUUID != "" || notes.Entries[0].AssociatedObject != "contact" {
		spew.Dump(notes)
		t.Fatal("Expected a contact note without customer UUID")
	}
}

func TestListEntityNotesWithFilters(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				query := r.URL.Query()
				if query.Get("customer_uuid") != "cus_00000000-0000-0000-0000-000000000000" ||
					query.Get("type") != "call" ||
					query.Get("author_email") != "john@example.com" ||
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
	params := &ListEntityNotesParams{
		Cursor:       Cursor{PerPage: 1},
		CustomerUUID: "cus_00000000-0000-0000-0000-000000000000",
		Type:         "call",
		AuthorEmail:  "john@example.com",
	}
	_, err := tested.ListEntityNotes(params)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestRetrieveEntityNote(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/notes/note_00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(customerEntityNoteExample))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	note, err := tested.RetrieveEntityNote("note_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if note.UUID != "note_00000000-0000-0000-0000-000000000000" {
		spew.Dump(note)
		t.Fatal("Unexpected result")
	}
}

func TestCreateEntityNote(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/notes" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(customerEntityNoteExample))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	note, err := tested.CreateEntityNote(&NewEntityNote{
		CustomerUUID: "cus_00000000-0000-0000-0000-000000000000",
		Type:         "note",
		Text:         "This is a note",
		AuthorEmail:  "john@example.com",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if note.UUID != "note_00000000-0000-0000-0000-000000000000" {
		spew.Dump(note)
		t.Fatal("Unexpected result")
	}
}

func TestCreateEntityNoteWithAssociatedObjectIdentifier(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				body := decodeNewEntityNote(t, r)
				if body.CustomerUUID != "" {
					t.Errorf("Unexpected customer_uuid %v", body.CustomerUUID)
				}
				expected := AssociatedObjectIdentifier{
					AssociatedObject: "contact",
					Method:           "uuid",
					Value:            "con_00000000-0000-0000-0000-000000000000",
				}
				if body.AssociatedObjectIdentifier == nil || *body.AssociatedObjectIdentifier != expected {
					t.Errorf("Unexpected associated_object_identifier %v", body.AssociatedObjectIdentifier)
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

	note, err := tested.CreateEntityNote(&NewEntityNote{
		AssociatedObjectIdentifier: &AssociatedObjectIdentifier{
			AssociatedObject: AssociatedObjectContact,
			Method:           AssociatedObjectIdentifierMethodUUID,
			Value:            "con_00000000-0000-0000-0000-000000000000",
		},
		Type:         "call",
		Text:         "Call with the contact",
		CallDuration: 60,
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if note.AssociatedObjectUUID != "con_00000000-0000-0000-0000-000000000000" || note.CustomerUUID != "" {
		spew.Dump(note)
		t.Fatal("Unexpected result")
	}
}

func TestUpdateEntityNote(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PATCH" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/notes/note_00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"uuid": "note_00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"type": "note",
					"text": "This is a new note",
					"author": "John Doe (john@example.com)",
					"call_duration": 0,
					"created_at": "2015-06-09T19:20:30Z",
					"updated_at": "2015-06-09T19:20:30Z"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	note, err := tested.UpdateEntityNote(&UpdateEntityNote{
		Text: "This is a new note",
	}, "note_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if note.Text != "This is a new note" {
		spew.Dump(note)
		t.Fatal("Unexpected result")
	}
}

func TestDeleteEntityNote(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/notes/note_00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusAccepted)
				//nolint
				w.Write([]byte(`{"message": "Note deleted"}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	err := tested.DeleteEntityNote("note_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func decodeNewEntityNote(t *testing.T, r *http.Request) NewEntityNote {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	var body NewEntityNote
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}
	return body
}

func decodeNewTask(t *testing.T, r *http.Request) NewTask {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}
	var body NewTask
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}
	return body
}
