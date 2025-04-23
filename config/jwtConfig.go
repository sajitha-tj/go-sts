package config

import (
	"context"
	"log"

	"github.com/ory/fosite"
	"github.com/sajitha-tj/go-sts/setup"
)

type JwtConfig struct {
	fosite.Config
}

func (c *JwtConfig) GetAccessTokenIssuer(ctx context.Context) string {
	issuerId := ctx.Value(CTX_KEY_ISSUER).(string)
	issuer, exists := setup.GetTempIssuerDBInstance().GetIssuer(issuerId)
	if !exists {
		log.Printf("No valid issuer for the ID: %s", issuerId)
		return "default-issuer"
	}
	return issuer
}
