package chartmogul

// TagsResult is necessary for the the result of AddTags.
type TagsResult struct {
	Tags []string `json:"tags"`
}

// TagsByEmail = input for EnrichmentAddTagsToCustomersWithEmail
type TagsByEmail struct {
	Email string   `json:"email"`
	Tags  []string `json:"tags"`
}

const (
	enrichmentCustomerTagsEndpoint = "customers/:uuid/attributes/tags"
	enrichmentTagsEndpoint         = "customers/attributes/tags"
)

// EnrichmentAddTagsToCustomer gives customer new tags.
//
// See https://dev.chartmogul.com/reference#add-tags-to-customer
func (api API) EnrichmentAddTagsToCustomer(customerUUID string, tags []string) (*TagsResult, error) {
	output := &TagsResult{}
	err := api.add(enrichmentCustomerTagsEndpoint,
		customerUUID,
		TagsResult{Tags: tags},
		output)
	return output, err
}

// EnrichmentAddTagsToCustomersWithEmail gives new tags to (multiple) customers
// identified by e-mail only.
//
// See https://dev.chartmogul.com/reference#add-tags-to-customers-with-email
func (api API) EnrichmentAddTagsToCustomersWithEmail(email string, tags []string) (*EnrichmentCustomers, error) {
	output := &EnrichmentCustomers{}
	err := api.create(enrichmentTagsEndpoint,
		&TagsByEmail{Tags: tags},
		output)
	// API doesn't have paging, but the struct does
	output.Page = 1
	output.TotalPages = 1
	output.CurrentPage = 1
	return output, err
}

// EnrichmentRemoveTagsFromCustomer deletes passed tags from customer of given UUID.
//
// See https://dev.chartmogul.com/reference#remove-tags-from-customer
func (api API) EnrichmentRemoveTagsFromCustomer(customerUUID string, tags []string) (*TagsResult, error) {
	output := &TagsResult{}
	err := api.deleteWhat(enrichmentCustomerTagsEndpoint,
		customerUUID,
		&TagsResult{Tags: tags},
		output)
	return output, err
}
