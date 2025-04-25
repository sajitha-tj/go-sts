package oauth_provider

import (
	"context"
	"log"
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
	)

	return oauth2Provider
}

// keyGetter is a function that retrieves the private key for signing JWT tokens.
// It is used by fosite to sign the tokens. This implementation returns private keys based on the issuer.
func keyGetter(ctx context.Context) (interface{}, error) {
	issuerId := ctx.Value(configs.CTX_ISSUER_ID_KEY).(string)
	privateKey, exists := issuer_repository.GetIssuerStoreInstance().GetPrivateKey(issuerId)
	if !exists {
		log.Printf("Private key not found for issuer ID: %s", issuerId)
		return nil, nil
	}
	return privateKey, nil
}

// JwtConfig is a struct that holds the configuration for JWT.
// It embeds the fosite.Config struct to inherit its properties.
// This struct is used to configure the JWT strategy in the OAuth2 provider.
type JwtConfig struct {
	fosite.Config
}

// GetAccessTokenIssuer retrieves the issuer for the access token from the context.
// It is used by fosite to determine the issuer of the access token.
func (c *JwtConfig) GetAccessTokenIssuer(ctx context.Context) string {
	issuerId := ctx.Value(configs.CTX_ISSUER_ID_KEY).(string)
	issuer, exists := issuer_repository.GetIssuerStoreInstance().GetIssuer(issuerId)
	if !exists {
		log.Printf("No valid issuer for the ID: %s", issuerId)
		return "default-issuer"
	}
	return issuer.IssuerUrl
}
