package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestListOpportunities(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/opportunities?customer_uuid=cus_00000000-0000-0000-0000-000000000000&per_page=1" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"entries": [{
						"uuid": "00000000-0000-0000-0000-000000000000",
						"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
						"owner": "test1@example.org",
						"pipeline": "New business 1",
						"pipeline_stage": "Discovery",
						"estimated_close_date": "2023-12-22",
						"currency": "USD",
						"amount_in_cents": 100,
						"type": "recurring",
						"forecast_category": "pipeline",
						"win_likelihood": 3,
						"custom": {"from_campaign": true},
						"created_at": "2024-03-13T07:33:28.356Z",
						"updated_at": "2024-03-13T07:33:28.356Z"
					}],
					"has_more": false,
					"cursor": "88abf99"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	params := &ListOpportunitiesParams{Cursor: Cursor{PerPage: 1}, CustomerUUID: "cus_00000000-0000-0000-0000-000000000000"}
	opportunities, err := tested.ListOpportunities(params)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if len(opportunities.Entries) == 0 {
		spew.Dump(opportunities)
		t.Fatal("Unexpected result")
	}
}

func TestRetrieveOpportunity(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/opportunities/00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"uuid": "00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"owner": "test1@example.org",
					"pipeline": "New business 1",
					"pipeline_stage": "Discovery",
					"estimated_close_date": "2023-12-22",
					"currency": "USD",
					"amount_in_cents": 100,
					"type": "recurring",
					"forecast_category": "pipeline",
					"win_likelihood": 3,
					"custom": {"from_campaign": true},
					"created_at": "2024-03-13T07:33:28.356Z",
					"updated_at": "2024-03-13T07:33:28.356Z"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}
	opportunity, err := tested.RetrieveOpportunity("00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if opportunity == nil {
		spew.Dump(opportunity)
		t.Fatal("Unexpected result")
	}
}

func TestCreateOpportunity(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/opportunities" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusCreated)
				//nolint
				w.Write([]byte(`{
					"uuid": "00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"owner": "test1@example.org",
					"pipeline": "New business 1",
					"pipeline_stage": "Discovery",
					"estimated_close_date": "2023-12-22",
					"currency": "USD",
					"amount_in_cents": 100,
					"type": "recurring",
					"forecast_category": "pipeline",
					"win_likelihood": 3,
					"custom": {"from_campaign": true},
					"created_at": "2024-03-13T07:33:28.356Z",
					"updated_at": "2024-03-13T07:33:28.356Z"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	opportunity, err := tested.CreateOpportunity(&NewOpportunity{
		CustomerUUID:       "cus_00000000-0000-0000-0000-000000000000",
		Owner:              "test1@example.org",
		Pipeline:           "New business 1",
		PipelineStage:      "Discovery",
		EstimatedCloseDate: "2023-12-22",
		Currency:           "USD",
		AmountInCents:      100,
		Type:               "recurring",
		ForecastCategory:   "pipeline",
		WinLikelihood:      3,
		Custom: []Custom{
			{
				Key:   "from_campaigng",
				Value: true,
			},
		},
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if opportunity.UUID != "00000000-0000-0000-0000-000000000000" {
		spew.Dump(opportunity)
		t.Fatal("Unexpected result")
	}
}

func TestUpdateOpportunity(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "PATCH" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/opportunities/00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
					"uuid": "00000000-0000-0000-0000-000000000000",
					"customer_uuid": "cus_00000000-0000-0000-0000-000000000000",
					"owner": "test1@example.org",
					"pipeline": "New business 1",
					"pipeline_stage": "Discovery",
					"estimated_close_date": "2024-12-22",
					"currency": "USD",
					"amount_in_cents": 100,
					"type": "recurring",
					"forecast_category": "pipeline",
					"win_likelihood": 3,
					"custom": {"from_campaign": true},
					"created_at": "2024-03-13T07:33:28.356Z",
					"updated_at": "2024-03-13T07:33:28.356Z"
				}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	opportunity, err := tested.UpdateOpportunity(&UpdateOpportunity{
		EstimatedCloseDate: "2024-12-22",
	}, "00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if opportunity.EstimatedCloseDate != "2024-12-22" {
		spew.Dump(opportunity)
		t.Fatal("Unexpected result")
	}
}

func TestDeleteOpportunity(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/opportunities/00000000-0000-0000-0000-000000000000" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		ApiKey: "token",
	}

	err := tested.DeleteOpportunity("00000000-0000-0000-0000-000000000000")

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}
