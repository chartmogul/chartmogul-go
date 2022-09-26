package chartmogul

import (
	"log"
	"testing"
)

type TestSetup struct {
	CustomerExternalID          string
	CustomerUUID                string
	DataSourceUUID              string
	InvoiceUUID                 string
	PlanUUID                    string
	SubscriptionEventExternalID string
	SubscriptionEventID         uint64
}

func setup() (*TestSetup, error) {
	dataSource, err := api.CreateDataSource("TestSubEventsDS01")
	if err != nil {
		log.Fatal(err)
		return &TestSetup{}, err
	}

	planDefinition := &Plan{
		DataSourceUUID: dataSource.UUID,
		ExternalID:     "TestSubEventsPlan01",
		Name:           "Super plan",
		IntervalCount:  1,
		IntervalUnit:   "month",
	}
	plan, err := api.CreatePlan(planDefinition)
	if err != nil {
		log.Fatal(err)
		return &TestSetup{}, err
	}

	customer, err := api.CreateCustomer(&NewCustomer{
		DataSourceUUID: dataSource.UUID,
		ExternalID:     "TestSubEventsCustomer01",
		Name:           "Test customer",
	})
	if err != nil {
		log.Fatal(err)
		return &TestSetup{}, err
	}

	invoicesDefinition := []*Invoice{
		{
			CustomerUUID:       customer.UUID,
			CustomerExternalID: customer.ExternalID,
			DataSourceUUID:     dataSource.UUID,
			Currency:           "USD",
			Date:               "2022-02-01",
			DueDate:            "2022-02-28",
			ExternalID:         "TestSubEventsInvoice01",
			LineItems: []*LineItem{{
				AmountInCents:          1000,
				ExternalID:             "TestSubEventsLineItem01",
				Quantity:               1,
				SubscriptionExternalID: "TestSubEventsSubscription01",
				Type:                   "subscription",
				ServicePeriodEnd:       "2022-02-28",
				ServicePeriodStart:     "2022-02-01",
				PlanUUID:               plan.UUID,
			}},
			Transactions: []*Transaction{{
				ExternalID: "TestSubEventsTransaction01",
				Type:       "payment",
				Result:     "successful",
				Date:       "2022-02-01",
			}},
		},
	}
	invoices, err := api.CreateInvoices(invoicesDefinition, customer.UUID)
	if err != nil {
		log.Fatal(err)
		return &TestSetup{}, err
	}
	subEventDefinition := &SubscriptionEvent{
		DataSourceUUID:         dataSource.UUID,
		CustomerExternalID:     customer.ExternalID,
		SubscriptionExternalID: "TestSubEventsSubscription01",
		EventDate:              "2022-02-15",
		EffectiveDate:          "2022-02-27",
		EventType:              "subscription_cancelled",
		ExternalID:             "TestSubEventsSubEvent01",
		Currency:               "PLN",
	}

	subscriptionEvent, err := api.CreateSubscriptionEvent(subEventDefinition)
	if err != nil {
		log.Fatal(err)
		return &TestSetup{}, err
	}

	return &TestSetup{
		CustomerExternalID:          customer.ExternalID,
		CustomerUUID:                customer.UUID,
		DataSourceUUID:              dataSource.UUID,
		InvoiceUUID:                 invoices.Invoices[0].UUID,
		PlanUUID:                    plan.UUID,
		SubscriptionEventExternalID: subscriptionEvent.ExternalID,
		SubscriptionEventID:         subscriptionEvent.ID,
	}, nil
}

func tearDown(testSetup *TestSetup) {
	api.DeleteSubscriptionEvent(&DeleteSubscriptionEvent{ID: testSetup.SubscriptionEventID})
	api.DeleteInvoice(testSetup.InvoiceUUID)
	api.DeleteCustomer(testSetup.CustomerUUID)
	api.DeletePlan(testSetup.PlanUUID)
	api.DeleteDataSource(testSetup.DataSourceUUID)
}

func TestListSubscriptionEvent(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}

	testSetup, err := setup()
	if err != nil {
		t.Error(err)
		return
	}
	defer tearDown(testSetup)

	result, err := api.ListSubscriptionEvents(&FilterSubscriptionEvents{DataSourceUUID: testSetup.DataSourceUUID}, &MetaCursor{})
	if err != nil {
		t.Error(err)
		return
	}

	if len(result.SubscriptionEvents) != 1 ||
		result.SubscriptionEvents[0].ID != testSetup.SubscriptionEventID ||
		result.SubscriptionEvents[0].ExternalID != testSetup.SubscriptionEventExternalID ||
		result.SubscriptionEvents[0].DataSourceUUID != testSetup.DataSourceUUID {
		t.Fatal("Unexpected result")
	}
}

func TestFilteredListSubscriptionEvent(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}

	testSetup, err := setup()
	if err != nil {
		t.Error(err)
		return
	}
	defer tearDown(testSetup)

	newSubEvent, err := api.CreateSubscriptionEvent(&SubscriptionEvent{
		DataSourceUUID:         testSetup.DataSourceUUID,
		CustomerExternalID:     testSetup.CustomerExternalID,
		SubscriptionExternalID: "TestSubEventsSubscription01",
		EventDate:              "2022-02-16",
		EffectiveDate:          "2022-02-28",
		EventType:              "subscription_cancelled",
		ExternalID:             "TestSubEventsSubEvent02",
		Currency:               "PLN",
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer api.DeleteSubscriptionEvent(&DeleteSubscriptionEvent{ID: newSubEvent.ID})

	result, err := api.ListSubscriptionEvents(&FilterSubscriptionEvents{ExternalID: newSubEvent.ExternalID}, &MetaCursor{})
	if err != nil {
		t.Error(err)
		return
	}

	if len(result.SubscriptionEvents) != 1 ||
		result.SubscriptionEvents[0].ID != newSubEvent.ID ||
		result.SubscriptionEvents[0].ExternalID != newSubEvent.ExternalID ||
		result.SubscriptionEvents[0].DataSourceUUID != testSetup.DataSourceUUID {
		t.Fatal("Unexpected result")
	}
}

func TestDeleteSubscriptionEventById(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}

	testSetup, err := setup()
	if err != nil {
		t.Error(err)
		return
	}
	defer tearDown(testSetup)

	err = api.DeleteSubscriptionEvent((&DeleteSubscriptionEvent{ID: testSetup.SubscriptionEventID}))
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestDeleteSubscriptionEventByExternalIdAndDataSourceUuid(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}

	testSetup, err := setup()
	if err != nil {
		t.Error(err)
		return
	}
	defer tearDown(testSetup)

	err = api.DeleteSubscriptionEvent((&DeleteSubscriptionEvent{DataSourceUUID: testSetup.DataSourceUUID, ExternalID: testSetup.SubscriptionEventExternalID}))
	if err != nil {
		t.Fatal(err)
		return
	}
}

func TestUpdateSubscriptionEventUsingId(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}

	testSetup, err := setup()
	if err != nil {
		t.Error(err)
		return
	}
	defer tearDown(testSetup)

	updateDefinition := &SubscriptionEvent{
		ID:       testSetup.SubscriptionEventID,
		Currency: "USD",
	}

	updatedSubEvent, err := api.UpdateSubscriptionEvent(updateDefinition)
	if err != nil {
		t.Error(err)
		return
	}

	if updatedSubEvent.Currency != "USD" {
		t.Errorf("Subscription Event's currency was not updated - expected: %v, actual: %v", "USD", updatedSubEvent.Currency)
	}
}

func TestUpdateSubscriptionEventUsingExternalIdAndDataSourceUuid(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}

	testSetup, err := setup()
	if err != nil {
		t.Error(err)
		return
	}

	defer tearDown(testSetup)
	updateDefinition := &SubscriptionEvent{
		ExternalID:     testSetup.SubscriptionEventExternalID,
		DataSourceUUID: testSetup.DataSourceUUID,
		Currency:       "CNY",
	}

	updatedSubEvent, err := api.UpdateSubscriptionEvent(updateDefinition)
	if err != nil {
		t.Error(err)
		return
	}

	if updatedSubEvent.Currency != "CNY" {
		t.Errorf("Subscription Event's currency was not updated - expected: %v, actual: %v", "CNY", updatedSubEvent.Currency)
	}
}
