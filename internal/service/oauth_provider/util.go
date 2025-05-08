package oauth_provider

import (
	"context"

	"github.com/ory/fosite"
	"github.com/sajitha-tj/go-sts/internal/configs"
	"github.com/sajitha-tj/go-sts/internal/repository/issuer_repository"
)

// JwtConfig is a struct that holds the configuration for JWT which implements the fosite.JwtConfig interface.
// It embeds the fosite.Config struct to inherit its properties.
// This struct is used to configure the JWT strategy in the OAuth2 provider.
type JwtConfig struct {
	fosite.Config
}

// GetAccessTokenIssuer retrieves the iss claim needed for the access token from the context.
// It is used by fosite to determine the issuer of the access token.
func (c *JwtConfig) GetAccessTokenIssuer(ctx context.Context) string {
	issuer := ctx.Value(configs.CTX_ISSUER_KEY).(issuer_repository.Issuer)
	return issuer.IssuerUrl
}

// keyGetter is a function that retrieves the private key needed to sign JWT tokens.
// It is used by fosite to sign the tokens. This implementation returns private keys based on the issuer.
func keyGetter(ctx context.Context) (interface{}, error) {
	issuer, ok := ctx.Value(configs.CTX_ISSUER_KEY).(issuer_repository.Issuer)
	if !ok {
		return nil, fosite.ErrServerError.WithHint("Failed to retrieve issuer from context")
	}
	return issuer.PrivateKey, nil
}
