package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/sajitha-tj/go-sts/internal/configs"
	"github.com/sajitha-tj/go-sts/internal/repository/issuer_repository"
)

// ctxMiddleware enriches the request context with required data.
// Following values are added to the context:
//   - issuer: The issuer object corresponding to the issuerId extracted from the request host.
func CtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add issuerId to the context of the requst
		issuerId, err := getIssuerId(r.Host)
		if err != nil {
			log.Println("Error occured while trying to read issuerId:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		issuer, exists := issuer_repository.GetIssuerStoreInstance().GetIssuer(issuerId)
		if !exists {
			log.Println("Error occured while trying to read issuerId:", err)
			http.Error(w, "Issuer not found", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), configs.CTX_ISSUER_KEY, issuer)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func getIssuerId(host string) (string, error) {
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
