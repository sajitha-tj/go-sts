package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/token/jwt"
	"github.com/sajitha-tj/go-sts/internal/service/authentication_service"
)

// OAuthRoutes sets up the OAuth2 routes for the given router, with the specified path prefix.
func OAuthRoutes(router *mux.Router, path string, provider *fosite.OAuth2Provider) {
	routes := router.PathPrefix(path).Subrouter()

	routes.HandleFunc("/authorize", authorizeHandler(*provider)).Methods("GET", "POST")
	routes.HandleFunc("/token", tokenHandler(*provider)).Methods("POST")
}

// authorizeHandler handles the authorization request and response.
func authorizeHandler(provider fosite.OAuth2Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println("Handling authorization request")
		ar, err := provider.NewAuthorizeRequest(ctx, r)
		if err != nil {
			log.Printf("Error creating authorization request: %v", err)
			provider.WriteAuthorizeError(ctx, w, ar, err)
			return
		}

		// Check if the user is authenticated (if !authenticated: return)
		if !authentication_service.HandleAuthentication(w, r, authentication_service.AuthenticationData{
			ResponseType: ar.GetResponseTypes()[0],
			ClientID:     ar.GetClient().GetID(),
			RedirectURI:  ar.GetRedirectURI().String(),
			Scope:        ar.GetRequestedScopes()[0],
			State:        ar.GetState(),
			Nonce:        "random_nonce",
		}) {
			return
		}

		// Check requested scopes
		for _, scope := range ar.GetRequestedScopes() {
			ar.GrantScope(scope)
		}

		// user is authenticated, then...
		mySessionData := newSession(r.Form.Get("username"))

		response, err := provider.NewAuthorizeResponse(ctx, ar, mySessionData)
		if err != nil {
			log.Printf("Error creating authorization response: %v", err)
			provider.WriteAuthorizeError(ctx, w, ar, err)
			return
		}

		provider.WriteAuthorizeResponse(ctx, w, ar, response)
		log.Println("Authorization response sent")
	}
}

// tokenHandler handles the token request and response.
func tokenHandler(provider fosite.OAuth2Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println("token request received")
		// Create an empty session object that will be passed to storage implementation to populate (unmarshal) the session into.
		mySessionData := newSession("")

		accessRequest, err := provider.NewAccessRequest(ctx, r, mySessionData)
		if err != nil {
			log.Println("Error creating access request from token handler:", err)
			provider.WriteAccessError(ctx, w, accessRequest, err)
			return
		}

		if mySessionData.GetUsername() == "peter" {
			log.Println("hey pete!")
		}

		response, err := provider.NewAccessResponse(ctx, accessRequest)
		if err != nil {
			log.Println("Error creating access response from token handler:", err)
			provider.WriteAccessError(ctx, w, accessRequest, err)
			return
		}

		// All done, send the response.
		provider.WriteAccessResponse(ctx, w, accessRequest, response)
		log.Println("access response created from token handler")
	}
}

// newSession creates a new JWT session with the given user.
func newSession(user string) *oauth2.JWTSession {
	return &oauth2.JWTSession{
		Username: user,
		JWTClaims: &jwt.JWTClaims{
			Subject: user,
			// Issuer:  "https://example.com",
		},
	}
}
