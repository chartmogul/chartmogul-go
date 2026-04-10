package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const createJsonImportExample = `{
	"id": "imp_12345",
	"data_source_uuid": "ds_abc123",
	"status": "queued",
	"external_id": "batch_001",
	"created_at": "2026-03-17T10:00:00.000Z",
	"updated_at": "2026-03-17T10:00:00.000Z"
}`

const retrieveJsonImportExample = `{
	"id": "imp_12345",
	"data_source_uuid": "ds_abc123",
	"status": "completed",
	"external_id": "batch_001",
	"created_at": "2026-03-17T10:00:00.000Z",
	"updated_at": "2026-03-17T10:05:00.000Z"
}`

func TestCreateJsonImport(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "POST"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/data_sources/ds_abc123/json_imports"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				w.Write([]byte(createJsonImportExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	result, err := tested.CreateJsonImport("ds_abc123", &JsonImportData{
		ExternalID: "batch_001",
		Customers: []*JsonImportCustomer{
			{
				ExternalID: "cus_ext_001",
				Name:       "Test Customer",
				Email:      "test@example.com",
			},
		},
		Plans: []*JsonImportPlan{
			{
				Name:          "Basic Plan",
				IntervalCount: 1,
				IntervalUnit:  "month",
				ExternalID:    "plan_ext_001",
			},
		},
	})

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
	if result.ID != "imp_12345" {
		t.Errorf("Expected ID imp_12345, got: %v", result.ID)
	}
	if result.Status != "queued" {
		t.Errorf("Expected status queued, got: %v", result.Status)
	}
	if result.ExternalID != "batch_001" {
		t.Errorf("Expected external_id batch_001, got: %v", result.ExternalID)
	}
}

func TestRetrieveJsonImport(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "GET"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/data_sources/ds_abc123/json_imports/imp_12345"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				w.Write([]byte(retrieveJsonImportExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	result, err := tested.RetrieveJsonImport("ds_abc123", "imp_12345")

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
	if result.ID != "imp_12345" {
		t.Errorf("Expected ID imp_12345, got: %v", result.ID)
	}
	if result.Status != "completed" {
		t.Errorf("Expected status completed, got: %v", result.Status)
	}
}
