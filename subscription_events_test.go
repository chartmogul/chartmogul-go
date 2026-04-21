package chartmogul

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
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

	result, err := api.ListSubscriptionEvents(&FilterSubscriptionEvents{DataSourceUUID: testSetup.DataSourceUUID}, &Cursor{})
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

	result, err := api.ListSubscriptionEvents(&FilterSubscriptionEvents{ExternalID: newSubEvent.ExternalID}, &Cursor{})
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

func TestToggleSubscriptionEventDisabled(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "PATCH"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/subscription_events/12345/disabled_state"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				defer r.Body.Close()
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatal("error should be nil")
				}
				expectedBody := `{"disabled":true}`
				if string(body) != expectedBody {
					t.Errorf("Requested body expected: %v, actual: %v", expectedBody, string(body))
				}
				w.Write([]byte(`{"id":12345,"disabled":true,"disabled_at":"2026-03-17T10:00:00Z"}`)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	result, err := tested.ToggleSubscriptionEventDisabled(12345, &ToggleSubscriptionEventDisabledParams{
		Disabled: true,
	})

	if err != nil {
		t.Fatalf("Not expected to fail: %v", err)
	}
	if result.ID != 12345 {
		t.Errorf("Expected ID 12345, got: %v", result.ID)
	}
	if result.Disabled == nil || *result.Disabled != true {
		t.Error("Expected disabled to be true")
	}
	if result.DisabledAt != "2026-03-17T10:00:00Z" {
		t.Errorf("Expected disabled_at, got: %v", result.DisabledAt)
	}
}

func TestToggleSubscriptionEventDisabledByExternalID(t *testing.T) {
	server := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				expectedMethod := "PATCH"
				if r.Method != expectedMethod {
					t.Errorf("Requested method expected: %v, actual: %v", expectedMethod, r.Method)
				}
				expected := "/v/subscription_events/disabled_state"
				path := r.URL.Path
				if path != expected {
					t.Errorf("Requested path expected: %v, actual: %v", expected, path)
					w.WriteHeader(http.StatusNotFound)
				}
				defer r.Body.Close()
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatal("error should be nil")
				}
				// Verify the body contains the expected structure
				if len(body) == 0 {
					t.Fatal("Expected non-empty body")
				}
				w.Write([]byte(`{"id":67890,"external_id":"evt_ext_1","data_source_uuid":"ds_abc","disabled":false}`)) //nolint
			}))
	defer server.Close()
	SetURL(server.URL + "/v/%v")

	var tested IApi = &API{
		ApiKey: "token",
	}
	params := &ToggleSubscriptionEventDisabledByExternalIDParams{
		Disabled: false,
	}
	params.SubscriptionEvent.DataSourceUUID = "ds_abc"
	params.SubscriptionEvent.ExternalID = "evt_ext_1"

	result, err := tested.ToggleSubscriptionEventDisabledByExternalID(params)

	if err != nil {
		t.Fatalf("Not expected to fail: %v", err)
	}
	if result.ID != 67890 {
		t.Errorf("Expected ID 67890, got: %v", result.ID)
	}
	if result.Disabled == nil || *result.Disabled != false {
		t.Error("Expected disabled to be false")
	}
}
