package chartmogul

import (
	"flag"
	"os"
	"testing"
)

var (
	cm    = flag.Bool("cm", false, "run integration library tests against ChartMogul")
	token = flag.String("token", "", "account token for CM test")
	key   = flag.String("key", "", "access key for CM test")
	api   = API{}
)

func TestMain(m *testing.M) {
	flag.Parse()
	if *key == "" || *token == "" {
		if *cm {
			panic("Please supply testing account key and token on cmd line to run live tests.")
		}
		*cm = false
	} else {
		api.AccountToken = *token
		api.AccessKey = *key
	}

	result := m.Run()

	os.Exit(result)
}

func TestPing(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}
	b, err := api.Ping()
	if err != nil {
		t.Error(err)
	} else if !b {
		t.Error("ping returned false")
	}
}

func TestIsInvoiceAndTransactionAlreadyExist(t *testing.T) {
	err := Errors(map[string]string{
		"transactions.external_id": "has already been taken",
		"external_id":              "The external ID for this invoice already exists in our system.",
	})
	if !err.IsInvoiceAndTransactionAlreadyExist() {
		t.Error("expected IsInvoiceAndTransactionAlreadyExist to be true")
	}
}

//TODO: unit tests against mocked HTTP server.
