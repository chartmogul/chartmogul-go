package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

const excludedCustomersListExample = `{
	"entries": [
		{
			"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
			"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000",
			"external_id": "ext-1",
			"name": "Acme",
			"email": "acme@example.com",
			"excluded_at": "2026-11-20T10:00:00.000Z"
		}
	],
	"cursor": "cursor-2",
	"has_more": true
}`

func TestListExcludedCustomers(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/excluded_customers?cursor=cursor-1&per_page=1" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(excludedCustomersListExample))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	params := &ListExcludedCustomersParams{Cursor: Cursor{PerPage: 1, Cursor: "cursor-1"}}
	excludedCustomers, err := tested.ListExcludedCustomers(params)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	expected := &ExcludedCustomers{
		Entries: []*ExcludedCustomer{{
			CustomerUUID:   "cus_00000000-0000-0000-0000-000000000000",
			DataSourceUUID: "ds_00000000-0000-0000-0000-000000000000",
			ExternalID:     "ext-1",
			Name:           "Acme",
			Email:          "acme@example.com",
			ExcludedAt:     "2026-11-20T10:00:00.000Z",
		}},
		Pagination: Pagination{Cursor: "cursor-2", HasMore: true},
	}
	if !reflect.DeepEqual(expected, excludedCustomers) {
		spew.Dump(excludedCustomers)
		t.Fatal("Unexpected result")
	}
}

func TestListExcludedCustomersWithoutParams(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.RequestURI != "/v/excluded_customers" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{"entries": [], "cursor": "", "has_more": false}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	excludedCustomers, err := tested.ListExcludedCustomers(nil)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(excludedCustomers.Entries) != 0 || excludedCustomers.HasMore {
		spew.Dump(excludedCustomers)
		t.Fatal("Unexpected result")
	}
}

func TestReincludeCustomer(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/excluded_customers/cus_00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusAccepted)
				//nolint
				w.Write([]byte(`{"customer_uuid": "cus_00000000-0000-0000-0000-000000000000", "status": "manual_reimport_required"}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	result, err := tested.ReincludeCustomer("cus_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	expected := &ReincludedCustomer{
		CustomerUUID: "cus_00000000-0000-0000-0000-000000000000",
		Status:       ReincludedCustomerStatusManualReimportRequired,
	}
	if !reflect.DeepEqual(expected, result) {
		spew.Dump(result)
		t.Fatal("Unexpected result")
	}
}

func TestReincludeCustomerNotExcluded(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				//nolint
				w.Write([]byte(`{"code": 404, "message": "Excluded customer does not exist", "param": "customer_uuid"}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	_, err := tested.ReincludeCustomer("cus_00000000-0000-0000-0000-000000000000")

	httpErr, ok := errors.Cause(err).(HTTPError)
	if !ok || httpErr.StatusCode() != http.StatusNotFound {
		spew.Dump(err)
		t.Fatal("Expected a 404 HTTPError")
	}
}
