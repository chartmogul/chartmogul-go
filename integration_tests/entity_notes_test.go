package integration

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v5"
	"github.com/davecgh/go-spew/spew"
	"github.com/parnurzeal/gorequest"
)

func TestEntityNotesIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := NewRecorderWithAuthFilter("./fixtures/entity_notes")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() //nolint

	api := &cm.API{
		ApiKey: os.Getenv("CHARTMOGUL_API_KEY"),
		Client: &http.Client{Transport: r},
	}

	gorequest.DisableTransportSwap = true

	ds, err := api.CreateDataSource("Test Entity Notes")
	if err != nil {
		t.Fatal(err)
	}
	defer api.DeleteDataSource(ds.UUID) //nolint

	cus1, err := api.CreateCustomer(&cm.NewCustomer{
		Name:           "Test Entity Notes",
		ExternalID:     "ext_customer_1",
		DataSourceUUID: ds.UUID,
	})
	if err != nil {
		t.Fatal(err)
	}

	contact, err := api.CreateContact(&cm.NewContact{
		CustomerUUID:   cus1.UUID,
		DataSourceUUID: ds.UUID,
		FirstName:      "Adam",
		LastName:       "Smith",
		Email:          "adam@smith.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	customerNote, err := api.CreateCustomerEntityNote(&cm.NewEntityNote{
		Type: "note",
		Text: "This is a customer note",
	}, cus1.UUID)
	if err != nil {
		t.Fatal(err)
	}
	customerNotes, err := api.ListCustomerEntityNotes(&cm.ListEntityNotesParams{
		Cursor: cm.Cursor{PerPage: 10},
	}, cus1.UUID)
	if err != nil {
		t.Fatal(err)
	}
	expectedCustomerNotes := &cm.EntityNotes{Entries: []*cm.EntityNote{customerNote}}
	expectedCustomerNotes.Cursor = customerNotes.Cursor
	if !reflect.DeepEqual(customerNotes, expectedCustomerNotes) {
		spew.Dump(customerNotes)
		spew.Dump(expectedCustomerNotes)
		t.Fatal("Customer notes are not equal!")
	}

	contactNote, err := api.CreateContactEntityNote(&cm.NewEntityNote{
		Type:         "call",
		Text:         "Call with the contact",
		CallDuration: 60,
	}, contact.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if contactNote.CustomerUUID != "" {
		spew.Dump(contactNote)
		t.Fatal("Contact note should not carry a customer UUID!")
	}
	contactNotes, err := api.ListContactEntityNotes(&cm.ListEntityNotesParams{
		Cursor: cm.Cursor{PerPage: 10},
	}, contact.UUID)
	if err != nil {
		t.Fatal(err)
	}
	expectedContactNotes := &cm.EntityNotes{Entries: []*cm.EntityNote{contactNote}}
	expectedContactNotes.Cursor = contactNotes.Cursor
	if !reflect.DeepEqual(contactNotes, expectedContactNotes) {
		spew.Dump(contactNotes)
		spew.Dump(expectedContactNotes)
		t.Fatal("Contact notes are not equal!")
	}

	retrievedNote, err := api.RetrieveEntityNote(contactNote.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(retrievedNote, contactNote) {
		spew.Dump(retrievedNote)
		t.Fatal("Created note is not equal!")
	}

	updatedNote, err := api.UpdateEntityNote(&cm.UpdateEntityNote{
		Text: "This is an updated note",
	}, retrievedNote.UUID)
	if err != nil {
		t.Fatal(err)
	}
	updatedRetrievedNote, err := api.RetrieveEntityNote(updatedNote.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(updatedNote, updatedRetrievedNote) {
		spew.Dump(updatedRetrievedNote)
		t.Fatal("Updated note is not equal!")
	}

	if err := api.DeleteEntityNote(customerNote.UUID); err != nil {
		t.Fatal(err)
	}
}
