package service

import (
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"

	"github.com/sajitha-tj/go-sts/internal/storage"
)

type OAuthService struct {
	oauth2Provider fosite.OAuth2Provider
}

func NewOauthService(storage *storage.Storage) *OAuthService {
	var secret = []byte("my super secret signing password")
	var config = &fosite.Config{
		AccessTokenLifespan: time.Minute * 30,
		GlobalSecret:        secret,
		// can add new issuer
	}

	oauth2Provider := compose.Compose(
		config,
		storage,
		compose.NewOAuth2HMACStrategy(config),
		compose.OAuth2AuthorizeExplicitFactory,
		// compose.OAuth2ClientCredentialsGrantFactory,
	)

	return &OAuthService{
		oauth2Provider: oauth2Provider,
	}
}
