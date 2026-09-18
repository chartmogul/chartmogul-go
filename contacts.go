package chartmogul

import "strings"

// Contact is the contact as represented in the API.
type Contact struct {
	UUID string `json:"uuid,omitempty"`
	// Basic info
	CustomerExternalID string `json:"customer_external_id,omitempty"`
	CustomerUUID       string `json:"customer_uuid,omitempty"`
	DataSourceUUID     string `json:"data_source_uuid,omitempty"`
	Email              string `json:"email,omitempty"`
	FirstName          string `json:"first_name,omitempty"`
	LastName           string `json:"last_name,omitempty"`
	LastSeen           string `json:"last_seen,omitempty"`
	LinkedIn           string `json:"linked_in,omitempty"`
	Notes              string `json:"notes,omitempty"`
	Phone              string `json:"phone,omitempty"`
	Position           uint32 `json:"position,omitempty"`
	Title              string `json:"title,omitempty"`
	Twitter            string `json:"twitter,omitempty"`
	// Using *string allows callers to explicitly clear this field
	// Passing nil means omit the field; passing "" means clear the field
	ExternalID *string                `json:"external_id,omitempty"`
	Custom     map[string]interface{} `json:"custom,omitempty"`
}

// UpdateContact allows updating contact on the update endpoint.
type UpdateContact struct {
	CustomerExternalID string `json:"customer_external_id,omitempty"`
	DataSourceUUID     string `json:"data_source_uuid,omitempty"`
	Email              string `json:"email,omitempty"`
	FirstName          string `json:"first_name,omitempty"`
	LastName           string `json:"last_name,omitempty"`
	LastSeen           string `json:"last_seen,omitempty"` // ISO 8601
	LinkedIn           string `json:"linked_in,omitempty"`
	Notes              string `json:"notes,omitempty"`
	Phone              string `json:"phone,omitempty"`
	Position           uint32 `json:"position,omitempty"`
	Title              string `json:"title,omitempty"`
	Twitter            string `json:"twitter,omitempty"`
	// Using *string allows callers to explicitly clear this field
	// Passing nil means omit the field; passing "" means clear the field
	ExternalID *string  `json:"external_id,omitempty"`
	Custom     []Custom `json:"custom,omitempty"`
}

// NewContact allows creating contact on a new endpoint.
type NewContact struct {
	// Optional: omit both to create a contact that is not linked to a customer.
	CustomerUUID   string `json:"customer_uuid,omitempty"`
	DataSourceUUID string `json:"data_source_uuid,omitempty"`

	// Optional
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	LastSeen  string `json:"last_seen,omitempty"` // ISO 8601
	LinkedIn  string `json:"linked_in,omitempty"`
	Notes     string `json:"notes,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Position  uint32 `json:"position,omitempty"`
	Title     string `json:"title,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	// Using *string allows callers to explicitly clear this field
	// Passing nil means omit the field; passing "" means clear the field
	ExternalID *string  `json:"external_id,omitempty"`
	Custom     []Custom `json:"custom,omitempty"`
}

// ListContactsParams = parameters for listing contacts in API.
type ListContactsParams struct {
	CustomerUUID       string `json:"customer_uuid,omitempty"`
	DataSourceUUID     string `json:"data_source_uuid,omitempty"`
	Email              string `json:"email,omitempty"`
	CustomerExternalID string `json:"customer_external_id,omitempty"`
	ExternalID         string `json:"external_id,omitempty"`
	Cursor
}

// Contacts is result of listing contacts in API.
type Contacts struct {
	Entries []*Contact `json:"entries,omitempty"`
	Pagination
}

const (
	singleContactEndpoint = "contacts/:uuid"
	contactsEndpoint      = "contacts"
	mergeContactsEndpoint = "contacts/:into_contact_uuid/merge/:from_contact_uuid"
)

// CreateContact loads the contact to Chartmogul
func (api API) CreateContact(newContact *NewContact) (*Contact, error) {
	result := &Contact{}
	return result, api.create(contactsEndpoint, newContact, result)
}

// RetrieveContact returns one contact as in API.
func (api API) RetrieveContact(contactUUID string) (*Contact, error) {
	result := &Contact{}
	return result, api.retrieve(singleContactEndpoint, contactUUID, result)
}

// UpdateContact updates one contact in API.
func (api API) UpdateContact(input *UpdateContact, contactUUID string) (*Contact, error) {
	output := &Contact{}
	return output, api.update(singleContactEndpoint, contactUUID, input, output)
}

// ListContacts lists all Contacts
func (api API) ListContacts(listContactsParams *ListContactsParams) (*Contacts, error) {
	result := &Contacts{}
	query := make([]interface{}, 0, 1)
	if listContactsParams != nil {
		query = append(query, *listContactsParams)
	}
	return result, api.list(contactsEndpoint, result, query...)
}

// MergeContacts merges two contacts.
func (api API) MergeContacts(intoContactUUID string, fromContactUUID string) (*Contact, error) {
	result := &Contact{}
	temp_path := strings.Replace(mergeContactsEndpoint, ":into_contact_uuid", intoContactUUID, 1)
	path := strings.Replace(temp_path, ":from_contact_uuid", fromContactUUID, 1)
	return result, api.create(path, nil, result)
}

// DeleteContact deletes one contact by UUID.
func (api API) DeleteContact(contactUUID string) error {
	return api.delete(singleContactEndpoint, contactUUID)
}

// ListContactTasks lists the tasks attached to the contact.
func (api API) ListContactTasks(listTasksParams *ListTasksParams, contactUUID string) (*Tasks, error) {
	if listTasksParams == nil {
		listTasksParams = &ListTasksParams{}
	}
	if listTasksParams.ContactUUID == "" {
		listTasksParams.ContactUUID = contactUUID
	}
	return api.ListTasks(listTasksParams)
}

// CreateContactTask creates a task attached to the contact.
func (api API) CreateContactTask(input *NewTask, contactUUID string) (*Task, error) {
	if input.CustomerUUID == "" && input.AssociatedObjectIdentifier == nil {
		input.AssociatedObjectIdentifier = contactIdentifier(contactUUID)
	}
	return api.CreateTask(input)
}

// ListContactEntityNotes lists the notes attached to the contact.
func (api API) ListContactEntityNotes(listEntityNotesParams *ListEntityNotesParams, contactUUID string) (*EntityNotes, error) {
	if listEntityNotesParams == nil {
		listEntityNotesParams = &ListEntityNotesParams{}
	}
	if listEntityNotesParams.ContactUUID == "" {
		listEntityNotesParams.ContactUUID = contactUUID
	}
	return api.ListEntityNotes(listEntityNotesParams)
}

// CreateContactEntityNote creates a note attached to the contact.
func (api API) CreateContactEntityNote(input *NewEntityNote, contactUUID string) (*EntityNote, error) {
	if input.CustomerUUID == "" && input.AssociatedObjectIdentifier == nil {
		input.AssociatedObjectIdentifier = contactIdentifier(contactUUID)
	}
	return api.CreateEntityNote(input)
}

func contactIdentifier(contactUUID string) *AssociatedObjectIdentifier {
	return &AssociatedObjectIdentifier{
		AssociatedObject: AssociatedObjectContact,
		Method:           AssociatedObjectIdentifierMethodUUID,
		Value:            contactUUID,
	}
}
