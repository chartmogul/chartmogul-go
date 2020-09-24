package chartmogul

import (
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestUploadFile tests creation & listing of customers + duplicate error.
func TestUploadFile(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Unexpected method %v", r.Method)
				}
				if r.RequestURI != "/v/data_sources/ds_uuid/uploads" {
					t.Errorf("Unexpected URI %v", r.RequestURI)
				}
				w.WriteHeader(http.StatusOK)
				//nolint
				w.Write([]byte(`{"id":"ecf1012d-c1b7-4e11-9134-a4f19896df37","original_name":"customers.csv","data_type":"customer","storage_path":"68/2e9ec2d3-3024-4694-9eac-2e399c1306b0/ecf1012d-c1b7-4e11-9134-a4f19896df37","percent_complete":0,"created_at":"2020-08-31T21:44:07Z","updated_at":"2020-08-31T21:44:07Z","batch_name":"Customers"}`))
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	tested := &API{
		AccountToken: "token",
		AccessKey:    "key",
		Client:       &http.Client{},
	}

	file := "./integration_tests/fixtures/invoices.csv"

	_, err := tested.UploadCSVFile(file, &CsvUploadRequest{
		DataSourceUUID: "ds_uuid",
		Type:           "invoice",
		BatchName:      "Invoices Upload",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}

	reader := strings.NewReader("Name,Email,Company,Country,State,City,Zip,External ID,Lead created at,Free trial started at\nJoe Doe,payment1@example.com,,US,IL,Chicago,60611,I-TVNB6F1J8NA6,,")

	_, err = tested.UploadCSVFile(reader, &CsvUploadRequest{
		DataSourceUUID: "ds_uuid",
		Type:           "customer",
		BatchName:      "Invoices Upload",
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}
