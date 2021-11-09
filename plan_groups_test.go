package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

const onePlanGroupExample = `{
	"name": "My plan group",
	"uuid": "plg_b53fdbfc-c5eb-4a61-a589-85146cf8d0ab",
	"plans_count": 2
}`

const listAllPlanGroupsExample = `{
  "plan_groups": [
    ` + onePlanGroupExample + `
  ],
  "current_page": 1,
  "total_pages": 1
}`

func TestNewPlanGroupAllListing(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(listAllPlanGroupsExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	result, err := tested.ListPlanGroups(&Cursor{})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(result.PlanGroups) != 1 ||
		result.PlanGroups[0].UUID != "plg_b53fdbfc-c5eb-4a61-a589-85146cf8d0ab" ||
		result.PlanGroups[0].PlansCount != 2 {
		spew.Dump(result)
		t.Fatal("Unexpected values")
	}
}

func TestDeletePlanGroup(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expected := "/v/plan_groups/plg_b53fdbfc-c5eb-4a61-a589-85146cf8d0ab"
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
	err := tested.DeletePlanGroup("plg_b53fdbfc-c5eb-4a61-a589-85146cf8d0ab")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}

func TestRetrievePlanGroup(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "GET"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/plan_groups/plg_b53fdbfc-c5eb-4a61-a589-85146cf8d0ab"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				w.Write([]byte(onePlanGroupExample)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	planGroup, err := tested.RetrievePlanGroup("plg_b53fdbfc-c5eb-4a61-a589-85146cf8d0ab")

	if planGroup.Name != "My plan group" || planGroup.PlansCount != 2 || planGroup.UUID != "plg_b53fdbfc-c5eb-4a61-a589-85146cf8d0ab" {
		spew.Dump(planGroup)
		t.Error("Unexpected plan group")
	}
	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}
