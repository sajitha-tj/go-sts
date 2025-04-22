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
	relaseId := ctx.Value("releaseId").(string)
	issuer, exists := setup.GetTempIssuerDBInstance().GetIssuer(relaseId)
	if !exists {
		log.Printf("Issuer not found for release ID: %s", relaseId)
		return "default-issuer"
	}
	return issuer
}
