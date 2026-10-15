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
	expectedHistoricalValues := map[string]interface{}{
		"custom": map[string]interface{}{
			"channel": []interface{}{
				map[string]interface{}{"value": "Google", "update_performed_at": nil, "update_performed_by": nil, "initial": true},
				map[string]interface{}{"value": "Facebook", "update_performed_at": "2026-09-25T10:00:00Z", "update_performed_by": "API", "initial": false},
			},
		},
	}
	if !reflect.DeepEqual(attributes.HistoricalValues, expectedHistoricalValues) {
		spew.Dump(attributes.HistoricalValues)
		t.Fatal("Unexpected historical_values")
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
		&AddCustomAttributesParams{Overrides: map[string]interface{}{"custom": map[string]interface{}{"channel": true}}},
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

func TestAddCustomAttributesWithEmailWithOptions(t *testing.T) {
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
	result, err := tested.AddCustomAttributesWithEmailWithOptions(
		"adam@example.com",
		[]*CustomAttribute{{Type: "String", Key: "channel", Value: "Facebook"}},
		&AddCustomAttributesWithEmailParams{Overrides: map[string]interface{}{"custom": map[string]interface{}{"channel": true}}},
	)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	expectedEntryOverrides := map[string]interface{}{
		"attributes": map[string]interface{}{"custom": map[string]interface{}{"channel": true}},
	}
	if len(result.Entries) != 1 || !reflect.DeepEqual(result.Entries[0].Overrides, expectedEntryOverrides) {
		spew.Dump(result)
		t.Fatal("Expected entry with customer-shaped overrides")
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
		&UpdateCustomAttributesParams{Overrides: map[string]interface{}{"custom": map[string]interface{}{"channel": true}}},
	)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if result.Custom["channel"] != "Twitter" {
		spew.Dump(result)
		t.Fatal("Unexpected custom attributes")
	}
	expectedOverrides := map[string]interface{}{
		"custom": map[string]interface{}{"channel": true},
	}
	if !reflect.DeepEqual(result.Overrides, expectedOverrides) {
		spew.Dump(result.Overrides)
		t.Fatal("Unexpected overrides")
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
		&RemoveCustomAttributesParams{Overrides: map[string]interface{}{"custom": map[string]interface{}{"age": false}}},
	)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if result.Message != "Custom attributes deleted from customer" {
		spew.Dump(result)
		t.Fatal("Unexpected message")
	}
	if !reflect.DeepEqual(result.Overrides, map[string]interface{}{}) {
		spew.Dump(result.Overrides)
		t.Fatal("Expected empty overrides")
	}
}

func TestLegacyCustomAttributeMethodsOmitOverrides(t *testing.T) {
	channelAttribute := []interface{}{
		map[string]interface{}{"type": "String", "key": "channel", "value": "Facebook"},
	}
	cases := []struct {
		name         string
		expectedVerb string
		expectedBody map[string]interface{}
		call         func(api *API) error
	}{
		{
			name:         "AddCustomAttributesToCustomer",
			expectedVerb: "POST",
			expectedBody: map[string]interface{}{"custom": channelAttribute},
			call: func(api *API) error {
				_, err := api.AddCustomAttributesToCustomer(
					"cus_00000000-0000-0000-0000-000000000000",
					[]*CustomAttribute{{Type: "String", Key: "channel", Value: "Facebook"}})
				return err
			},
		},
		{
			name:         "AddCustomAttributesWithEmail",
			expectedVerb: "POST",
			expectedBody: map[string]interface{}{"email": "adam@example.com", "custom": channelAttribute},
			call: func(api *API) error {
				_, err := api.AddCustomAttributesWithEmail(
					"adam@example.com",
					[]*CustomAttribute{{Type: "String", Key: "channel", Value: "Facebook"}})
				return err
			},
		},
		{
			name:         "UpdateCustomAttributesOfCustomer",
			expectedVerb: "PUT",
			expectedBody: map[string]interface{}{"custom": map[string]interface{}{"channel": "Twitter"}},
			call: func(api *API) error {
				_, err := api.UpdateCustomAttributesOfCustomer(
					"cus_00000000-0000-0000-0000-000000000000",
					map[string]interface{}{"channel": "Twitter"})
				return err
			},
		},
		{
			name:         "RemoveCustomAttributes",
			expectedVerb: "DELETE",
			expectedBody: map[string]interface{}{"custom": []interface{}{"channel"}},
			call: func(api *API) error {
				_, err := api.RemoveCustomAttributes(
					"cus_00000000-0000-0000-0000-000000000000",
					[]string{"channel"})
				return err
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						if r.Method != tc.expectedVerb {
							t.Errorf("Unexpected method %v", r.Method)
						}
						assertRequestBody(t, r, tc.expectedBody)
						w.Header().Set("Content-Type", "application/json")
						w.Write([]byte(`{"custom": {}, "overrides": {}}`)) //nolint
					}))
			defer server.Close()
			SetURL(server.URL + "/v/%v")

			if err := tc.call(&API{ApiKey: "token"}); err != nil {
				spew.Dump(err)
				t.Fatal("Not expected to fail")
			}
		})
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
