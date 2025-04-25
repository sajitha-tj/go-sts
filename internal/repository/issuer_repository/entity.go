package issuer_repository

import "crypto/rsa"

type Issuer struct {
	IssuerId   string
	IssuerUrl  string
	PrivateKey *rsa.PrivateKey
}
