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
	Note               string                 `json:"note,omitempty"`
	Phone              string                 `json:"phone,omitempty"`
	Position           uint32                 `json:"position,omitempty"`
	Title              string                 `json:"title,omitempty"`
	Twitter            string                 `json:"twitter,omitempty"`
	Custom             map[string]interface{} `json:"custom,omitempty"`
}

// UpdateContact allows updating contact on the update endpoint.
type UpdateContact struct {
	CustomerExternalID string   `json:"customer_external_id,omitempty"`
	DataSourceUUID     string   `json:"data_source_uuid,omitempty"`
	FirstName          string   `json:"first_name,omitempty"`
	LastName           string   `json:"last_name,omitempty"`
	LinkedIn           string   `json:"linked_in,omitempty"`
	Note               string   `json:"note,omitempty"`
	Phone              string   `json:"phone,omitempty"`
	Position           uint32   `json:"position,omitempty"`
	Title              string   `json:"title,omitempty"`
	Twitter            string   `json:"twitter,omitempty"`
	Custom             []Custom `json:"custom,omitempty"`
}

// NewContact allows creating contact on a new endpoint.
type NewContact struct {
	// Obligatory
	CustomerUUID   string `json:"customer_uuid,omitempty"`
	DataSourceUUID string `json:"data_source_uuid,omitempty"`

	//Optional
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	LinkedIn  string   `json:"linked_in,omitempty"`
	Note      string   `json:"note,omitempty"`
	Phone     string   `json:"phone,omitempty"`
	Position  uint32   `json:"position,omitempty"`
	Title     string   `json:"title,omitempty"`
	Twitter   string   `json:"twitter,omitempty"`
	Custom    []Custom `json:"custom,omitempty"`
}

type Custom struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"string,omitempty"`
}

// ListContactsParams = parameters for listing contacts in API.
type ListContactsParams struct {
	CustomerUUID   string `json:"customer_uuid,omitempty"`
	DataSourceUUID string `json:"data_source_uuid,omitempty"`
	PaginationWithCursor
}

// Contacts is result of listing contacts in API.
type Contacts struct {
	Entries []*Contact `json:"entries,omitempty"`
	Cursor  string     `json:"cursor,omitempty"`
	PerPage uint32     `json:"per_page"`
	HasMore bool       `json:"has_more,omitempty"`
}

const (
	singleContactEndpoint = "contacts/:uuid"
	contactsEndpoint      = "contacts"
	mergeContactsEndpoint = "contacts/:into_contact_uuid/merge/:from_contact_uuid"
)

// CreateContact loads the contact to Chartmogul
//
// See https://dev.chartmogul.com/reference/create-a-contact-contacts
func (api API) CreateContact(newContact *NewContact) (*Contact, error) {
	result := &Contact{}
	return result, api.create(contactsEndpoint, newContact, result)
}

// RetrieveContact returns one contact as in API.
//
// See https://dev.chartmogul.com/reference/retrieve-a-contact
func (api API) RetrieveContact(contactUUID string) (*Contact, error) {
	result := &Contact{}
	return result, api.retrieve(singleContactEndpoint, contactUUID, result)
}

// UpdateContact updates one contact in API.
//
// See https://dev.chartmogul.com/reference/retrieve-a-contact
func (api API) UpdateContact(input *UpdateContact, contactUUID string) (*Contact, error) {
	output := &Contact{}
	return output, api.update(singleContactEndpoint, contactUUID, input, output)
}

// ListContacts lists all Contacts
//
// See https://dev.chartmogul.com/reference/list-all-contacts
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
// See https://dev.chartmogul.com/reference/merge-contacts
func (api API) MergeContacts(intoContactUUID string, fromContactUUID string) (*Contact, error) {
	result := &Contact{}
	temp_path := strings.Replace(mergeContactsEndpoint, ":into_contact_uuid", intoContactUUID, 1)
	path := strings.Replace(temp_path, ":from_contact_uuid", fromContactUUID, 1)
	return result, api.create(path, nil, result)
}

// DeleteContact deletes one contact by UUID.
//
// See https://dev.chartmogul.com/reference/delete-a-contact
func (api API) DeleteContact(contactUUID string) error {
	return api.delete(singleContactEndpoint, contactUUID)
}
