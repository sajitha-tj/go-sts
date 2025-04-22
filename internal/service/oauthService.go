package service

import (
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"

	"github.com/sajitha-tj/go-sts/config"
	"github.com/sajitha-tj/go-sts/internal/lib"
	"github.com/sajitha-tj/go-sts/internal/storage"
)

type OAuthService struct {
	oauth2Provider fosite.OAuth2Provider
}

func NewOauthService(storage *storage.Storage) *OAuthService {
	var secret = []byte(config.GetConfigInstance().SIGNING_SECRET)
	var fositeConfigs = &fosite.Config{
		AccessTokenLifespan:        time.Minute * 30,
		GlobalSecret:               secret,
		SendDebugMessagesToClients: true,
	}

	var jwtConfig = &config.JwtConfig{Config: *fositeConfigs}

	oauth2Provider := compose.Compose(
		fositeConfigs,
		storage,
		// compose.NewOAuth2HMACStrategy(config),
		compose.NewOAuth2JWTStrategy(lib.KeyGetter, compose.NewOAuth2HMACStrategy(jwtConfig), jwtConfig),
		compose.OAuth2AuthorizeExplicitFactory,
	)

	return &OAuthService{
		oauth2Provider: oauth2Provider,
	}
}
