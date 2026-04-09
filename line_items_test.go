package chartmogul

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const retrieveLineItemExample = `{
	"uuid": "li_d72e6843-5793-41d0-bfdf-0269514c9c56",
	"external_id": "li_ext_123",
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
	"disabled": false
}`

const createLineItemsExample = `{
	"line_items": [
		{
			"uuid": "li_new_123",
			"type": "one_time",
			"amount_in_cents": 2500,
			"description": "Setup Fee"
		}
	]
}`

func TestRetrieveLineItem(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "GET"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/line_items/li_d72e6843-5793-41d0-bfdf-0269514c9c56"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				w.Write([]byte(retrieveLineItemExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	lineItem, err := tested.RetrieveLineItem("li_d72e6843-5793-41d0-bfdf-0269514c9c56")

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
	if lineItem.UUID != "li_d72e6843-5793-41d0-bfdf-0269514c9c56" {
		t.Errorf("Expected UUID li_d72e6843-5793-41d0-bfdf-0269514c9c56, got: %v", lineItem.UUID)
	}
	if lineItem.AmountInCents != 5000 {
		t.Errorf("Expected amount_in_cents 5000, got: %v", lineItem.AmountInCents)
	}
}

func TestRetrieveLineItemWithParams(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Expected GET, got: %v", r.Method)
				}
				validationType := r.URL.Query().Get("validation_type")
				if validationType != "all" {
					t.Errorf("Expected validation_type=all, got: %v", validationType)
				}
				withDisabled := r.URL.Query().Get("with_disabled")
				if withDisabled != "true" {
					t.Errorf("Expected with_disabled=true, got: %v", withDisabled)
				}
				w.Write([]byte(retrieveLineItemExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	trueBool := true
	lineItem, err := tested.RetrieveLineItem("li_d72e6843-5793-41d0-bfdf-0269514c9c56", &RetrieveLineItemParams{
		ValidationType: "all",
		WithDisabled:   &trueBool,
	})

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
	if lineItem.UUID != "li_d72e6843-5793-41d0-bfdf-0269514c9c56" {
		t.Errorf("Expected UUID li_d72e6843-5793-41d0-bfdf-0269514c9c56, got: %v", lineItem.UUID)
	}
}

func TestCreateLineItems(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "POST"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/import/invoices/inv_123/line_items"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				w.Write([]byte(createLineItemsExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	result, err := tested.CreateLineItems("inv_123", []*LineItem{
		{
			Type:          "one_time",
			AmountInCents: 2500,
			Description:   "Setup Fee",
		},
	})

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
	if len(result.LineItems) != 1 {
		t.Fatalf("Expected 1 line item, got: %v", len(result.LineItems))
	}
	if result.LineItems[0].UUID != "li_new_123" {
		t.Errorf("Expected UUID li_new_123, got: %v", result.LineItems[0].UUID)
	}
}

func TestUpdateLineItem(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "PATCH"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/line_items/li_123"
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
				expectedBody := `{"amount_in_cents":7500}`
				if string(body) != expectedBody {
					t.Errorf("Requested body expected: %v, actual: %v", expectedBody, string(body))
				}
				w.Write([]byte(`{"uuid":"li_123","type":"subscription","amount_in_cents":7500}`)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	amount := 7500
	lineItem, err := tested.UpdateLineItem("li_123", &UpdateLineItemParams{
		AmountInCents: &amount,
	})

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
	if lineItem.UUID != "li_123" {
		t.Errorf("Expected UUID li_123, got: %v", lineItem.UUID)
	}
	if lineItem.AmountInCents != 7500 {
		t.Errorf("Expected amount_in_cents 7500, got: %v", lineItem.AmountInCents)
	}
}

func TestToggleLineItemDisabled(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "PATCH"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/line_items/li_123/disabled_state"
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
				w.Write([]byte(`{"uuid":"li_123","type":"subscription","amount_in_cents":5000,"disabled":true,"disabled_at":"2026-03-17T10:00:00.000Z","disabled_by":"user@example.com"}`)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	lineItem, err := tested.ToggleLineItemDisabled("li_123", &ToggleLineItemDisabledParams{
		Disabled: true,
	})

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
	if lineItem.UUID != "li_123" {
		t.Errorf("Expected UUID li_123, got: %v", lineItem.UUID)
	}
	if lineItem.Disabled == nil || *lineItem.Disabled != true {
		t.Error("Expected disabled to be true")
	}
	if lineItem.DisabledAt != "2026-03-17T10:00:00.000Z" {
		t.Errorf("Expected disabled_at, got: %v", lineItem.DisabledAt)
	}
}

func TestDeleteLineItem(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expected := "/v/line_items/li_123"
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
	err := tested.DeleteLineItem("li_123")

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
}
