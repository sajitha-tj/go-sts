package lib

import (
	"encoding/json"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/jwt"
)

// // JWTSession Container for the JWT session.
// type JWTSession struct {
// 	JWTClaims *jwt.JWTClaims
// 	JWTHeader *jwt.Headers
// 	ExpiresAt map[fosite.TokenType]time.Time
// 	Username  string
// 	Subject   string
// }

// type JWTClaims struct {
// 	Subject    string
// 	Issuer     string
// 	Audience   []string
// 	JTI        string
// 	IssuedAt   time.Time
// 	NotBefore  time.Time
// 	ExpiresAt  time.Time
// 	Scope      []string
// 	Extra      map[string]interface{}
// 	ScopeField JWTScopeFieldEnum
// }

func NewSession(user string) *oauth2.JWTSession {
	return &oauth2.JWTSession{
		Username: user,
		JWTClaims: &jwt.JWTClaims{
			Subject: user,
			// Issuer:  "https://example.com",
		},
	}
}

func GetSerializedSession(session fosite.Session) (string, error) {
	sessionData, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	return string(sessionData), nil
}

func DeserializeSession(sessionData string, session fosite.Session) error {
	return json.Unmarshal([]byte(sessionData), &session)
}
