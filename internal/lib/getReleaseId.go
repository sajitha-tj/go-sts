package lib

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func GetReleaseId(host string) (string, error) {
	// slice from . and get the first element
	hostParts := strings.Split(host, ".")
	if len(hostParts) == 0 {
		return "", fmt.Errorf("invalid host: %s", host)
	}
	releaseId := hostParts[0]
	// check if releaseId is empty
	if releaseId == "" {
		return "", fmt.Errorf("invalid host: %s", host)
	}
	// check if releaseId is a valid UUID
	if _, err := uuid.Parse(releaseId); err != nil {
		return "", fmt.Errorf("invalid releaseId: %s", releaseId)
	}
	return releaseId, nil
}
