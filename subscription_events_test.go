package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

const firstSubscriptionEventExample = `{
	"id": 73966836,
	"data_source_uuid": "ds_1fm3eaac-62d0-31ec-clf4-4bf0mbe81aba",
	"customer_external_id": "scus_023",
	"subscription_set_external_id": "sub_set_ex_id_1",
	"subscription_external_id": "sub_0023",
	"plan_external_id": "p_ex_id_1",
	"event_date": "2022-04-09T11:17:14Z",
	"effective_date": "2022-04-09T10:04:13Z",
	"event_type": "subscription_cancelled",
	"external_id": "ex_id_1",
	"errors": {},
	"created_at": "2022-04-09T11:17:14Z",
	"updated_at": "2022-04-09T11:17:14Z",
	"quantity": 1,
	"currency": "USD",
	"amount_in_cents": 1000,
	"tax_amount_in_cents": 19,
	"retracted_event_id": null
}`

const secondSubscriptionEventExample = `{
	"id": 73966837,
	"data_source_uuid": "ds_1fm3eaac-62d0-31ec-clf4-4bf0mbe81aba",
	"customer_external_id": "scus_024",
	"subscription_set_external_id": "sub_set_ex_id_2",
	"subscription_external_id": "sub_0024",
	"plan_external_id": "p_ex_id_2",
	"event_date": "2022-04-09T11:17:14Z",
	"effective_date": "2022-04-09T10:04:13Z",
	"event_type": "subscription_cancelled",
	"external_id": "ex_id_2",
	"errors": {},
	"created_at": "2022-04-09T11:17:14Z",
	"updated_at": "2022-04-09T11:17:14Z",
	"quantity": 1,
	"currency": "USD",
	"amount_in_cents": 1000,
	"tax_amount_in_cents": 19,
	"retracted_event_id": null
}`

func TestListSubscriptionEvent(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/subscription_events" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(
					`{
						"subscription_events": [
							` + firstSubscriptionEventExample + `
						],
						"meta": {
							"next_key": 67048503,
							"prev_key": null,
							"before_key": "2022-04-10T22:27:35.834Z",
							"page": 1,
							"total_pages": 166
						}
					}`,
				))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	actual, err := tested.ListSubscriptionEvents(&FilterSubscriptionEvents{}, &MetaCursor{})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(actual.SubscriptionEvents) != 1 ||
		actual.SubscriptionEvents[0].ID != 73966836 ||
		actual.SubscriptionEvents[0].DataSourceUUID != "ds_1fm3eaac-62d0-31ec-clf4-4bf0mbe81aba" ||
		actual.SubscriptionEvents[0].CustomerExternalID != "scus_023" ||
		actual.SubscriptionEvents[0].EventDate != "2022-04-09T11:17:14Z" ||
		actual.SubscriptionEvents[0].EventType != "subscription_cancelled" ||
		actual.SubscriptionEvents[0].CreatedAt != "2022-04-09T11:17:14Z" ||
		actual.SubscriptionEvents[0].UpdatedAt != "2022-04-09T11:17:14Z" ||
		actual.SubscriptionEvents[0].EffectiveDate != "2022-04-09T10:04:13Z" {
		spew.Dump(actual)
		t.Fatal("Unexpected result")
	}
}

func TestFilteredListSubscriptionEvent(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/subscription_events?external_id=ex_id_2" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(
					`{
						"subscription_events": [
							` + secondSubscriptionEventExample + `
						],
						"meta": {
							"next_key": 67048503,
							"prev_key": null,
							"before_key": "2022-04-10T22:27:35.834Z",
							"page": 1,
							"total_pages": 166
						}
					}`,
				))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	actual, err := tested.ListSubscriptionEvents(&FilterSubscriptionEvents{ExternalID: "ex_id_2"}, &MetaCursor{})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(actual.SubscriptionEvents) != 1 ||
		actual.SubscriptionEvents[0].ID != 73966837 {
		spew.Dump(actual)
		t.Fatal("Unexpected result")
	}
}

func TestDeleteSubscriptionEventById(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/subscription_events" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}

				w.WriteHeader(http.StatusNoContent)
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	err := tested.DeleteSubscriptionEvent(&DeleteSubscriptionEvent{ID: 123})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestDeleteSubscriptionEventByExternalIdAndDataSourceUuid(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/subscription_events" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	err := tested.DeleteSubscriptionEvent(&DeleteSubscriptionEvent{DataSourceUUID: "ds_123", ExternalID: "ex_id_1"})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}
