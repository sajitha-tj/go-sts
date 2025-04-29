package session_repository

import (
	"net/url"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/sajitha-tj/go-sts/internal/repository/client_repository"
	"golang.org/x/text/language"
)

type StoredRequest struct {
	ID                string                   `json:"id" gorethink:"id"`
	RequestedAt       time.Time                `json:"requestedAt" gorethink:"requestedAt"`
	Client            client_repository.Client `json:"client" gorethink:"client"`
	RequestedScope    fosite.Arguments         `json:"scopes" gorethink:"scopes"`
	GrantedScope      fosite.Arguments         `json:"grantedScopes" gorethink:"grantedScopes"`
	Form              url.Values               `json:"form" gorethink:"form"`
	Session           oauth2.JWTSession        `json:"session" gorethink:"session"`
	RequestedAudience fosite.Arguments         `json:"requestedAudience"`
	GrantedAudience   fosite.Arguments         `json:"grantedAudience"`
	Lang              language.Tag             `json:"-"`
}
