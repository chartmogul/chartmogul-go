package chartmogul

// EntityNote is a note or call log attached to a customer or a contact, as represented in the API.
type EntityNote struct {
	UUID string `json:"uuid"`
	// Basic info
	CustomerUUID         string `json:"customer_uuid,omitempty"`          // empty for notes attached to a contact
	AssociatedObject     string `json:"associated_object,omitempty"`      // "customer" | "contact"
	AssociatedObjectUUID string `json:"associated_object_uuid,omitempty"` // cus_… or con_…
	Type                 string `json:"type"`
	Text                 string `json:"text"`
	Author               string `json:"author"`
	CallDuration         uint32 `json:"call_duration"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
}

// UpdateEntityNote allows updating a note on the update endpoint.
type UpdateEntityNote struct {
	Text         string `json:"text,omitempty"`
	AuthorEmail  string `json:"author_email,omitempty"`
	CallDuration uint32 `json:"call_duration,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

// NewEntityNote allows creating a note on the new endpoint.
type NewEntityNote struct {
	// Exactly one of CustomerUUID or AssociatedObjectIdentifier is required.
	CustomerUUID               string                      `json:"customer_uuid,omitempty"`
	AssociatedObjectIdentifier *AssociatedObjectIdentifier `json:"associated_object_identifier,omitempty"`

	// Obligatory
	Type string `json:"type"` // "note" | "call"

	// Optional
	AuthorEmail  string `json:"author_email,omitempty"`
	Text         string `json:"text,omitempty"`
	CallDuration uint32 `json:"call_duration,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}

// ListEntityNotesParams = parameters for listing notes in API.
type ListEntityNotesParams struct {
	CustomerUUID string `json:"customer_uuid,omitempty"`
	ContactUUID  string `json:"contact_uuid,omitempty"`
	Type         string `json:"type,omitempty"` // "note" | "call"
	AuthorEmail  string `json:"author_email,omitempty"`
	Cursor
}

// EntityNotes is result of listing notes in API.
type EntityNotes struct {
	Entries []*EntityNote `json:"entries"`
	Pagination
}

const (
	singleEntityNoteEndpoint = "notes/:uuid"
	entityNotesEndpoint      = "notes"
)

// CreateEntityNote loads the note to Chartmogul.
func (api API) CreateEntityNote(input *NewEntityNote) (*EntityNote, error) {
	result := &EntityNote{}
	return result, api.create(entityNotesEndpoint, input, result)
}

// RetrieveEntityNote returns one note as in API.
func (api API) RetrieveEntityNote(entityNoteUUID string) (*EntityNote, error) {
	result := &EntityNote{}
	return result, api.retrieve(singleEntityNoteEndpoint, entityNoteUUID, result)
}

// UpdateEntityNote updates one note in API.
// The API responds with 304 and an empty body when the input carries no updatable field,
// which is reported as an error.
func (api API) UpdateEntityNote(input *UpdateEntityNote, entityNoteUUID string) (*EntityNote, error) {
	output := &EntityNote{}
	return output, api.update(singleEntityNoteEndpoint, entityNoteUUID, input, output)
}

// ListEntityNotes lists all notes.
func (api API) ListEntityNotes(listEntityNotesParams *ListEntityNotesParams) (*EntityNotes, error) {
	result := &EntityNotes{}
	query := make([]interface{}, 0, 1)
	if listEntityNotesParams != nil {
		query = append(query, *listEntityNotesParams)
	}
	return result, api.list(entityNotesEndpoint, result, query...)
}

// DeleteEntityNote deletes one note by UUID.
func (api API) DeleteEntityNote(entityNoteUUID string) error {
	return api.delete(singleEntityNoteEndpoint, entityNoteUUID)
}
