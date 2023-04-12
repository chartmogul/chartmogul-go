package integration

import (
	"log"
	"net/http"
	"reflect"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v3"
	"github.com/davecgh/go-spew/spew"
	"github.com/parnurzeal/gorequest"

	"github.com/dnaeon/go-vcr/recorder"
)

func TestContactsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := recorder.New("./fixtures/contacts")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() //nolint

	api := &cm.API{}
	api.SetClient(&http.Client{Transport: r})
	gorequest.DisableTransportSwap = true

	ds, err := api.CreateDataSource("Test Contact")
	if err != nil {
		t.Fatal(err)
	}
	defer api.DeleteDataSource(ds.UUID) //nolint

	cus1, err := api.CreateCustomer(&cm.NewCustomer{
		Name:           "Test Contact",
		Email:          "briwa@chartmogul.com",
		ExternalID:     "ext_customer_1",
		DataSourceUUID: ds.UUID,
	})
	if err != nil {
		t.Fatal(err)
	}

	newContactParams := &cm.NewContact{
		CustomerUUID:   cus1.UUID,
		DataSourceUUID: ds.UUID,
		FirstName:      "Adam",
		LastName:       "Smith",
		LinkedIn:       "https://linkedin.com/linkedin",
		Notes:          "Heading\nBody\nFooter",
		Phone:          "+1234567890",
		Position:       1,
		Title:          "CEO",
		Twitter:        "https://twitter.com/twitter",
		Custom: []cm.Custom{
			{
				Key:   "Facebook",
				Value: "https://www.facebook.com/adam.smith",
			},
			{
				Key:   "date_of_birth",
				Value: "1985-01-22",
			},
		},
	}
	newContact, err := api.CreateContact(newContactParams)
	if err != nil {
		t.Fatal(err)
	}
	retrievedContact, err := api.RetrieveContact(newContact.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(retrievedContact, newContact) {
		spew.Dump(retrievedContact)
		t.Fatal("Created contact is not equal!")
	}

	updatedContactParams := &cm.UpdateContact{
		FirstName: "Bill",
		LastName:  "Thompson",
		LinkedIn:  "https://linkedin.com/bill-linkedin",
		Notes:     "New Heading\nNew Body\nNew Footer",
		Phone:     "+987654321",
		Position:  10,
		Title:     "CTO",
		Twitter:   "https://twitter.com/bill-twitter",
		Custom: []cm.Custom{
			{
				Key:   "Facebook",
				Value: "https://www.facebook.com/bill.thompson",
			},
			{
				Key:   "date_of_birth",
				Value: "1990-01-01",
			},
		},
	}
	updatedContact, err := api.UpdateContact(updatedContactParams, retrievedContact.UUID)
	if err != nil {
		t.Fatal(err)
	}
	updatedRetrievedContact, err := api.RetrieveContact(updatedContact.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(updatedContact, updatedRetrievedContact) {
		spew.Dump(updatedRetrievedContact)
		t.Fatal("Updated contact is not equal!")
	}

	otherContactParams := &cm.NewContact{
		CustomerUUID:   cus1.UUID,
		DataSourceUUID: ds.UUID,
		FirstName:      "",
		LastName:       "",
		LinkedIn:       "",
		Notes:          "",
		Phone:          "",
		Position:       1,
		Title:          "",
		Twitter:        "",
	}
	otherContact, err := api.CreateContact(otherContactParams)
	if err != nil {
		t.Fatal(err)
	}
	mergedContact, err := api.MergeContacts(otherContact.UUID, updatedContact.UUID)
	if err != nil {
		t.Fatal(err)
	}

	if mergedContact.Notes != updatedContact.Notes {
		spew.Dump(mergedContact.Notes, updatedContact.Notes)
		t.Fatal("Contact is not equal!")
	}

	deleteErr := api.DeleteContact(mergedContact.UUID)
	if deleteErr != nil {
		t.Fatal(err)
	}
}
