package chartmogul

import "strings"

// Contact is the contact as represented in the API.
type Contact struct {
	UUID string `json:"uuid,omitempty"`
	// Basic info
	CustomerExternalID string                 `json:"customer_external_id,omitempty"`
	CustomerUUID       string                 `json:"customer_uuid,omitempty"`
	DataSourceUUID     string                 `json:"data_source_uuid,omitempty"`
	FirstName          string                 `json:"first_name,omitempty"`
	LastName           string                 `json:"last_name,omitempty"`
	LinkedIn           string                 `json:"linked_in,omitempty"`
	Notes              string                 `json:"notes,omitempty"`
	Phone              string                 `json:"phone,omitempty"`
	Position           uint32                 `json:"position,omitempty"`
	Title              string                 `json:"title,omitempty"`
	Twitter            string                 `json:"twitter,omitempty"`
	// Using *string allows callers to explicitly clear this field
	// Passing nil means omit the field; passing "" means clear the field
	ExternalID         *string                `json:"external_id,omitempty"`
	Custom             map[string]interface{} `json:"custom,omitempty"`
}

// UpdateContact allows updating contact on the update endpoint.
type UpdateContact struct {
	CustomerExternalID string   `json:"customer_external_id,omitempty"`
	DataSourceUUID     string   `json:"data_source_uuid,omitempty"`
	FirstName          string   `json:"first_name,omitempty"`
	LastName           string   `json:"last_name,omitempty"`
	LinkedIn           string   `json:"linked_in,omitempty"`
	Notes              string   `json:"notes,omitempty"`
	Phone              string   `json:"phone,omitempty"`
	Position           uint32   `json:"position,omitempty"`
	Title              string   `json:"title,omitempty"`
	Twitter            string   `json:"twitter,omitempty"`
	// Using *string allows callers to explicitly clear this field
	// Passing nil means omit the field; passing "" means clear the field
	ExternalID         *string  `json:"external_id,omitempty"`
	Custom             []Custom `json:"custom,omitempty"`
}

// NewContact allows creating contact on a new endpoint.
type NewContact struct {
	// Obligatory
	CustomerUUID   string `json:"customer_uuid,omitempty"`
	DataSourceUUID string `json:"data_source_uuid,omitempty"`

	// Optional
	FirstName  string   `json:"first_name,omitempty"`
	LastName   string   `json:"last_name,omitempty"`
	LinkedIn   string   `json:"linked_in,omitempty"`
	Notes      string   `json:"notes,omitempty"`
	Phone      string   `json:"phone,omitempty"`
	Position   uint32   `json:"position,omitempty"`
	Title      string   `json:"title,omitempty"`
	Twitter    string   `json:"twitter,omitempty"`
	// Using *string allows callers to explicitly clear this field
	// Passing nil means omit the field; passing "" means clear the field
	ExternalID *string  `json:"external_id,omitempty"`
	Custom     []Custom `json:"custom,omitempty"`
}

// ListContactsParams = parameters for listing contacts in API.
type ListContactsParams struct {
	CustomerUUID   string `json:"customer_uuid,omitempty"`
	DataSourceUUID string `json:"data_source_uuid,omitempty"`
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
//
func (api API) CreateContact(newContact *NewContact) (*Contact, error) {
	result := &Contact{}
	return result, api.create(contactsEndpoint, newContact, result)
}

// RetrieveContact returns one contact as in API.
//
func (api API) RetrieveContact(contactUUID string) (*Contact, error) {
	result := &Contact{}
	return result, api.retrieve(singleContactEndpoint, contactUUID, result)
}

// UpdateContact updates one contact in API.
//
func (api API) UpdateContact(input *UpdateContact, contactUUID string) (*Contact, error) {
	output := &Contact{}
	return output, api.update(singleContactEndpoint, contactUUID, input, output)
}

// ListContacts lists all Contacts
//
func (api API) ListContacts(listContactsParams *ListContactsParams) (*Contacts, error) {
	result := &Contacts{}
	query := make([]interface{}, 0, 1)
	if listContactsParams != nil {
		query = append(query, *listContactsParams)
	}
	return result, api.list(contactsEndpoint, result, query...)
}

// MergeContact merges two contacts.
//
func (api API) MergeContacts(intoContactUUID string, fromContactUUID string) (*Contact, error) {
	result := &Contact{}
	temp_path := strings.Replace(mergeContactsEndpoint, ":into_contact_uuid", intoContactUUID, 1)
	path := strings.Replace(temp_path, ":from_contact_uuid", fromContactUUID, 1)
	return result, api.create(path, nil, result)
}

// DeleteContact deletes one contact by UUID.
//
func (api API) DeleteContact(contactUUID string) error {
	return api.delete(singleContactEndpoint, contactUUID)
}
