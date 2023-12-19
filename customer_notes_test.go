package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestListNotes(t *testing.T) {
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
						"created_at": "2015-06-09T19:20:30Z",
						"updated_at": "2015-06-09T19:20:30Z"
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
	params := &ListNotesParams{Cursor: Cursor{PerPage: 1}, CustomerUUID: "cus_00000000-0000-0000-0000-000000000000"}
	customer_notes, err := tested.ListNotes(params)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(customer_notes.Entries) == 0 {
		spew.Dump(customer_notes)
		t.Fatal("Unexpected result")
	}
}

func TestRetrieveNote(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/customer_notes/note_00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"uuid": "note_00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"type": "note",
					"author": "John Doe (john@example.com)",
					"text": "This is a note",
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
	customer_note, err := tested.RetrieveNote("note_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if customer_note == nil {
		spew.Dump(customer_note)
		t.Fatal("Unexpected result")
	}
}

func TestCreateNote(t *testing.T) {
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
					"text": "This is a note",
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

	customer_notes, err := tested.CreateNote(&NewNote{
		CustomerUUID: "cus_00000000-0000-0000-0000-000000000000",
		Type:         "note",
		Text:         "This is a note",
		AuthorEmail:  "john@example.com",
		CallDuration: 0,
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if customer_notes.UUID != "note_00000000-0000-0000-0000-000000000000" {
		spew.Dump(customer_notes)
		t.Fatal("Unexpected result")
	}
}

func TestUpdateNote(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PATCH" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/customer_notes/note_00000000-0000-0000-0000-000000000000" {
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

	customer_note, err := tested.UpdateNote(&UpdateNote{
		Text: "This is a new note",
	}, "note_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if customer_note.Text != "This is a new note" {
		spew.Dump(customer_note)
		t.Fatal("Unexpected result")
	}
}

func TestDeleteNote(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/customer_notes/note_00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	err := tested.DeleteNote("note_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}
