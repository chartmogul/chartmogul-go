package chartmogul

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestConnectSubscriptions(t *testing.T) {
	expected := map[string]interface{}{
		"subscriptions": []interface{}{
			map[string]interface{}{
				"data_source_uuid": "ds_uuid1",
				"external_id":      "ext_id1",
			},
			map[string]interface{}{
				"data_source_uuid": "ds_uuid2",
				"external_id":      "ext_id2",
			},
		},
	}
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(202)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte("{}")) //nolint

				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
					return
				}
				var incoming interface{}
				err = json.Unmarshal(body, &incoming)
				if err != nil {
					t.Error(err)
					return
				}
				if !reflect.DeepEqual(expected, incoming) {
					spew.Dump(expected, incoming)
					t.Error("Doesn't equal expected value")
					return
				}
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	err := tested.ConnectSubscriptions("cus_uuid", []Subscription{
		{
			ExternalID:     "ext_id1",
			DataSourceUUID: "ds_uuid1",
		},
		{
			ExternalID:     "ext_id2",
			DataSourceUUID: "ds_uuid2",
		},
	})
	if err != nil {
		spew.Dump(err)
		t.Fatal("Expected to retry")
	}
}

func TestCancelSubscriptionParams(t *testing.T) {
	emptySlice := []string{}
	notEmptySlice := []string{"some-date"}
	testCases := map[string]struct {
		param *CancelSubscriptionParams
		exp   string
	}{
		"clearing cancellation history": {
			param: &CancelSubscriptionParams{
				CancellationDates: &emptySlice,
			},
			exp: `{"cancellation_dates":[]}`,
		},
		"setting cancellation history": {
			param: &CancelSubscriptionParams{
				CancellationDates: &notEmptySlice,
			},
			exp: `{"cancellation_dates":["some-date"]}`,
		},
		"not using cancellation history": {
			param: &CancelSubscriptionParams{
				CancelledAt: "some-date",
			},
			exp: `{"cancelled_at":"some-date"}`,
		},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got, err := json.Marshal(tc.param)
			if err != nil {
				spew.Dump(err)
				t.Error("Expected not error")
			}
			if string(got) != tc.exp {
				spew.Dump(tc.exp, string(got))
				t.Error("Doesn't equal expected value")
			}
		})
	}
}
