package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

const accountExample = `{
	"name": "Example Test Company",
	"currency": "EUR",
	"time_zone": "Europe/Berlin",
	"week_start_on": "sunday"
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
		AccountToken: "token",
		AccessKey:    "key",
	}
	account, err := tested.RetrieveAccount()

	if account.Name != "Example Test Company" || account.Currency != "EUR" || account.TimeZone != "Europe/Berlin" || account.WeekStartOn != "sunday" {
		spew.Dump(account)
		t.Error("Unexpected account details")
	}
	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}
