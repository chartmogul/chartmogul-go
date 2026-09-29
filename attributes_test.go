package chartmogul

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestRetrieveCustomersAttributesWithOptions(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				query := r.URL.Query()
				if query.Get("with_overrides") != "true" {
					t.Errorf("Expected with_overrides=true, got: %v", query.Get("with_overrides"))
				}
				if query.Get("attributes_with_history") != "custom.channel" {
					t.Errorf("Expected attributes_with_history=custom.channel, got: %v", query.Get("attributes_with_history"))
				}
				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write([]byte(`{
					"tags": ["important"],
					"custom": {"channel": "Facebook"},
					"overrides": {"custom": {"channel": true}},
					"historical_values": {"custom": {"channel": [
						{"value": "Google", "update_performed_at": null, "update_performed_by": null, "initial": true},
						{"value": "Facebook", "update_performed_at": "2026-09-25T10:00:00Z", "update_performed_by": "API", "initial": false}
					]}}
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	trueBool := true
	attributes, err := tested.RetrieveCustomersAttributesWithOptions("cus_00000000-0000-0000-0000-000000000000", &RetrieveCustomersAttributesParams{
		WithOverrides:         &trueBool,
		AttributesWithHistory: "custom.channel",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	expectedOverrides := map[string]interface{}{
		"custom": map[string]interface{}{"channel": true},
	}
	if !reflect.DeepEqual(attributes.Overrides, expectedOverrides) {
		spew.Dump(attributes.Overrides)
		t.Fatal("Unexpected overrides")
	}
	if attributes.HistoricalValues == nil {
		spew.Dump(attributes)
		t.Fatal("Expected historical_values")
	}
}

func TestRetrieveCustomersAttributesWithOptionsNil(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.URL.RawQuery != "" {
					t.Errorf("Expected no query, got: %v", r.URL.RawQuery)
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"custom": {"channel": "Facebook"}}`)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	attributes, err := tested.RetrieveCustomersAttributesWithOptions("cus_00000000-0000-0000-0000-000000000000", nil)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if attributes.Overrides != nil {
		spew.Dump(attributes)
		t.Fatal("Expected no overrides")
	}
}

func TestAddCustomAttributesToCustomerWithOptions(t *testing.T) {
	expectedBody := map[string]interface{}{
		"custom": []interface{}{
			map[string]interface{}{"type": "String", "key": "channel", "value": "Facebook"},
		},
		"overrides": map[string]interface{}{
			"custom": map[string]interface{}{"channel": true},
		},
	}

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				assertRequestBody(t, r, expectedBody)
				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write([]byte(`{
					"custom": {"channel": "Facebook"},
					"overrides": {"custom": {"channel": true}}
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	result, err := tested.AddCustomAttributesToCustomerWithOptions(
		"cus_00000000-0000-0000-0000-000000000000",
		[]*CustomAttribute{{Type: "String", Key: "channel", Value: "Facebook"}},
		&AddCustomAttributesOptions{Overrides: map[string]interface{}{"custom": map[string]interface{}{"channel": true}}},
	)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if !reflect.DeepEqual(result.Overrides, map[string]interface{}{"custom": map[string]interface{}{"channel": true}}) {
		spew.Dump(result)
		t.Fatal("Unexpected overrides")
	}
}

func TestAddCustomAttributesByEmailWithOptions(t *testing.T) {
	expectedBody := map[string]interface{}{
		"email": "adam@example.com",
		"custom": []interface{}{
			map[string]interface{}{"type": "String", "key": "channel", "value": "Facebook"},
		},
		"overrides": map[string]interface{}{
			"custom": map[string]interface{}{"channel": true},
		},
	}

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				assertRequestBody(t, r, expectedBody)
				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write([]byte(`{
					"entries": [{
						"uuid": "cus_00000000-0000-0000-0000-000000000000",
						"overrides": {"attributes": {"custom": {"channel": true}}}
					}]
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	result, err := tested.AddCustomAttributesByEmailWithOptions(
		"adam@example.com",
		[]*CustomAttribute{{Type: "String", Key: "channel", Value: "Facebook"}},
		&AddCustomAttributesWithEmailOptions{Overrides: map[string]interface{}{"custom": map[string]interface{}{"channel": true}}},
	)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(result.Entries) != 1 || result.Entries[0].Overrides == nil {
		spew.Dump(result)
		t.Fatal("Expected entry with overrides")
	}
}

func TestUpdateCustomAttributesOfCustomerWithOptions(t *testing.T) {
	expectedBody := map[string]interface{}{
		"custom": map[string]interface{}{"channel": "Twitter"},
		"overrides": map[string]interface{}{
			"custom": map[string]interface{}{"channel": true},
		},
	}

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PUT" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				assertRequestBody(t, r, expectedBody)
				w.Header().Set("Content-Type", "application/json")
				//nolint
				w.Write([]byte(`{
					"custom": {"channel": "Twitter"},
					"overrides": {"custom": {"channel": true}}
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	result, err := tested.UpdateCustomAttributesOfCustomerWithOptions(
		"cus_00000000-0000-0000-0000-000000000000",
		map[string]interface{}{"channel": "Twitter"},
		&UpdateCustomAttributesOptions{Overrides: map[string]interface{}{"custom": map[string]interface{}{"channel": true}}},
	)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if result.Custom["channel"] != "Twitter" {
		spew.Dump(result)
		t.Fatal("Unexpected custom attributes")
	}
}

func TestRemoveCustomAttributesWithOptions(t *testing.T) {
	expectedBody := map[string]interface{}{
		"custom": []interface{}{"age"},
		"overrides": map[string]interface{}{
			"custom": map[string]interface{}{"age": false},
		},
	}

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				assertRequestBody(t, r, expectedBody)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				//nolint
				w.Write([]byte(`{
					"custom": {"channel": "Facebook"},
					"overrides": {},
					"message": "Custom attributes deleted from customer"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	result, err := tested.RemoveCustomAttributesWithOptions(
		"cus_00000000-0000-0000-0000-000000000000",
		[]string{"age"},
		&RemoveCustomAttributesOptions{Overrides: map[string]interface{}{"custom": map[string]interface{}{"age": false}}},
	)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if result.Message != "Custom attributes deleted from customer" {
		spew.Dump(result)
		t.Fatal("Unexpected message")
	}
}

func assertRequestBody(t *testing.T, r *http.Request, expected map[string]interface{}) {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
		return
	}
	var incoming interface{}
	if err := json.Unmarshal(body, &incoming); err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(expected, incoming) {
		spew.Dump(expected, incoming)
		t.Error("Request body doesn't equal expected value")
	}
}
