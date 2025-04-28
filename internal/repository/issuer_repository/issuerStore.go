package issuer_repository

import (
	"crypto/rand"
	"crypto/rsa"
)

type IssuerStore struct {
	issuers map[string]Issuer
}

func newIssuerStore() *IssuerStore {
	issuerOne := Issuer{
		IssuerId:   "123e4567-e89b-12d3-a456-426614174000",
		IssuerUrl:  "https://my-issuer-one.com",
		PrivateKey: generatePrivateKey(),
	}
	issuerTwo := Issuer{
		IssuerId:   "53d45148-f68f-4c1e-8aa8-0a2108a06daa",
		IssuerUrl:  "https://my-issuer-two.com",
		PrivateKey: generatePrivateKey(),
	}
	// Initialize the issuer store with some example issuers and keys
	issuers := map[string]Issuer{
		issuerOne.IssuerId: issuerOne,
		issuerTwo.IssuerId: issuerTwo,
	}

	return &IssuerStore{
		issuers: issuers,
	}
}

var issuerStoreInstance *IssuerStore

// GetIssuerStoreInstance returns a singleton instance of the IssuerStore.
func GetIssuerStoreInstance() *IssuerStore {
	if issuerStoreInstance == nil {
		issuerStoreInstance = newIssuerStore()
	}
	return issuerStoreInstance
}

// GetIssuer retrieves the issuer URL for a given issuer ID.
func (s *IssuerStore) GetIssuer(issuerId string) (Issuer, bool) {
	issuer, exists := s.issuers[issuerId]
	// Return the issuer and whether it exists
	return issuer, exists
}

// generatePrivateKey generates a new RSA private key.
func generatePrivateKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	return privateKey
}