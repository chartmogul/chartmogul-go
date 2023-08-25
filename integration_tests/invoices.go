package integration

import (
	"testing"
	"time"

	cm "github.com/chartmogul/chartmogul-go/v4"
)

// createTestInvoicesForCustomer creates a set of test invoices for a specified plan and customer.
// It returns the created invoices.
//
//nolint:gomnd
func createTestInvoicesForCustomer(api *cm.API, cus cm.Customer, planUUID string, t *testing.T) *cm.Invoices {
	testInvoices := []*cm.Invoice{
		{
			Date:               time.Now().Format(time.RFC3339),
			ExternalID:         "INV_to_be_retrieved",
			CustomerUUID:       cus.UUID,
			CustomerExternalID: cus.ExternalID,
			Currency:           "EUR",
			LineItems: []*cm.LineItem{
				{
					Type:          "one_time",
					AmountInCents: 4500,
					Description:   "fake_item",
					Quantity:      2,
				},
				{
					Type:                      "subscription",
					AmountInCents:             10000,
					ExternalID:                "ext_line_item",
					SubscriptionExternalID:    "ext_subscription",
					SubscriptionSetExternalID: "ext_subscription_set",
					PlanUUID:                  planUUID,
					Quantity:                  10,
					ServicePeriodStart:        "2017-05-01T00:00:00.000Z",
					ServicePeriodEnd:          "2017-05-31T00:00:00.000Z",
				},
			},
			Transactions: []*cm.Transaction{
				{
					Date:   time.Now().Add(3 * time.Hour).Format(time.RFC3339),
					Result: "successful",
					Type:   "payment",
				},
			},
		},
		{
			Date:               time.Now().Format(time.RFC3339),
			ExternalID:         "INV_to_be_retrieved1",
			CustomerUUID:       cus.UUID,
			CustomerExternalID: cus.ExternalID,
			Currency:           "EUR",
			LineItems: []*cm.LineItem{
				{
					Type:          "one_time",
					AmountInCents: 4200,
					Description:   "fake_item1",
					Quantity:      2,
				},
				{
					Type:                      "subscription",
					AmountInCents:             11000,
					ExternalID:                "ext_line_item1",
					SubscriptionExternalID:    "ext_subscription1",
					SubscriptionSetExternalID: "ext_subscription_set1",
					PlanUUID:                  planUUID,
					Quantity:                  11,
					ServicePeriodStart:        "2017-05-01T00:00:00.000Z",
					ServicePeriodEnd:          "2017-05-31T00:00:00.000Z",
				},
			},
			Transactions: []*cm.Transaction{
				{
					Date:   time.Now().Add(3 * time.Hour).Format(time.RFC3339),
					Result: "successful",
					Type:   "payment",
				},
			},
		},
	}

	inv, err := api.CreateInvoices(testInvoices, cus.UUID)
	if err != nil {
		t.Fatal(err)
	}

	return inv
}
