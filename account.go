package chartmogul

const (
	accountEndpoint = "account"
)

// Account details in ChartMogul
type Account struct {
	Name        string `json:"name"`
	Currency    string `json:"currency"`
	TimeZone    string `json:"time_zone"`
	WeekStartOn string `json:"week_start_on"`
}

// RetrieveAccount returns details of current account.
func (api API) RetrieveAccount() (*Account, error) {
	result := &Account{}
	accountUUID := ""
	return result, api.retrieve(accountEndpoint, accountUUID, result)
}
