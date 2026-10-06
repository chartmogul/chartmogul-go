package chartmogul

// attributesDefinition is internal struct used to define multiple new custom attributes.
type attributesDefinition struct {
	Email     string                 `json:"email,omitempty"`
	Custom    []*CustomAttribute     `json:"custom"`
	Overrides map[string]interface{} `json:"overrides,omitempty"`
}

// CustomAttributes contains updated custom attributes.
type CustomAttributes struct {
	Custom    map[string]interface{} `json:"custom"`
	Overrides map[string]interface{} `json:"overrides,omitempty"`
	Message   string                 `json:"message,omitempty"`
}

type deleteCustomAttrs struct {
	Custom    []string               `json:"custom"`
	Overrides map[string]interface{} `json:"overrides,omitempty"`
}

// RetrieveCustomersAttributesParams optional query parameters for RetrieveCustomersAttributes.
type RetrieveCustomersAttributesParams struct {
	WithOverrides         *bool  `json:"with_overrides,omitempty"`
	AttributesWithHistory string `json:"attributes_with_history,omitempty"` // Comma-separated attribute names
}

// AddCustomAttributesOptions optional inputs for AddCustomAttributesToCustomerWithOptions.
type AddCustomAttributesOptions struct {
	Overrides map[string]interface{}
}

// AddCustomAttributesWithEmailOptions optional inputs for AddCustomAttributesWithEmailWithOptions.
type AddCustomAttributesWithEmailOptions struct {
	Overrides map[string]interface{}
}

// UpdateCustomAttributesOptions optional inputs for UpdateCustomAttributesOfCustomerWithOptions.
type UpdateCustomAttributesOptions struct {
	Overrides map[string]interface{}
}

// RemoveCustomAttributesOptions optional inputs for RemoveCustomAttributesWithOptions.
type RemoveCustomAttributesOptions struct {
	Overrides map[string]interface{}
}

// CustomAttribute = typed custom attribute.
type CustomAttribute struct {
	Type   string      `json:"type"`
	Key    string      `json:"key"`
	Value  interface{} `json:"value"`
	Source string      `json:"source,omitempty"`
}

// AttributeWithSource covers the special case when you need to update customer's attribute
// and chage the source which shows in the ChartMogul UI.
type AttributeWithSource struct {
	Value  interface{} `json:"value"`
	Source string      `json:"source"`
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
func (api API) RetrieveCustomersAttributes(customerUUID string) (*Attributes, error) {
	return api.RetrieveCustomersAttributesWithOptions(customerUUID, nil)
}

// RetrieveCustomersAttributesWithOptions returns attributes for given customer UUID, with query options.
// A nil opts behaves like RetrieveCustomersAttributes.
func (api API) RetrieveCustomersAttributesWithOptions(customerUUID string, opts *RetrieveCustomersAttributesParams) (*Attributes, error) {
	output := &Attributes{}
	if opts != nil {
		return output, api.retrieveWithParams(customersAttributesEndpoint, customerUUID, output, *opts)
	}
	err := api.retrieve(customersAttributesEndpoint, customerUUID, output)
	return output, err
}

// AddCustomAttributesToCustomer adds custom attributes to specific customer.
func (api API) AddCustomAttributesToCustomer(customerUUID string, customAttributes []*CustomAttribute) (*CustomAttributes, error) {
	return api.AddCustomAttributesToCustomerWithOptions(customerUUID, customAttributes, nil)
}

// AddCustomAttributesToCustomerWithOptions adds custom attributes to a specific customer, with options
// such as override flags. A nil opts behaves like AddCustomAttributesToCustomer.
func (api API) AddCustomAttributesToCustomerWithOptions(customerUUID string, customAttributes []*CustomAttribute, opts *AddCustomAttributesOptions) (*CustomAttributes, error) {
	input := &attributesDefinition{Custom: customAttributes}
	if opts != nil {
		input.Overrides = opts.Overrides
	}
	output := &CustomAttributes{}
	err := api.add(customerCustomAttributesEndpoint, customerUUID, input, output)
	return output, err
}

// AddCustomAttributesWithEmail adds custom attributes to customers with specific email.
func (api API) AddCustomAttributesWithEmail(email string, customAttributes []*CustomAttribute) (*Customers, error) {
	return api.AddCustomAttributesWithEmailWithOptions(email, customAttributes, nil)
}

// AddCustomAttributesWithEmailWithOptions adds custom attributes to customers with specific email, with options
// such as override flags. A nil opts behaves like AddCustomAttributesWithEmail.
func (api API) AddCustomAttributesWithEmailWithOptions(email string, customAttributes []*CustomAttribute, opts *AddCustomAttributesWithEmailOptions) (*Customers, error) {
	input := &attributesDefinition{Email: email, Custom: customAttributes}
	if opts != nil {
		input.Overrides = opts.Overrides
	}
	output := &Customers{}
	err := api.create(customAttributesEndpoint, input, output)
	return output, err
}

// UpdateCustomAttributesOfCustomer updates custom attributes of a specific customer.
func (api API) UpdateCustomAttributesOfCustomer(customerUUID string, customAttributes map[string]interface{}) (*CustomAttributes, error) {
	return api.UpdateCustomAttributesOfCustomerWithOptions(customerUUID, customAttributes, nil)
}

// UpdateCustomAttributesOfCustomerWithOptions updates custom attributes of a specific customer, with options
// such as override flags. A nil opts behaves like UpdateCustomAttributesOfCustomer.
func (api API) UpdateCustomAttributesOfCustomerWithOptions(customerUUID string, customAttributes map[string]interface{}, opts *UpdateCustomAttributesOptions) (*CustomAttributes, error) {
	input := &CustomAttributes{Custom: customAttributes}
	if opts != nil {
		input.Overrides = opts.Overrides
	}
	output := &CustomAttributes{}
	err := api.putTo(customerCustomAttributesEndpoint, customerUUID, input, output)
	return output, err
}

// RemoveCustomAttributes removes a list of custom attributes from a specific customer.
func (api API) RemoveCustomAttributes(customerUUID string, customAttributes []string) (*CustomAttributes, error) {
	return api.RemoveCustomAttributesWithOptions(customerUUID, customAttributes, nil)
}

// RemoveCustomAttributesWithOptions removes a list of custom attributes from a specific customer, with options
// such as override flags. A nil opts behaves like RemoveCustomAttributes.
func (api API) RemoveCustomAttributesWithOptions(customerUUID string, customAttributes []string, opts *RemoveCustomAttributesOptions) (*CustomAttributes, error) {
	input := &deleteCustomAttrs{Custom: customAttributes}
	if opts != nil {
		input.Overrides = opts.Overrides
	}
	output := &CustomAttributes{}
	err := api.deleteWhat(customerCustomAttributesEndpoint, customerUUID, input, output)
	return output, err
}
