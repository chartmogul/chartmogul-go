package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// Test list activities tests.
func TestListCustomerActivities(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/customers/cus_8bc55ab6-c3b5-11eb-ac45-2f9a49d75af7/activities" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{"entries": [{
									"activity-arr": 24000,
									"activity-mrr": 2000,
									"activity-mrr-movement": 2000,
									"currency": "USD",
									"currency-sign": "$",
									"date": "2015-06-09T13:16:00+00:00",
									"description": "purchased the Silver Monthly plan (1)",
									"id": 48730,
									"type": "new_biz",
									"subscription-external-id": "1"
								}],
								"has_more": false,
								"per_page": 200,
								"page": 1}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	activities, err := tested.MetricsListCustomerActivities(&PaginationWithCursor{}, "cus_8bc55ab6-c3b5-11eb-ac45-2f9a49d75af7")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(activities.Entries) == 0 ||
		activities.Entries[0].Date != "2015-06-09T13:16:00+00:00" {
		spew.Dump(activities)
		t.Fatal("Unexpected result")
	}
}
