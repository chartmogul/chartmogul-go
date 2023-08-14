package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestListPlans(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("Unexpected method %v", r.Method)
			}
			if r.RequestURI != "/v/plans?data_source_uuid=" {
				t.Errorf("Unexpected URI %v", r.RequestURI)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"plans": [{
					"uuid": "plan_00000000-0000-0000-0000-000000000000",
					"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000",
					"external_id": "plan001"
				}]
			}`))
		}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	plans, err := tested.ListPlans(&ListPlansParams{})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(plans.Plans) == 0 {
		spew.Dump(plans)
		t.Fatal("Unexpected result")
	}
}

func TestRetrievePlan(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				t.Errorf("Unexpected method %v", r.Method)
			}
			if r.RequestURI != "/v/plans/plan_00000000-0000-0000-0000-000000000000" {
				t.Errorf("Unexpected URI %v", r.RequestURI)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"uuid": "plan_00000000-0000-0000-0000-000000000000",
				"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000",
				"external_id": "plan001"
			}`))
		}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	plan, err := tested.RetrievePlan("plan_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if plan == nil || plan.UUID == "" {
		spew.Dump(plan)
		t.Fatal("Unexpected result")
	}
}

func TestCreatePlan(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/plans" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusCreated)
				// Mock response for created plan
				w.Write([]byte(`{
					"uuid": "plan_00000000-0000-0000-0000-000000000000",
					"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000"
				}`))
			}))

	defer server.Close()
	SetURL(server.URL + "/%v")

	tested := &API{
		ApiKey: "token",
	}

	plan, err := tested.CreatePlan(&Plan{
		DataSourceUUID: "ds_00000000-0000-0000-0000-000000000001",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if plan.UUID != "plan_00000000-0000-0000-0000-000000000000" {
		spew.Dump(plan)
		t.Fatal("Unexpected result")
	}
}

func TestUpdatePlan(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PATCH" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/plans/plan_00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				// Mock response for updated plan
				w.Write([]byte(`{
					"uuid": "plan_00000000-0000-0000-0000-000000000000",
					"data_source_uuid": "ds_00000000-0000-0000-0000-000000000000",
					"name": "test"
				}`))
			}))

	defer server.Close()
	SetURL(server.URL + "/%v")

	tested := &API{
		ApiKey: "token",
	}

	updatedPlan, err := tested.UpdatePlan(&Plan{
		DataSourceUUID: "ds_00000000-0000-0000-0000-000000000000",
		Name:           "test",
	},
		"plan_00000000-0000-0000-0000-000000000000",
	)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if updatedPlan.Name != "test" {
		spew.Dump(updatedPlan)
		t.Fatal("Unexpected result")
	}
}

func TestDeletePlan(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/plans/plan_00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusNoContent)
			}))

	defer server.Close()
	SetURL(server.URL + "/%v")

	tested := &API{
		ApiKey: "token",
	}

	err := tested.DeletePlan("plan_00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}
