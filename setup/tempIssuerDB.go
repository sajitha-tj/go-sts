package setup

// This is a temporary issuer database that is used for testing purposes.
// It keeps a list of issuers and their corresponding private keys with respective release IDs.
// Database is hardcoded and stored in memory.

import (
	"crypto/rand"
	"crypto/rsa"
)

type TempIssuerDB struct {
	issuers map[string]string
	keyMap  map[string]*rsa.PrivateKey
}

func newTempIssuerDB() *TempIssuerDB {
	// Initialize the issuer database with some example issuers and keys
	myIssuers := map[string]string{
		"123e4567-e89b-12d3-a456-426614174000": "https://my-issuer-one.com",
		"53d45148-f68f-4c1e-8aa8-0a2108a06daa": "https://my-issuer-two.com",
	}
	myKeyMap := make(map[string]*rsa.PrivateKey)
	for k := range myIssuers {
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			panic(err)
		}
		myKeyMap[k] = privateKey
	}

	return &TempIssuerDB{
		issuers: myIssuers,
		keyMap:  myKeyMap,
	}
}

var tempIssuerDB *TempIssuerDB

// GetTempIssuerDBInstance returns a singleton instance of the TempIssuerDB.
func GetTempIssuerDBInstance() *TempIssuerDB {
	if tempIssuerDB == nil {
		tempIssuerDB = newTempIssuerDB()
	}
	return tempIssuerDB
}

// GetIssuer retrieves the issuer URL for a given release ID.
func (db *TempIssuerDB) GetIssuer(releaseId string) (string, bool) {
	issuer, exists := db.issuers[releaseId]
	return issuer, exists
}

// GetPrivateKey retrieves the private key for a given release ID.
func (db *TempIssuerDB) GetPrivateKey(releaseId string) (*rsa.PrivateKey, bool) {
	privateKey, exists := db.keyMap[releaseId]
	return privateKey, exists
}
