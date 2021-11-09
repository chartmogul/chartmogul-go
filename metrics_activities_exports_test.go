package chartmogul

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

// Tests creation on an activity export.
func TestCreateActivitiesExport(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/activities_export" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
							      "id": "7f554dba-4a41-4cb2-9790-2045e4c3a5b1",
							      "status": "pending",
							      "file_url": null,
							      "params": {
							        "kind": "activities",
							        "params": {
							          "activity_type": "contraction",
							          "start_date": "2020-01-01",
							          "end_date": "2020-12-31"
							        }
							      },
							      "expires_at": null,
							      "created_at": "2021-07-12T14:46:56+00:00"
							    }`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	activitiesExport, err := tested.MetricsCreateActivitiesExport(&CreateMetricsActivitiesExportParam{
		StartDate: "2020-01-01",
		EndDate:   "2020-12-31",
		Type:      "contraction",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if activitiesExport.ID != "7f554dba-4a41-4cb2-9790-2045e4c3a5b1" || activitiesExport.Status != "pending" {
		spew.Dump(activitiesExport)
		t.Fatal("Unexpected result")
	}
}

// Tests retrieval of an activity export
func TestRetrieveActivitiesExport(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/activities_export/7f554dba-4a41-4cb2-9790-2045e4c3a5b1" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{
							      "id": "7f554dba-4a41-4cb2-9790-2045e4c3a5b1",
							      "status": "succeeded",
							      "file_url": "https://chartmogul-customer-export.s3.eu-west-1.amazonaws.com/activities-acme-corp-91e1ca88-d747-4e25-83d9-2b752033bdba.zip",
							      "params": {
							        "kind": "activities",
							        "params": {
							          "activity_type": "contraction",
							          "start_date": "2020-01-01",
							          "end_date": "2020-12-31"
							        }
							      },
							      "expires_at": "2021-07-19T14:46:58+00:00",
							      "created_at": "2021-07-12T14:46:56+00:00"
							    }`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}

	var activitiesExportID = "7f554dba-4a41-4cb2-9790-2045e4c3a5b1"
	activitiesExport, err := tested.MetricsRetrieveActivitiesExport(activitiesExportID)

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
	if activitiesExport.ID != activitiesExportID || activitiesExport.Status != "succeeded" {
		spew.Dump(activitiesExport)
		t.Fatal("Unexpected result")
	}
}
