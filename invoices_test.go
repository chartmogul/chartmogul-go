package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

const oneInvoiceExample = `{
	"uuid": "inv_565c73b2-85b9-49c9-a25e-2b7df6a677c9",
	"customer_uuid": "cus_f466e33d-ff2b-4a11-8f85-417eb02157a7",
	"external_id": "INV0001",
	"date": "2015-11-01T00:00:00.000Z",
	"due_date": "2015-11-15T00:00:00.000Z",
	"currency": "USD",
	"line_items": [
		{
			"uuid": "li_d72e6843-5793-41d0-bfdf-0269514c9c56",
			"external_id": null,
			"type": "subscription",
			"subscription_uuid": "sub_e6bc5407-e258-4de0-bb43-61faaf062035",
			"subscription_external_id": "sub_external_id_123",
			"plan_uuid": "pl_eed05d54-75b4-431b-adb2-eb6b9e543206",
			"prorated": false,
			"service_period_start": "2015-11-01T00:00:00.000Z",
			"service_period_end": "2015-12-01T00:00:00.000Z",
			"amount_in_cents": 5000,
			"quantity": 1,
			"discount_code": "PSO86",
			"discount_amount_in_cents": 1000,
			"tax_amount_in_cents": 900,
			"account_code": null
		},
		{
			"uuid": "li_0cc8c112-beac-416d-af11-f35744ca4e83",
			"external_id": null,
			"type": "one_time",
			"description": "Setup Fees",
			"amount_in_cents": 2500,
			"quantity": 1,
			"discount_code": "PSO86",
			"discount_amount_in_cents": 500,
			"tax_amount_in_cents": 450,
			"account_code": null
		}
	],
	"transactions": [
		{
			"uuid": "tr_879d560a-1bec-41bb-986e-665e38a2f7bc",
			"external_id": null,
			"type": "payment",
			"date": "2015-11-05T00:14:23.000Z",
			"result": "successful"
		}
	]
}`

const listAllInvoicesExample = `{
  "invoices": [
    ` + oneInvoiceExample + `
  ],
  "current_page": 1,
  "total_pages": 1
}`

const retrieveInvoiceExample = `{
	"uuid": "inv_123",
	"external_id": "INV0001",
	"date": "2015-11-01T00:00:00.000Z",
	"due_date": "2015-11-15T00:00:00.000Z",
	"currency": "USD",
	"line_items": [
		{
			"uuid": "li_d72e6843-5793-41d0-bfdf-0269514c9c56",
			"external_id": null,
			"type": "subscription",
			"subscription_uuid": "sub_e6bc5407-e258-4de0-bb43-61faaf062035",
			"subscription_external_id": "sub_external_id_123",
			"subscription_set_external_id": "set_external_id_123",
			"plan_uuid": "pl_eed05d54-75b4-431b-adb2-eb6b9e543206",
			"prorated": false,
			"service_period_start": "2015-11-01T00:00:00.000Z",
			"service_period_end": "2015-12-01T00:00:00.000Z",
			"amount_in_cents": 5000,
			"quantity": 1,
			"discount_code": "PSO86",
			"discount_amount_in_cents": 1000,
			"tax_amount_in_cents": 900,
			"account_code": null
		}
	],
	"transactions": [
		{
			"uuid": "tr_879d560a-1bec-41bb-986e-665e38a2f7bc",
			"external_id": null,
			"type": "payment",
			"date": "2015-11-05T00:14:23.000Z",
			"result": "successful"
		}
	]
}`

func TestNewInvoicesAllListing(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(listAllInvoicesExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		AccountToken: "token",
		AccessKey:    "key",
	}
	result, err := tested.ListAllInvoices(&ListAllInvoicesParams{
		CustomerUUID: "cus_f466e33d-ff2b-4a11-8f85-417eb02157a7",
		ExternalID:   "INV0001",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(result.Invoices) != 1 ||
		result.Invoices[0].CustomerUUID != "cus_f466e33d-ff2b-4a11-8f85-417eb02157a7" ||
		result.Invoices[0].LineItems[0].AmountInCents != 5000 {
		spew.Dump(result)
		t.Fatal("Unexpected values")
	}
}

func TestDeleteInvoice(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expected := "/v/invoices/inv_123"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		AccountToken: "token",
		AccessKey:    "key",
	}
	err := tested.DeleteInvoice("inv_123")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestRetrieveInvoice(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "GET"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/invoices/inv_123"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				w.Write([]byte(retrieveInvoiceExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		AccountToken: "token",
		AccessKey:    "key",
	}
	invoice, err := tested.RetrieveInvoice("inv_123")

	if len(invoice.LineItems) != 1 || len(invoice.Transactions) != 1 || invoice.UUID != "inv_123" {
		spew.Dump(invoice)
		t.Error("Unexpected invoice")
	}
	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}
