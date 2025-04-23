package middleware

import (
	"context"
	"log"
	"net/http"
	
	"github.com/sajitha-tj/go-sts/config"
	"github.com/sajitha-tj/go-sts/internal/lib"
)

// ctxMiddleware enriches the request context with required data.
// Values added to the context:
//   - issuerId: The issuerId extracted from the request's Host header.
func CtxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add issuerId to the context of the requst
		issuerId, err := lib.GetIssuerId(r.Host)
		if err != nil {
			log.Println("Error occured while trying to read issuerId:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), config.CTX_ISSUER_ID_KEY, issuerId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
