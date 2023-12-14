package chartmogul

const subscriptionEventsEndpoint = "subscription_events"

type SubscriptionEvent struct {
	ID                        uint64      `json:"id,omitempty"`
	DataSourceUUID            string      `json:"data_source_uuid,omitempty"`
	CustomerExternalID        string      `json:"customer_external_id,omitempty"`
	SubscriptionSetExternalID string      `json:"subscription_set_external_id,omitempty"`
	SubscriptionExternalID    string      `json:"subscription_external_id,omitempty"`
	PlanExternalID            string      `json:"plan_external_id,omitempty"`
	EventDate                 string      `json:"event_date,omitempty"`
	EffectiveDate             string      `json:"effective_date,omitempty"`
	EventType                 string      `json:"event_type,omitempty"`
	ExternalID                string      `json:"external_id,omitempty"`
	Errors                    interface{} `json:"errors,omitempty"`
	CreatedAt                 string      `json:"created_at,omitempty"`
	UpdatedAt                 string      `json:"updated_at,omitempty"`
	Quantity                  int32       `json:"quantity,omitempty"`
	Currency                  string      `json:"currency,omitempty"`
	AmountInCents             int32       `json:"amount_in_cents,omitempty"`
	TaxAmountInCents          int32       `json:"tax_amount_in_cents,omitempty"`
	RetractedEventId          string      `json:"retracted_event_id,omitempty"`
}

type SubscriptionEvents struct {
	SubscriptionEvents []*SubscriptionEvent `json:"subscription_events"`
	Pagination
}

type DeleteSubscriptionEvent struct {
	ID             uint64 `json:"id,omitempty"`
	DataSourceUUID string `json:"data_source_uuid,omitempty"`
	ExternalID     string `json:"external_id,omitempty"`
}

type FilterSubscriptionEvents struct {
	CustomerExternalID     string `json:"customer_external_id,omitempty"`
	DataSourceUUID         string `json:"data_source_uuid,omitempty"`
	EffectiveDate          string `json:"effective_date,omitempty"`
	EventDate              string `json:"event_date,omitempty"`
	EventType              string `json:"event_type,omitempty"`
	ExternalID             string `json:"external_id,omitempty"`
	PlanExternalID         string `json:"plan_external_id,omitempty"`
	SubscriptionExternalID string `json:"subscription_external_id,omitempty"`
}

type DeleteSubscriptionEventParams struct {
	Params *DeleteSubscriptionEvent `json:"subscription_event"`
}

type SubscriptionEventParams struct {
	Params *SubscriptionEvent `json:"subscription_event"`
}

func (api API) ListSubscriptionEvents(filters *FilterSubscriptionEvents, cursor *Cursor) (*SubscriptionEvents, error) {
	result := &SubscriptionEvents{}
	query := make([]interface{}, 0, 1)
	if cursor != nil {
		query = append(query, *cursor)
	}
	if filters != nil {
		query = append(query, *filters)
	}

	return result, api.list(subscriptionEventsEndpoint, result, query...)
}

func (api API) CreateSubscriptionEvent(newSubscriptionEvent *SubscriptionEvent) (*SubscriptionEvent, error) {
	result := &SubscriptionEvent{}
	return result, api.create(subscriptionEventsEndpoint, SubscriptionEventParams{Params: newSubscriptionEvent}, result)
}

func (api API) UpdateSubscriptionEvent(subscriptionEvent *SubscriptionEvent) (*SubscriptionEvent, error) {
	result := &SubscriptionEvent{}
	return result, api.update(subscriptionEventsEndpoint, "", SubscriptionEventParams{Params: subscriptionEvent}, result)
}

func (api API) DeleteSubscriptionEvent(deleteParams *DeleteSubscriptionEvent) error {
	return api.deleteWithData(
		subscriptionEventsEndpoint,
		DeleteSubscriptionEventParams{Params: deleteParams},
	)
}
