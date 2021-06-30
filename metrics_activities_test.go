package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// Test list activities tests.
func TestListActivities(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/activities?type=new_biz" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{"entries": [{
                                              "description": "purchased the plan_11 plan",
                                              "activity-mrr-movement": 6000,
                                              "activity-mrr": 6000,
                                              "activity-arr": 72000,
                                              "date": "2020-05-06T01:00:00",
                                              "type": "new_biz",
                                              "currency": "USD",
                                              "subscription-external-id": "sub_2",
                                              "plan-external-id": "11",
                                              "customer-name": "customer_2",
                                              "customer-uuid": "8bc55ab6-c3b5-11eb-ac45-2f9a49d75af7",
                                              "customer-external-id": "customer_2",
                                              "billing-connector-uuid": "99076cb8-97a1-11eb-8798-a73b507e7929",
                                              "uuid": "f1a49735-21c7-4e3f-9ddc-67927aaadcf4"
                                            }],
											"has_more": false,
										    "per_page": 200}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		AccountToken: "token",
		AccessKey:    "key",
	}
	activities, err := tested.MetricsListActivities(&MetricsListActivitiesParams{Type: "new_biz"})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(activities.Entries) == 0 ||
		activities.Entries[0].UUID != "f1a49735-21c7-4e3f-9ddc-67927aaadcf4" {
		spew.Dump(activities)
		t.Fatal("Unexpected result")
	}
}
