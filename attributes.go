package chartmogul

// AttributesResult is necessary for the the GET attributes call, has just one field.
type AttributesResult struct {
	Attributes *Attributes `json:"attributes"`
}

// attributesDefinition is internal struct used to define multiple new custom attributes.
type attributesDefinition struct {
	Email  string             `json:"email,omitempty"`
	Custom []*CustomAttribute `json:"custom"`
}

// CustomAttributes contains updated custom attributes.
type CustomAttributes struct {
	Custom map[string]interface{} `json:"custom"`
}

type deleteCustomAttrs struct {
	Custom []string `json:"custom"`
}

// CustomAttribute = typed custom attribute.
type CustomAttribute struct {
	Type  string      `json:"type"`
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

const (
	customersAttributesEndpoint      = "customers/:uuid/attributes"
	customerCustomAttributesEndpoint = "customers/:uuid/attributes/custom"
	customAttributesEndpoint         = "customers/attributes/custom"

	// AttrTypeString is one of the possible data types for custom attributes.
	AttrTypeString = "String"
	// AttrTypeInteger is one of the possible data types for custom attributes.
	AttrTypeInteger = "Integer"
	// AttrTypeTimestamp is one of the possible data types for custom attributes.
	AttrTypeTimestamp = "Timestamp"
	// AttrTypeBoolean is one of the possible data types for custom attributes.
	AttrTypeBoolean = "Boolean"
)

// RetrieveCustomersAttributes returns attributes for given customer UUID.
//
// See https://dev.chartmogul.com/reference#customer-attributes
func (api API) RetrieveCustomersAttributes(customerUUID string) (*AttributesResult, error) {
	output := &AttributesResult{}
	err := api.retrieve(customersAttributesEndpoint, customerUUID, output)
	return output, err
}

// AddCustomAttributesToCustomer adds custom attributes to specific customer.
//
// See https://dev.chartmogul.com/reference#add-custom-attributes-to-customer
func (api API) AddCustomAttributesToCustomer(customerUUID string, customAttributes []*CustomAttribute) (*CustomAttributes, error) {
	output := &CustomAttributes{}
	err := api.add(customerCustomAttributesEndpoint,
		customerUUID,
		&attributesDefinition{Custom: customAttributes},
		output)
	return output, err
}

// AddCustomAttributesWithEmail adds custom attributes to customers with specific email.
//
// See https://dev.chartmogul.com/reference#add-custom-attributes-to-customers-with-email
func (api API) AddCustomAttributesWithEmail(email string, customAttributes []*CustomAttribute) (*Customers, error) {
	output := &Customers{}
	err := api.create(customAttributesEndpoint,
		&attributesDefinition{Email: email, Custom: customAttributes},
		output)
	return output, err
}

// UpdateCustomAttributesOfCustomer updates custom attributes of a specific customer.
//
// See https://dev.chartmogul.com/reference#update-custom-attributes-of-customer
func (api API) UpdateCustomAttributesOfCustomer(customerUUID string, customAttributes map[string]interface{}) (*CustomAttributes, error) {
	output := &CustomAttributes{}
	err := api.putTo(customerCustomAttributesEndpoint,
		customerUUID,
		&CustomAttributes{Custom: customAttributes},
		output)
	return output, err
}

// RemoveCustomAttributes removes a list of custom attributes from a specific customer.
//
// See https://dev.chartmogul.com/reference#remove-custom-attributes-from-customer
func (api API) RemoveCustomAttributes(customerUUID string, customAttributes []string) (*CustomAttributes, error) {
	output := &CustomAttributes{}
	err := api.putTo(customerCustomAttributesEndpoint,
		customerUUID,
		&deleteCustomAttrs{Custom: customAttributes},
		output)
	return output, err
}
