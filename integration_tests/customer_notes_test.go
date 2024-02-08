package integration

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v4"
	"github.com/davecgh/go-spew/spew"
	"github.com/parnurzeal/gorequest"
)

func TestCustomerNotesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := NewRecorderWithAuthFilter("./fixtures/notes")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() //nolint

	api := &cm.API{
		ApiKey: os.Getenv("CHARTMOGUL_API_KEY"),
		Client: &http.Client{Transport: r},
	}

	gorequest.DisableTransportSwap = true

	ds, err := api.CreateDataSource("Test Customer Notes")
	if err != nil {
		t.Fatal(err)
	}
	defer api.DeleteDataSource(ds.UUID) //nolint

	cus1, err := api.CreateCustomer(&cm.NewCustomer{
		Name:           "Test Customer Notes",
		ExternalID:     "ext_customer_1",
		DataSourceUUID: ds.UUID,
	})
	if err != nil {
		t.Fatal(err)
	}

	newNoteParams := &cm.NewNote{
		CustomerUUID: cus1.UUID,
		Type:         "note",
		Text:         "This is a note",
	}
	newNote, err := api.CreateNote(newNoteParams)
	if err != nil {
		t.Fatal(err)
	}
	allNotes, err := api.ListNotes(&cm.ListNotesParams{
		CustomerUUID: cus1.UUID,
		Cursor:       cm.Cursor{PerPage: 10},
	})
	if err != nil {
		t.Fatal(err)
	}

	var expectedAllNotes *cm.Notes = &cm.Notes{
		Entries: []*cm.Note{newNote},
	}
	expectedAllNotes.Cursor = allNotes.Cursor
	expectedAllNotes.HasMore = false

	if !reflect.DeepEqual(allNotes, expectedAllNotes) {
		spew.Dump(allNotes)
		spew.Dump(expectedAllNotes)
		t.Fatal("All notes are not equal!")
	}

	retrievedNote, err := api.RetrieveNote(newNote.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(retrievedNote, newNote) {
		spew.Dump(retrievedNote)
		t.Fatal("Created note is not equal!")
	}

	updatedNoteParams := &cm.UpdateNote{
		Text: "This is an updated note",
	}
	updatedNote, err := api.UpdateNote(updatedNoteParams, retrievedNote.UUID)
	if err != nil {
		t.Fatal(err)
	}
	updatedRetrievedNote, err := api.RetrieveNote(updatedNote.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(updatedNote, updatedRetrievedNote) {
		spew.Dump(updatedRetrievedNote)
		t.Fatal("Updated note is not equal!")
	}

	otherNoteParams := &cm.NewNote{
		CustomerUUID: cus1.UUID,
		Type:         "note",
		Text:         "This is another note",
	}
	otherNote, err := api.CreateNote(otherNoteParams)
	if err != nil {
		t.Fatal(err)
	}

	deleteErr := api.DeleteNote(otherNote.UUID)
	if deleteErr != nil {
		t.Fatal(err)
	}
}
