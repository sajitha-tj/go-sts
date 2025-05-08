package oauth_provider

import (
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"

	"github.com/sajitha-tj/go-sts/internal/configs"
	"github.com/sajitha-tj/go-sts/internal/storage"
)

type Provider interface {
	GetStorage() *storage.Storage
	fosite.OAuth2Provider
}

type OauthProvider struct {
	storage *storage.Storage
	fosite.OAuth2Provider
}

// NewOauthProvider initializes a new fosite OAuth2Provider instance with necessary configurations.
func NewOauthProvider(storage *storage.Storage) Provider {
	var secret = []byte(configs.GetConfig().FositeConfigs.Secret)
	var fositeConfigs = &fosite.Config{
		AccessTokenLifespan:        time.Minute * 30,
		GlobalSecret:               secret,
		SendDebugMessagesToClients: true,
		RefreshTokenScopes:         []string{"offline"},
	}

	var jwtConfig = &JwtConfig{Config: *fositeConfigs}

	oauth2Provider := compose.Compose(
		fositeConfigs,
		storage,
		compose.NewOAuth2JWTStrategy(keyGetter, compose.NewOAuth2HMACStrategy(jwtConfig), jwtConfig),
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.OAuth2StatelessJWTIntrospectionFactory,
		compose.OAuth2RefreshTokenGrantFactory,
	)

	return &OauthProvider{
		storage:        storage,
		OAuth2Provider: oauth2Provider,
	}
}

// GetStorage returns the storage instance used by the provider.
func (p *OauthProvider) GetStorage() *storage.Storage {
	return p.storage
}
