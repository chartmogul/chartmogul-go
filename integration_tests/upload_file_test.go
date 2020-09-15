package integration

import (
	cm "github.com/chartmogul/chartmogul-go"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
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

	api := &cm.API{AccountToken: "5d490f3d8bf161bdc914fff6708bdc6b", AccessKey: "cb26b59dae523d0d62adfc143f24866e"}
	api.SetClient(&http.Client{Transport: r})
	gorequest.DisableTransportSwap = true

	/*ds, err := api.CreateDataSource("Test upload file")
	if err != nil {
		t.Fatal(err)
	}
	defer api.DeleteDataSource(ds.UUID) //nolint
	*/
	//f, _ := filepath.Abs("./fixtures/invoices.csv")
	//csvContent, _ := ioutil.ReadFile(f)

	uploadJob, err := api.UploadCSVFile("./fixtures/invoices.csv", &cm.CsvUploadRequest{
		DataSourceUUID: "035e73a2-f43e-11ea-9f8b-0ff773344dee",
		DataType:       "invoice",
		BatchName:      "Invoices Upload"})

	if err != nil {
		t.Fatal(err)
	}

	log.Print(uploadJob)
}
