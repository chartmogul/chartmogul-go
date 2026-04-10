package chartmogul

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

const retrieveTransactionExample = `{
	"uuid": "tr_879d560a-1bec-41bb-986e-665e38a2f7bc",
	"external_id": "tx_ext_123",
	"type": "payment",
	"date": "2015-11-05T00:14:23.000Z",
	"result": "successful",
	"disabled": false,
	"disabled_at": "",
	"disabled_by": "",
	"user_created": false
}`

func TestRetrieveTransaction(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "GET"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/transactions/tr_879d560a-1bec-41bb-986e-665e38a2f7bc"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				w.Write([]byte(retrieveTransactionExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	transaction, err := tested.RetrieveTransaction("tr_879d560a-1bec-41bb-986e-665e38a2f7bc")

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
	if transaction.UUID != "tr_879d560a-1bec-41bb-986e-665e38a2f7bc" {
		t.Errorf("Expected UUID tr_879d560a-1bec-41bb-986e-665e38a2f7bc, got: %v", transaction.UUID)
	}
	if transaction.Type != "payment" {
		t.Errorf("Expected type payment, got: %v", transaction.Type)
	}
}

func TestRetrieveTransactionWithParams(t *testing.T) {
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
				w.Write([]byte(retrieveTransactionExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	transaction, err := tested.RetrieveTransaction("tr_879d560a-1bec-41bb-986e-665e38a2f7bc", &RetrieveTransactionParams{
		ValidationType: "all",
	})

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
	if transaction.UUID != "tr_879d560a-1bec-41bb-986e-665e38a2f7bc" {
		t.Errorf("Expected UUID tr_879d560a-1bec-41bb-986e-665e38a2f7bc, got: %v", transaction.UUID)
	}
}

func TestUpdateTransaction(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "PATCH"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/transactions/tr_123"
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
				expectedBody := `{"date":"2015-12-01T00:00:00.000Z","result":"failed"}`
				if string(body) != expectedBody {
					t.Errorf("Requested body expected: %v, actual: %v", expectedBody, string(body))
				}
				w.Write([]byte(`{"uuid":"tr_123","type":"payment","date":"2015-12-01T00:00:00.000Z","result":"failed"}`)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	transaction, err := tested.UpdateTransaction("tr_123", &UpdateTransactionParams{
		Date:   "2015-12-01T00:00:00.000Z",
		Result: "failed",
	})

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
	if transaction.UUID != "tr_123" {
		t.Errorf("Expected UUID tr_123, got: %v", transaction.UUID)
	}
	if transaction.Result != "failed" {
		t.Errorf("Expected result failed, got: %v", transaction.Result)
	}
}

func TestToggleTransactionDisabled(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "PATCH"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/transactions/tr_123/disabled_state"
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
				w.Write([]byte(`{"uuid":"tr_123","type":"payment","date":"2015-11-05T00:14:23.000Z","result":"successful","disabled":true,"disabled_at":"2026-03-17T10:00:00.000Z","disabled_by":"user@example.com"}`)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	transaction, err := tested.ToggleTransactionDisabled("tr_123", &ToggleTransactionDisabledParams{
		Disabled: true,
	})

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
	if transaction.UUID != "tr_123" {
		t.Errorf("Expected UUID tr_123, got: %v", transaction.UUID)
	}
	if transaction.Disabled == nil || *transaction.Disabled != true {
		t.Error("Expected disabled to be true")
	}
}

func TestDeleteTransaction(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expected := "/v/transactions/tr_123"
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
	err := tested.DeleteTransaction("tr_123")

	if err != nil {
		t.Fatal("Not expected to fail:", err)
	}
}
