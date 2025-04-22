package lib

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"log"
)

func KeyGetter(ctx context.Context) (interface{}, error) {
	keys, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Error generating RSA key:", err)
		return nil, err
	}
	return keys, nil
}
