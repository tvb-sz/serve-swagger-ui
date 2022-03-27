package conf

// microsoft oauth config
type microsoft struct {
	ClientID     string `json:"client_id"`     // microsoft oauth client_id
	ClientSecret string `json:"client_secret"` // microsoft oauth client_secret
	Tenant       string `json:"tenant"`        // microsoft oauth tenant
}
