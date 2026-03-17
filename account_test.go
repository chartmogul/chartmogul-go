package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

const accountExample = `{
	"id": "acc_12345",
	"name": "Example Test Company",
	"currency": "EUR",
	"time_zone": "Europe/Berlin",
	"week_start_on": "sunday"
}`

const accountWithIncludeExample = `{
	"id": "acc_12345",
	"name": "Example Test Company",
	"currency": "EUR",
	"time_zone": "Europe/Berlin",
	"week_start_on": "sunday",
	"churn_recognition": {"mode": "at_period_end"},
	"churn_when_zero_mrr": {"enabled": true}
}`

func TestRetrieveAccount(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "GET"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/account"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				w.Write([]byte(accountExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	account, err := tested.RetrieveAccount()

	if account.ID != "acc_12345" || account.Name != "Example Test Company" || account.Currency != "EUR" || account.TimeZone != "Europe/Berlin" || account.WeekStartOn != "sunday" {
		spew.Dump(account)
		t.Error("Unexpected account details")
	}
	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestRetrieveAccountWithInclude(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "GET"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				// Check include query param
				include := r.URL.Query().Get("include")
				if include != "churn_recognition,churn_when_zero_mrr" {
					t.Errorf("Expected include param, got: %v", include)
				}
				w.Write([]byte(accountWithIncludeExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	account, err := tested.RetrieveAccount(&RetrieveAccountParams{
		Include: "churn_recognition,churn_when_zero_mrr",
	})

	if account.ID != "acc_12345" || account.Name != "Example Test Company" {
		spew.Dump(account)
		t.Error("Unexpected account details")
	}
	if account.ChurnRecognition == nil {
		t.Error("Expected churn_recognition to be present")
	}
	if account.ChurnWhenZeroMRR == nil {
		t.Error("Expected churn_when_zero_mrr to be present")
	}
	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}
