package lib

import (
	"context"
	"log"

	"github.com/sajitha-tj/go-sts/config"
	"github.com/sajitha-tj/go-sts/setup"
)

func KeyGetter(ctx context.Context) (interface{}, error) {
	issuerId := ctx.Value(config.CTX_KEY_ISSUER).(string)
	privateKey, exists := setup.GetTempIssuerDBInstance().GetPrivateKey(issuerId)
	if !exists {
		log.Printf("Private key not found for issuer ID: %s", issuerId)
		return nil, nil
	}
	return privateKey, nil
}
