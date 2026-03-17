package chartmogul

import (
	"io/ioutil"
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
			"transaction_fees_currency": "EUR",
			"discount_description": "5 EUR",
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
			"transaction_fees_currency": "EUR",
			"discount_description": "2 EUR",
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
	"disabled": false,
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
			"transaction_fees_currency": "EUR",
			"discount_description": "5 EUR",
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
	],
	"errors": {}
}`

const retrieveInvoiceWithEditHistoryExample = `{
	"uuid": "inv_123",
	"external_id": "INV0001",
	"date": "2015-11-01T00:00:00.000Z",
	"due_date": "2015-11-15T00:00:00.000Z",
	"currency": "USD",
	"disabled": false,
	"disabled_at": "",
	"disabled_by": "",
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
			"transaction_fees_currency": "EUR",
			"discount_description": "5 EUR",
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
	],
	"edit_history_summary": {
		"values_changed": {
			"date": {
				"original_value": "2015-07-31T21:39:06.000Z",
				"edited_value": "2015-07-14T21:39:06.000Z"
			},
			"amount_in_cents": {
				"original_value": 11000,
				"edited_value": 10000
			}
		},
		"latest_edit_author": "test@example.com",
		"latest_edit_performed_at": "2015-08-04T12:36:50.574Z"
	},
	"errors": {}
}`

const retrieveInvoiceWithErrorsExample = `{
	"uuid": "inv_456",
	"external_id": "INV0002",
	"date": "2015-11-01T00:00:00.000Z",
	"due_date": "2015-11-15T00:00:00.000Z",
	"currency": "USD",
	"disabled": false,
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
			"transaction_fees_currency": "EUR",
			"discount_description": "5 EUR",
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
	],
	"errors": {
		"currency": [
			"The invoice currency must be specified for each invoice."
		],
		"date": [
			"Invoice date is required.",
			"Invoice date must be in the past."
		]
	}
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
		ApiKey: "token",
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

func TestListAllInvoicesWithValidationType(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Check query parameter
				validationType := r.URL.Query().Get("validation_type")
				if validationType != "all" {
					t.Errorf("Expected validation_type=all, got: %v", validationType)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(listAllInvoicesExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	result, err := tested.ListAllInvoices(&ListAllInvoicesParams{
		ValidationType: "all",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(result.Invoices) != 1 {
		spew.Dump(result)
		t.Fatal("Unexpected values")
	}
}

func TestListAllInvoicesWithAllParams(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				// Check query parameters
				query := r.URL.Query()
				if query.Get("validation_type") != "valid" {
					t.Errorf("Expected validation_type=valid, got: %v", query.Get("validation_type"))
				}
				if query.Get("include_edit_histories") != "true" {
					t.Errorf("Expected include_edit_histories=true, got: %v", query.Get("include_edit_histories"))
				}
				if query.Get("with_disabled") != "true" {
					t.Errorf("Expected with_disabled=true, got: %v", query.Get("with_disabled"))
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(listAllInvoicesExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	trueBool := true
	result, err := tested.ListAllInvoices(&ListAllInvoicesParams{
		ValidationType:       "valid",
		IncludeEditHistories: &trueBool,
		WithDisabled:         &trueBool,
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(result.Invoices) != 1 {
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
		ApiKey: "token",
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
		ApiKey: "token",
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

func TestRetrieveInvoiceWithValidationType(t *testing.T) {
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
				// Check query parameter
				validationType := r.URL.Query().Get("validation_type")
				if validationType != "all" {
					t.Errorf("Expected validation_type=all, got: %v", validationType)
				}
				w.Write([]byte(retrieveInvoiceExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	invoice, err := tested.RetrieveInvoice("inv_123", &RetrieveInvoiceParams{
		ValidationType: "all",
	})

	if len(invoice.LineItems) != 1 || len(invoice.Transactions) != 1 || invoice.UUID != "inv_123" {
		spew.Dump(invoice)
		t.Error("Unexpected invoice")
	}
	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestRetrieveInvoiceWithAllParams(t *testing.T) {
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
				// Check query parameters
				query := r.URL.Query()
				if query.Get("validation_type") != "invalid" {
					t.Errorf("Expected validation_type=invalid, got: %v", query.Get("validation_type"))
				}
				if query.Get("include_edit_histories") != "true" {
					t.Errorf("Expected include_edit_histories=true, got: %v", query.Get("include_edit_histories"))
				}
				if query.Get("with_disabled") != "false" {
					t.Errorf("Expected with_disabled=false, got: %v", query.Get("with_disabled"))
				}
				w.Write([]byte(retrieveInvoiceWithEditHistoryExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	trueBool := true
	falseBool := false
	invoice, err := tested.RetrieveInvoice("inv_123", &RetrieveInvoiceParams{
		ValidationType:       "invalid",
		IncludeEditHistories: &trueBool,
		WithDisabled:         &falseBool,
	})

	if len(invoice.LineItems) != 1 || len(invoice.Transactions) != 1 || invoice.UUID != "inv_123" {
		spew.Dump(invoice)
		t.Error("Unexpected invoice")
	}

	// Check new fields
	if invoice.Disabled == nil || *invoice.Disabled != false {
		t.Error("Expected disabled field to be false")
	}

	if invoice.EditHistorySummary == nil {
		t.Error("Expected edit_history_summary to be present")
	} else {
		if invoice.EditHistorySummary.LatestEditAuthor != "test@example.com" {
			t.Errorf("Expected latest_edit_author to be 'test@example.com', got: %v", invoice.EditHistorySummary.LatestEditAuthor)
		}
		if invoice.EditHistorySummary.LatestEditPerformedAt != "2015-08-04T12:36:50.574Z" {
			t.Errorf("Expected latest_edit_performed_at to be '2015-08-04T12:36:50.574Z', got: %v", invoice.EditHistorySummary.LatestEditPerformedAt)
		}
		if len(invoice.EditHistorySummary.ValuesChanged) != 2 {
			t.Errorf("Expected 2 values changed, got: %v", len(invoice.EditHistorySummary.ValuesChanged))
		}
	}

	if invoice.Errors == nil {
		t.Error("Expected errors field to be present")
	}

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestRetrieveInvoiceWithErrors(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "GET"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/invoices/inv_456"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				w.Write([]byte(retrieveInvoiceWithErrorsExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	invoice, err := tested.RetrieveInvoice("inv_456")

	if invoice.UUID != "inv_456" {
		t.Errorf("Expected UUID to be 'inv_456', got: %v", invoice.UUID)
	}

	// Check errors field
	if invoice.Errors == nil {
		t.Fatal("Expected errors field to be present")
	}

	currencyErrors := (*invoice.Errors)["currency"]
	if len(currencyErrors) != 1 {
		t.Errorf("Expected 1 currency error, got: %v", len(currencyErrors))
	} else if currencyErrors[0] != "The invoice currency must be specified for each invoice." {
		t.Errorf("Unexpected currency error message: %v", currencyErrors[0])
	}

	dateErrors := (*invoice.Errors)["date"]
	if len(dateErrors) != 2 {
		t.Errorf("Expected 2 date errors, got: %v", len(dateErrors))
	} else {
		if dateErrors[0] != "Invoice date is required." {
			t.Errorf("Unexpected first date error message: %v", dateErrors[0])
		}
		if dateErrors[1] != "Invoice date must be in the past." {
			t.Errorf("Unexpected second date error message: %v", dateErrors[1])
		}
	}

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestUpdateInvoiceStatus(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "PATCH"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/data_sources/ds_123/invoices/INV0001/status"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				defer r.Body.Close()
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatal("error should be nil")
				}
				expectedBody := `{"status":"void"}`
				if string(body) != expectedBody {
					t.Errorf("Requested body expected: %v, actual: %v", expectedBody, string(body))
				}
				w.Write([]byte("{}")) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	err := tested.UpdateInvoiceStatus("ds_123", "INV0001", &UpdateInvoiceStatusParams{
		Status: "void",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestToggleInvoiceDisabled(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "PATCH"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/invoices/inv_123/disabled_state"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				defer r.Body.Close()
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatal("error should be nil")
				}
				expectedBody := `{"disabled":true}`
				if string(body) != expectedBody {
					t.Errorf("Requested body expected: %v, actual: %v", expectedBody, string(body))
				}
				w.Write([]byte(`{"uuid":"inv_123","external_id":"INV0001","date":"2015-11-01T00:00:00.000Z","currency":"USD","disabled":true,"disabled_at":"2026-03-17T10:00:00.000Z","disabled_by":"user@example.com","line_items":[]}`)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	invoice, err := tested.ToggleInvoiceDisabled("inv_123", &ToggleInvoiceDisabledParams{
		Disabled: true,
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if invoice.UUID != "inv_123" {
		t.Errorf("Expected UUID inv_123, got: %v", invoice.UUID)
	}
	if invoice.Disabled == nil || *invoice.Disabled != true {
		t.Error("Expected disabled to be true")
	}
	if invoice.DisabledAt != "2026-03-17T10:00:00.000Z" {
		t.Errorf("Expected disabled_at, got: %v", invoice.DisabledAt)
	}
}

func TestRetrieveInvoiceWithLineItemErrors(t *testing.T) {
	const lineItemErrorsExample = `{
		"uuid": "inv_789",
		"external_id": "INV0003",
		"date": "2015-11-01T00:00:00.000Z",
		"currency": "USD",
		"disabled": false,
		"line_items": [
			{
				"uuid": "li_abc",
				"type": "subscription",
				"amount_in_cents": 5000,
				"errors": {"service_period_start": ["is required"]}
			}
		],
		"transactions": []
	}`
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(lineItemErrorsExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	invoice, err := tested.RetrieveInvoice("inv_789")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(invoice.LineItems) != 1 {
		t.Fatalf("Expected 1 line item, got: %v", len(invoice.LineItems))
	}
	if invoice.LineItems[0].Errors == nil {
		t.Error("Expected line item errors to be present")
	}
}

func TestCreateInvoiceFullRefund(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "POST"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/import/customers/uuid/invoices"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				defer r.Body.Close()
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					spew.Dump(err)
					t.Fatal("error should be nil")
				}
				expected = `{"invoices":[{"currency":"","date":"","external_id":"","line_items":null,"transactions":[{"date":"","external_id":"tx-id","result":"successful","type":"payment"}]}]}`
				if string(body) != expected {
					t.Errorf("Requested body expected: %v, actual: %v", expected, string(body))
				}
				w.Write([]byte("{}")) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	_, err := tested.CreateInvoices([]*Invoice{{
		Transactions: []*Transaction{{
			ExternalID: "tx-id",
			Type:       "payment",
			Result:     "successful",
		}},
	}}, "uuid")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestCreateInvoicePartialRefund(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "POST"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/import/customers/uuid/invoices"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				defer r.Body.Close()
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					spew.Dump(err)
					t.Fatal("error should be nil")
				}
				expected = `{"invoices":[{"currency":"","date":"","external_id":"","line_items":null,"transactions":[{"amount_in_cents":200,"date":"","external_id":"tx-id","result":"successful","type":"payment"}]}]}`
				if string(body) != expected {
					t.Errorf("Requested body expected: %v, actual: %v", expected, string(body))
				}
				w.Write([]byte("{}")) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	amount := 200
	_, err := tested.CreateInvoices([]*Invoice{{
		Transactions: []*Transaction{{
			AmountInCents: &amount,
			ExternalID:    "tx-id",
			Type:          "payment",
			Result:        "successful",
		}},
	}}, "uuid")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}
