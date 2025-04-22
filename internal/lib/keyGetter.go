package lib

import (
	"context"
	"log"

	"github.com/sajitha-tj/go-sts/setup"
)

func KeyGetter(ctx context.Context) (interface{}, error) {
	releaseId := ctx.Value("releaseId").(string)
	privateKey, exists := setup.GetTempIssuerDBInstance().GetPrivateKey(releaseId)
	if !exists {
		log.Printf("Private key not found for release ID: %s", releaseId)
		return nil, nil
	}
	return privateKey, nil
}
