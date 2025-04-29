package client_repository

import (
	"github.com/mohae/deepcopy"
	"github.com/ory/fosite"

	"github.com/sajitha-tj/go-sts/internal/lib"
)

type Client struct {
	ClientID       string              `json:"client_id"`
	ClientSecret   string              `json:"client_secret"`
	RotatedSecrets lib.JSONStringArray `json:"rotated_secrets"`
	RedirectURIs   lib.JSONStringArray `json:"redirect_uris"`
	GrantTypes     lib.FositeArguments `json:"grant_types"`
	ResponseTypes  lib.FositeArguments `json:"response_types"`
	Scopes         lib.FositeArguments `json:"scopes"`
	Public         bool                `json:"public"`
	Audience       lib.FositeArguments `json:"audience"`
}

// GetID returns the client ID.
func (c Client) GetID() string {
	return c.ClientID
}

// GetHashedSecret returns the hashed secret as it is stored in the store.
func (c Client) GetHashedSecret() []byte {
	return []byte(c.ClientSecret)
}

// GetRedirectURIs returns the client's allowed redirect URIs.
func (c Client) GetRedirectURIs() []string {
	return c.RedirectURIs
}

// GetGrantTypes returns the client's allowed grant types.
func (c Client) GetGrantTypes() fosite.Arguments {
	return fosite.Arguments(c.GrantTypes)
}

// GetResponseTypes returns the client's allowed response types.
// All allowed combinations of response types have to be listed, each combination having
// response types of the combination separated by a space.
func (c Client) GetResponseTypes() fosite.Arguments {
	return fosite.Arguments(c.ResponseTypes)
}

// GetScopes returns the scopes this client is allowed to request.
func (c Client) GetScopes() fosite.Arguments {
	return fosite.Arguments(c.Scopes)
}

// IsPublic returns true, if this client is marked as public.
func (c Client) IsPublic() bool {
	return c.Public
}

// GetAudience returns the allowed audience(s) for this client.
func (c Client) GetAudience() fosite.Arguments {
	return fosite.Arguments(c.Audience)
}

func (c Client) CloneWithSecret(secret string) Client {
	newClient := deepcopy.Copy(c).(Client)
	newClient.ClientSecret = secret
	newClient.RotatedSecrets = nil
	return newClient
}