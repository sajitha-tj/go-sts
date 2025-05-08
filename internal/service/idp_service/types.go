package idp_service

import "net/url"

type LoginFormData struct {
	FlowId string
}

type UserClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AcceptLoginRequestPayload struct {
	FlowId     string      `json:"flowId"`
	Success    bool        `json:"success"`
	UserClaims *UserClaims `json:"userClaims"`
}

type AcceptLoginResponseData struct {
	ResponseType []string `json:"response_type"`
	RedirectURI  url.URL  `json:"redirect_uri"`
	ClientID     string   `json:"client_id"`
	Scope        []string `json:"scope"`
	State        string   `json:"state"`
}
