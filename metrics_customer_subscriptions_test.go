package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// Test list activities tests.
func TestListCustomerSubscriptions(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/customers/cus_8bc55ab6-c3b5-11eb-ac45-2f9a49d75af7/subscriptions" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{"entries": [{
                                                "id": 9306830,
                                                "external_id": "sub_0001",
                                                "plan": "PRO Plan (10,000 active cust.) monthly",
                                                "quantity": 1,
                                                "mrr": 70800,
                                                "arr": 849600,
                                                "status": "active",
                                                "billing-cycle": "month",
                                                "billing-cycle-count": 1,
                                                "start-date": "2015-12-20T08:26:49-05:00",
                                                "end-date": "2016-03-20T09:26:49-05:00",
                                                "currency": "USD",
                                                "currency-sign": "$"
                                            }],
                                "has_more": false,
                                "per_page": 200,
                                "page": 1}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		AccountToken: "token",
		AccessKey:    "key",
	}
	activities, err := tested.MetricsListCustomerSubscriptions(&Cursor{}, "cus_8bc55ab6-c3b5-11eb-ac45-2f9a49d75af7")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(activities.Entries) == 0 ||
		activities.Entries[0].ExternalID != "sub_0001" {
		spew.Dump(activities)
		t.Fatal("Unexpected result")
	}
}
