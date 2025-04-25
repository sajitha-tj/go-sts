package lib

import (
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/jwt"
)

// newSession creates a new JWT session with the given user.
func NewSession(user string) *oauth2.JWTSession {
	return &oauth2.JWTSession{
		Username: user,
		JWTClaims: &jwt.JWTClaims{
			Subject: user,
			Extra: map[string]interface{}{
				"extra_claim_1": "extra_value_1",
				"extra_claim_2": "extra_value_2",
			},
		},
	}
}
