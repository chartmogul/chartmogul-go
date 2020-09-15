package chartmogul

import (
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
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
				if r.RequestURI != "/v/../data_platform/v1/data_sources/:data_source_uuid/uploads" {
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
	}

	f, _ := filepath.Abs("./invoices.csv")
	csvContent, _ := ioutil.ReadFile(f)

	_, err := tested.UploadCSVFile(&CsvUploadRequest{
		DataSourceUUID: "ds_uuid",
		DataType:       "invoice",
		BatchName:      "Invoices Upload",
		File:           csvContent,
	})

	if err != nil {
		spew.Dump(err)
		t.Fatal("Not expected to fail")
	}
}
