package chartmogul

import (
	"flag"
	"os"
	"testing"
)

var (
	cm      = flag.Bool("cm", false, "run integration library tests against ChartMogul")
	api_key = flag.String("api_key", "", "API key for CM test")
	api     = API{}
)

func TestMain(m *testing.M) {
	flag.Parse()
	if *api_key == "" {
		if *cm {
			panic("Please supply testing API key on cmd line to run live tests.")
		}
		*cm = false
	} else {
		api.ApiKey = *api_key
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

var matchersCases = map[string]struct {
	Errors
	expectations map[string]bool
}{
	"invoice & transaction exist": {
		Errors{
			"transactions.external_id": "has already been taken",
			"external_id":              "The external ID for this invoice already exists in our system.",
		}, map[string]bool{
			"IsAlreadyExists":                     false,
			"IsInvoiceAndTransactionAlreadyExist": true,
			"IsInvoiceAndItsEntitiesAlreadyExist": true,
		},
	},
	"invoice exists": {
		Errors{
			"external_id": "The external ID for this invoice already exists in our system.",
		}, map[string]bool{
			"IsAlreadyExists":                     true,
			"IsInvoiceAndTransactionAlreadyExist": false,
			"IsInvoiceAndItsEntitiesAlreadyExist": true,
		},
	},
	"invoice, line items and transactions exist": {
		Errors{
			"external_id":              "The external ID for this invoice already exists in our system.",
			"transactions.external_id": "has already been taken",
			"line_items.external_id":   "The external ID for this line item already exists in our system.",
		}, map[string]bool{
			"IsAlreadyExists":                     false,
			"IsInvoiceAndTransactionAlreadyExist": false,
			"IsInvoiceAndItsEntitiesAlreadyExist": true,
		},
	},
	"transaction exists": {
		Errors{
			"transactions.external_id": "has already been taken",
		}, map[string]bool{
			"IsAlreadyExists":                     true,
			"IsInvoiceAndTransactionAlreadyExist": false,
			"IsInvoiceAndItsEntitiesAlreadyExist": false,
		},
	},
}

func TestErrorMatchers(t *testing.T) {
	for testName, testCase := range matchersCases {
		for fnName, expected := range testCase.expectations {
			t := t
			expected := expected

			var fn func() bool
			switch fnName {
			case "IsAlreadyExists":
				fn = testCase.Errors.IsAlreadyExists
			case "IsInvoiceAndTransactionAlreadyExist":
				fn = testCase.Errors.IsInvoiceAndTransactionAlreadyExist
			case "IsInvoiceAndItsEntitiesAlreadyExist":
				fn = testCase.Errors.IsInvoiceAndItsEntitiesAlreadyExist
			}
			t.Run(testName+"/"+fnName, func(t *testing.T) {
				if fn() != expected {
					t.Error("unexpected match")
				}
			})
		}
	}
}
