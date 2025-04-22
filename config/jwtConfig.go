package config

import (
	"context"

	"github.com/ory/fosite"
)

type JwtConfig struct {
	fosite.Config
}

func (c *JwtConfig) GetAccessTokenIssuer(ctx context.Context) string {
	return "my-issuer"
}
