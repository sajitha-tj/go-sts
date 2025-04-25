package oauth_provider

import (
	"context"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"

	"github.com/sajitha-tj/go-sts/internal/configs"
	"github.com/sajitha-tj/go-sts/internal/repository/issuer_repository"
	"github.com/sajitha-tj/go-sts/internal/storage"
)

// NewOauthProvider initializes a new fosite OAuth2Provider instance with necessary configurations.
func NewOauthProvider(storage *storage.Storage) fosite.OAuth2Provider {
	var secret = []byte(configs.GetConfig().FositeConfigs.Secret)
	var fositeConfigs = &fosite.Config{
		AccessTokenLifespan:        time.Minute * 30,
		GlobalSecret:               secret,
		SendDebugMessagesToClients: true,
	}

	var jwtConfig = &JwtConfig{Config: *fositeConfigs}

	oauth2Provider := compose.Compose(
		fositeConfigs,
		storage,
		compose.NewOAuth2JWTStrategy(keyGetter, compose.NewOAuth2HMACStrategy(jwtConfig), jwtConfig),
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.OAuth2StatelessJWTIntrospectionFactory,
	)

	return oauth2Provider
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

// JwtConfig is a struct that holds the configuration for JWT.
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
