package lib

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func GetIssuerId(host string) (string, error) {
	// slice from . and get the first element
	hostParts := strings.Split(host, ".")
	if len(hostParts) == 0 {
		return "", fmt.Errorf("invalid host: %s", host)
	}
	issuerId := hostParts[0]
	// check if issuerId is empty
	if issuerId == "" {
		return "", fmt.Errorf("invalid host: %s", host)
	}
	// check if issuerId is a valid UUID
	if _, err := uuid.Parse(issuerId); err != nil {
		return "", fmt.Errorf("invalid issuerId: %s", issuerId)
	}
	return issuerId, nil
}
