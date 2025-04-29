package dcr_service

type ClientRegistrationRequest struct {
	ClientID      string   `json:"client_id"`
	RedirectURIs  []string `json:"redirect_uris"`
	GrantTypes    []string `json:"grant_types"`
	ResponseTypes []string `json:"response_types"`
	Scopes         []string `json:"scopes"`
	Public        bool     `json:"public"`
	Audience      []string `json:"audience"`
}
