package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ory/fosite"
	"github.com/sajitha-tj/go-sts/internal/lib"
	"github.com/sajitha-tj/go-sts/internal/service/authentication_service"
)

// OAuthRoutes sets up the OAuth2 routes for the given router, with the specified path prefix.
func OAuthRoutes(router *mux.Router, path string, service *authentication_service.AuthenticationService, provider *fosite.OAuth2Provider) {
	routes := router.PathPrefix(path).Subrouter()

	routes.HandleFunc("/authorize", authorizeHandler(*service, *provider)).Methods("GET", "POST")
	routes.HandleFunc("/token", tokenHandler(*provider)).Methods("POST")
	routes.HandleFunc("/introspect", introspectHandler(*provider)).Methods("POST")
}

// authorizeHandler handles the authorization request and response.
func authorizeHandler(service authentication_service.AuthenticationService, provider fosite.OAuth2Provider) http.HandlerFunc {
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
		if !service.HandleAuthentication(w, r, authentication_service.AuthenticationData{
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
		mySessionData := lib.NewSession(r.Form.Get("username"))

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
		mySessionData := lib.NewSession("")

		accessRequest, err := provider.NewAccessRequest(ctx, r, mySessionData)
		if err != nil {
			log.Println("Error creating access request from token handler:", err)
			provider.WriteAccessError(ctx, w, accessRequest, err)
			return
		}

		if mySessionData.GetUsername() == "peter" {
			log.Println("hey pete!")
		}

		// Check and grant requested scopes
		// This only grants the scopes that are requested by the request.
		// scope validation happens from the NewAccessRequest method provided by fosite.
		for _, scope := range accessRequest.GetRequestedScopes() {
			accessRequest.GrantScope(scope)
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

// introspectHandler handles the introspection request and response.
func introspectHandler(provider fosite.OAuth2Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println("introspection request received")
		mySessionData := lib.NewSession("")
		responder, err := provider.NewIntrospectionRequest(ctx, r, mySessionData)
		if err != nil {
			log.Println("Error creating introspection request:", err)
			provider.WriteIntrospectionError(ctx, w, err)
			return
		}

		provider.WriteIntrospectionResponse(ctx, w, responder)
		log.Println("introspection response sent")
	}
}