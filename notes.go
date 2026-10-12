package chartmogul

import "log"

// Note is the customer note as represented in the API.
//
// Deprecated: Use EntityNote instead.
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
//
// Deprecated: Use UpdateEntityNote instead.
type UpdateNote struct {
	Text         string `json:"text,omitempty"`
	AuthorEmail  string `json:"author_email,omitempty"`
	CallDuration uint32 `json:"call_duration,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

// NewNote allows creating note on a new endpoint.
//
// Deprecated: Use NewEntityNote instead.
type NewNote struct {
	// Obligatory
	CustomerUUID string `json:"customer_uuid"`
	Type         string `json:"type"`

	//Optional
	AuthorEmail  string `json:"author_email,omitempty"`
	Text         string `json:"text,omitempty"`
	CallDuration uint32 `json:"call_duration,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

// ListNotesParams = parameters for listing customer notes in API.
//
// Deprecated: Use ListEntityNotesParams instead.
type ListNotesParams struct {
	CustomerUUID string `json:"customer_uuid,omitempty"`
	Cursor
}

// Notes is result of listing customer notes in API.
//
// Deprecated: Use EntityNotes instead.
type Notes struct {
	Entries []*Note `json:"entries"`
	Pagination
}

const (
	singleCustomerNoteEndpoint = "customer_notes/:uuid"
	customerNotesEndpoint      = "customer_notes"
)

// CreateNote loads the customer note to Chartmogul.
//
// Deprecated: Use CreateEntityNote instead.
func (api API) CreateNote(input *NewNote) (*Note, error) {
	log.Println("[DEPRECATED] CreateNote is deprecated. Use CreateEntityNote instead.")
	result := &Note{}
	return result, api.create(customerNotesEndpoint, input, result)
}

// RetrieveNote returns one customer note as in API.
//
// Deprecated: Use RetrieveEntityNote instead.
func (api API) RetrieveNote(customerNoteUUID string) (*Note, error) {
	log.Println("[DEPRECATED] RetrieveNote is deprecated. Use RetrieveEntityNote instead.")
	result := &Note{}
	return result, api.retrieve(singleCustomerNoteEndpoint, customerNoteUUID, result)
}

// UpdateNote updates one customer note in API.
//
// Deprecated: Use UpdateEntityNote instead.
func (api API) UpdateNote(input *UpdateNote, customerNoteUUID string) (*Note, error) {
	log.Println("[DEPRECATED] UpdateNote is deprecated. Use UpdateEntityNote instead.")
	output := &Note{}
	return output, api.update(singleCustomerNoteEndpoint, customerNoteUUID, input, output)
}

// ListNotes lists all customer notes.
//
// Deprecated: Use ListEntityNotes instead.
func (api API) ListNotes(listNotesParams *ListNotesParams) (*Notes, error) {
	log.Println("[DEPRECATED] ListNotes is deprecated. Use ListEntityNotes instead.")
	result := &Notes{}
	query := make([]interface{}, 0, 1)
	if listNotesParams != nil {
		query = append(query, *listNotesParams)
	}
	return result, api.list(customerNotesEndpoint, result, query...)
}

// DeleteNote deletes one customer note by UUID.
//
// Deprecated: Use DeleteEntityNote instead.
func (api API) DeleteNote(customerNoteUUID string) error {
	log.Println("[DEPRECATED] DeleteNote is deprecated. Use DeleteEntityNote instead.")
	return api.delete(singleCustomerNoteEndpoint, customerNoteUUID)
}
