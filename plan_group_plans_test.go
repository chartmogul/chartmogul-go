package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

const listAllPlanGroupsPlansExample = `{
	"plans": [
		{
			"name": "A Test Plan",
			"uuid": "pl_2f7f4360-3633-0138-f01f-62b37fb4c770",
			"data_source_uuid":"ds_424b9628-5405-11ea-ab18-e3a0f4f45097",
			"interval_count":1,
			"interval_unit":"month",
			"external_id":"2f7f4360-3633-0138-f01f-62b37fb4c770"
		},
		{
			"name":"Another Test Plan",
			"uuid": "pl_3011ead0-3633-0138-f020-62b37fb4c770",
			"data_source_uuid": "ds_424b9628-5405-11ea-ab18-e3a0f4f45097",
			"interval_count": 1,
			"interval_unit": "month",
			"external_id": "3011ead0-3633-0138-f020-62b37fb4c770"
		}
	],
	"current_page": 1,
	"total_pages": 1
}`

func TestNewPlanGroupPlansAllListing(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(listAllPlanGroupsPlansExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		AccountToken: "token",
		AccessKey:    "key",
	}
	result, err := tested.ListPlanGroupPlans(&Cursor{}, "plg_b53fdbfc-c5eb-4a61-a589-85146cf8d0ab")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(result.Plans) != 2 ||
		result.Plans[0].UUID != "pl_2f7f4360-3633-0138-f01f-62b37fb4c770" ||
		result.Plans[1].UUID != "pl_3011ead0-3633-0138-f020-62b37fb4c770" {
		spew.Dump(result)
		t.Fatal("Unexpected values")
	}
}
