package integration

import (
	"log"
	"net/http"
	"os"
	"testing"

	cm "github.com/chartmogul/chartmogul-go/v4"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-test/deep"
	"github.com/parnurzeal/gorequest"
)

func TestRetrieveMetrics(t *testing.T) {
	if testing.Short() {
		t.Skip("Integration test.")
	}

	r, err := NewRecorderWithAuthFilter("./fixtures/retrieve_metrics")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Stop() //nolint

	api := &cm.API{
		ApiKey: os.Getenv("CHARTMOGUL_API_KEY"),
		Client: &http.Client{Transport: r},
	}
	gorequest.DisableTransportSwap = true

	all := &cm.MetricsResult{
		Entries: []*cm.AllMetrics{
			{
				Date:                              "2022-04-30",
				Mrr:                               68604881,
				MrrPercentageChange:               0.0,
				Arr:                               823258572,
				ArrPercentageChange:               0.0,
				CustomerChurnRate:                 2.85,
				CustomerChurnRatePercentageChange: 0.0,
				MrrChurnRate:                      -0.67,
				MrrChurnRatePercentageChange:      0.0,
				Ltv:                               2194315,
				LtvPercentageChange:               0.0,
				Customers:                         1097,
				CustomersPercentageChange:         0.0,
				Asp:                               22776,
				AspPercentageChange:               0.0,
				Arpa:                              62538,
				ArpaPercentageChange:              0.0,
			},
			{
				Date:                              "2022-05-31",
				Mrr:                               68954549,
				MrrPercentageChange:               0.51,
				Arr:                               827454588,
				ArrPercentageChange:               0.51,
				CustomerChurnRate:                 2.92,
				CustomerChurnRatePercentageChange: 2.46,
				MrrChurnRate:                      0.91,
				MrrChurnRatePercentageChange:      235.82,
				Ltv:                               2142876,
				LtvPercentageChange:               -2.34,
				Customers:                         1102,
				CustomersPercentageChange:         0.46,
				Asp:                               30531,
				AspPercentageChange:               34.05,
				Arpa:                              62572,
				ArpaPercentageChange:              0.05,
			},
		},
		Summary: &cm.AllSummary{
			CurrentMrr:                        69499091,
			PreviousMrr:                       68998840,
			MrrPercentageChange:               0.73,
			CurrentArr:                        833989092,
			PreviousArr:                       827986080,
			ArrPercentageChange:               0.73,
			CurrentCustomerChurnRate:          2.92,
			PreviousCustomerChurnRate:         2.85,
			CustomerChurnRatePercentageChange: 2.46,
			CurrentMrrChurnRate:               0.91,
			PreviousMrrChurnRate:              -0.67,
			MrrChurnRatePercentageChange:      235.82,
			CurrentLtv:                        2142876,
			PreviousLtv:                       2194315,
			LtvPercentageChange:               -2.34,
			CurrentCustomers:                  1100,
			PreviousCustomers:                 1099,
			CustomersPercentageChange:         0.09,
			CurrentAsp:                        31486,
			PreviousAsp:                       19737,
			AspPercentageChange:               59.53,
			CurrentArpa:                       63180,
			PreviousArpa:                      62783,
			ArpaPercentageChange:              0.63,
		},
	}
	retrieved_all, err := api.MetricsRetrieveAll(&cm.MetricsFilter{
		StartDate: "2022-04-01",
		EndDate:   "2022-05-31",
		Interval:  "month",
	})
	if err != nil {
		t.Fatal(err)
	}
	diff_all := deep.Equal(retrieved_all, all)
	if diff_all != nil {
		spew.Dump(all)
		t.Errorf("compare failed: %#v", diff_all)
	}

	mrr := &cm.MRRResult{
		Entries: []*cm.MRRMetrics{
			{
				Date:             "2022-04-30",
				MRR:              68604881,
				PercentageChange: 0.0,
				MRRNewBusiness:   865500,
				MRRExpansion:     2386603,
				MRRContraction:   -508125,
				MRRChurn:         -1660000,
				MRRReactivation:  230000,
			},
			{
				Date:             "2022-05-31",
				MRR:              68954549,
				PercentageChange: 0.51,
				MRRNewBusiness:   977000,
				MRRExpansion:     1559209,
				MRRContraction:   -770875,
				MRRChurn:         -1730666,
				MRRReactivation:  315000,
			},
		},
		Summary: &cm.Summary{
			Current:          69499091,
			Previous:         68998840,
			PercentageChange: 0.73,
		},
	}
	retrieved_mrr, err := api.MetricsRetrieveMRR(&cm.MetricsFilter{
		StartDate: "2022-04-01",
		EndDate:   "2022-05-31",
		Interval:  "month",
	})
	if err != nil {
		t.Fatal(err)
	}
	diff_mrr := deep.Equal(retrieved_mrr, mrr)
	if diff_mrr != nil {
		spew.Dump(all)
		t.Errorf("compare failed: %#v", diff_mrr)
	}
}
