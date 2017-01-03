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
	enrichmentCustomersAttributesEndpoint      = "customers/:uuid/attributes"
	enrichmentCustomerCustomAttributesEndpoint = "customers/:uuid/attributes/custom"
	enrichmentCustomAttributesEndpoint         = "customers/attributes/custom"
)

// EnrichmentRetrieveCustomersAttributes returns attributes for given customer UUID.
//
// See https://dev.chartmogul.com/reference#customer-attributes
func (api API) EnrichmentRetrieveCustomersAttributes(customerUUID string) (*AttributesResult, error) {
	output := &AttributesResult{}
	err := api.retrieve(enrichmentCustomersAttributesEndpoint, customerUUID, output)
	return output, err
}

// EnrichmentAddCustomAttributesToCustomer adds custom attributes to specific customer.
//
// See https://dev.chartmogul.com/reference#add-custom-attributes-to-customer
func (api API) EnrichmentAddCustomAttributesToCustomer(customerUUID string, customAttributes []*CustomAttribute) (*CustomAttributes, error) {
	output := &CustomAttributes{}
	err := api.add(enrichmentCustomerCustomAttributesEndpoint,
		customerUUID,
		&attributesDefinition{Custom: customAttributes},
		output)
	return output, err
}

// EnrichmentAddCustomAttributesWithEmail adds custom attributes to customers with specific email.
//
// See https://dev.chartmogul.com/reference#add-custom-attributes-to-customers-with-email
func (api API) EnrichmentAddCustomAttributesWithEmail(email string, customAttributes []*CustomAttribute) (*EnrichmentCustomers, error) {
	output := &EnrichmentCustomers{}
	err := api.create(enrichmentCustomAttributesEndpoint,
		&attributesDefinition{Email: email, Custom: customAttributes},
		output)
	return output, err
}

// EnrichmentUpdateCustomAttributesOfCustomer updates custom attributes of a specific customer.
//
// See https://dev.chartmogul.com/reference#update-custom-attributes-of-customer
func (api API) EnrichmentUpdateCustomAttributesOfCustomer(customerUUID string, customAttributes map[string]interface{}) (*CustomAttributes, error) {
	output := &CustomAttributes{}
	err := api.putTo(enrichmentCustomerCustomAttributesEndpoint,
		customerUUID,
		&CustomAttributes{Custom: customAttributes},
		output)
	return output, err
}

// EnrichmentRemoveCustomAttributes removes a list of custom attributes from a specific customer.
//
// See https://dev.chartmogul.com/reference#remove-custom-attributes-from-customer
func (api API) EnrichmentRemoveCustomAttributes(customerUUID string, customAttributes []string) (*CustomAttributes, error) {
	output := &CustomAttributes{}
	err := api.putTo(enrichmentCustomerCustomAttributesEndpoint,
		customerUUID,
		&deleteCustomAttrs{Custom: customAttributes},
		output)
	return output, err
}
