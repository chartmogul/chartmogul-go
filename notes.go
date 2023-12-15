package chartmogul

// Note is the customer note as represented in the API.
type Note struct {
	UUID string `json:"uuid"`
	// Basic info
	CustomerUUID string `json:"customer_uuid"`
	Type         string `json:"type"`
	Text         string `json:"text"`
	Author       string `json:"author"`
	CallDuration uint32 `json:"call_duration"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// UpdateNote allows updating note on the update endpoint.
type UpdateNote struct {
	Text         string `json:"text"`
	AuthorEmail  string `json:"author_email"`
	CallDuration uint32 `json:"call_duration"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// NewNote allows creating note on a new endpoint.
type NewNote struct {
	// Obligatory
	CustomerUUID string `json:"customer_uuid"`
	Type         string `json:"type"`

	//Optional
	AuthorEmail  string `json:"author_email"`
	Text         string `json:"text"`
	CallDuration uint32 `json:"call_duration"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// ListNoteParams = parameters for listing customer notes in API.
type ListNotesParams struct {
	CustomerUUID string `json:"customer_uuid"`
	Cursor
}

// Notes is result of listing customer notes in API.
type Notes struct {
	Entries []*Note `json:"entries"`
	Pagination
}

const (
	singleCustomerNoteEndpoint = "customer_notes/:uuid"
	customerNotesEndpoint      = "customer_notes"
)

// CreateCustomerNote loads the customer note to Chartmogul
//
// See https://dev.chartmogul.com/reference/create-a-customer-note
func (api API) CreateNote(input *NewNote) (*Note, error) {
	result := &Note{}
	return result, api.create(customerNotesEndpoint, input, result)
}

// RetrieveCustomerNote returns one customer note as in API.
//
// See https://dev.chartmogul.com/reference/retrieve-a-customer-note
func (api API) RetrieveNote(customerNoteUUID string) (*Note, error) {
	result := &Note{}
	return result, api.retrieve(singleCustomerNoteEndpoint, customerNoteUUID, result)
}

// UpdateNote updates one customer note in API.
//
// See https://dev.chartmogul.com/reference/update-a-customer-note
func (api API) UpdateNote(input *UpdateNote, customerNoteUUID string) (*Note, error) {
	output := &Note{}
	return output, api.update(singleCustomerNoteEndpoint, customerNoteUUID, input, output)
}

// ListNotes lists all Notes
//
// See https://dev.chartmogul.com/reference/list-all-customer-notes
func (api API) ListNotes(listNotesParams *ListNotesParams) (*Notes, error) {
	result := &Notes{}
	query := make([]interface{}, 0, 1)
	if listNotesParams != nil {
		query = append(query, *listNotesParams)
	}
	return result, api.list(customerNotesEndpoint, result, query...)
}

// DeleteNote deletes one customer note by UUID.
//
// See https://dev.chartmogul.com/reference/delete-a-customer-note
func (api API) DeleteNote(customerNoteUUID string) error {
	return api.delete(singleCustomerNoteEndpoint, customerNoteUUID)
}
