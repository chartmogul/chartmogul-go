package integration

import (
	cm "github.com/chartmogul/chartmogul-go"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestUploadCSVFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := recorder.New("./fixtures/upload_csv_file")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() //nolint

	api := &cm.API{}
	api.SetClient(&http.Client{Transport: r})
	gorequest.DisableTransportSwap = true

	ds, err := api.CreateDataSource("Testing CSV batch upload file")
	if err != nil {
		t.Fatal(err)
	}
	defer api.DeleteDataSource(ds.UUID) //nolint

	uploadJob, err := api.UploadCSVFile("./fixtures/invoices.csv", &cm.CsvUploadRequest{
		DataSourceUUID: strings.Replace(ds.UUID, "ds_", "", 1),
		Type:           "invoice",
		BatchName:      "Invoices Upload"})

	if err != nil {
		t.Fatal(err)
	}
	if uploadJob.OriginalName != "invoices.csv" {
		t.Fatal("Wrong Original Name")
	}

	if uploadJob.BatchName != "Invoices Upload" {
		t.Fatal("Wrong batch_name")
	}
}
