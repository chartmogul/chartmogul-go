package chartmogul

// TagsResult is necessary for the the result of AddTags.
type TagsResult struct {
	Tags []string `json:"tags"`
}

// TagsByEmail = input for AddTagsToCustomersWithEmail
type TagsByEmail struct {
	Email string   `json:"email"`
	Tags  []string `json:"tags"`
}

const (
	customerTagsEndpoint = "customers/:uuid/attributes/tags"
	tagsEndpoint         = "customers/attributes/tags"
)

// AddTagsToCustomer gives customer new tags.
//
// See https://dev.chartmogul.com/v1.0/reference#tags
func (api API) AddTagsToCustomer(customerUUID string, tags []string) (*TagsResult, error) {
	output := &TagsResult{}
	err := api.add(customerTagsEndpoint,
		customerUUID,
		TagsResult{Tags: tags},
		output)
	return output, err
}

// AddTagsToCustomersWithEmail gives new tags to (multiple) customers
// identified by e-mail only.
//
// See https://dev.chartmogul.com/v1.0/reference#tags
func (api API) AddTagsToCustomersWithEmail(email string, tags []string) (*Customers, error) {
	output := &Customers{}
	err := api.create(tagsEndpoint,
		&TagsByEmail{Tags: tags},
		output)
	// API doesn't have paging, but the struct does
	output.Page = 1
	output.TotalPages = 1
	output.CurrentPage = 1
	return output, err
}

// RemoveTagsFromCustomer deletes passed tags from customer of given UUID.
//
// See https://dev.chartmogul.com/v1.0/reference#tags
func (api API) RemoveTagsFromCustomer(customerUUID string, tags []string) (*TagsResult, error) {
	output := &TagsResult{}
	err := api.deleteWhat(customerTagsEndpoint,
		customerUUID,
		&TagsResult{Tags: tags},
		output)
	return output, err
}
