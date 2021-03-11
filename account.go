package chartmogul

const (
	accountEndpoint = "account"
)

// Account details in ChartMogul
type Account struct {
	Name        string `json:"name,omitempty"`
	Currency    string `json:"currency,omitempty"`
	TimeZone    string `json:"time_zone,omitempty"`
	WeekStartOn string `json:"week_start_on,omitempty"`

	Errors Errors `json:"errors,omitempty"`
}

// RetrieveAccount returns details of current account.
func (api API) RetrieveAccount() (*Account, error) {
	result := &Account{}
	accountUUID := ""
	return result, api.retrieve(accountEndpoint, accountUUID, result)
}
