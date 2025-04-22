package lib

import (
	"context"
)

// DefaultAccessTokenIssuerProvider is a default implementation of AccessTokenIssuerProvider.
type AccessTokenIssuerProvider struct {
	issuer string
}

// NewDefaultAccessTokenIssuerProvider creates a new DefaultAccessTokenIssuerProvider.
func NewAccessTokenIssuerProvider(issuer string) *AccessTokenIssuerProvider {
	return &AccessTokenIssuerProvider{issuer: issuer}
}

// GetAccessTokenIssuer returns the access token issuer.
func (p *AccessTokenIssuerProvider) GetAccessTokenIssuer(ctx context.Context) string {
	return p.issuer
}
