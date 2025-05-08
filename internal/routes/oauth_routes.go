package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ory/fosite"
	"github.com/sajitha-tj/go-sts/internal/lib"
	"github.com/sajitha-tj/go-sts/internal/middleware"
	"github.com/sajitha-tj/go-sts/internal/service/oauth_provider"
)

// OAuthRoutes sets up the OAuth2 routes for the given router, with the specified path prefix.
func OAuthRoutes(router *mux.Router, path string, p *oauth_provider.Provider) {
	routes := router.PathPrefix(path).Subrouter()
	routes.Use(middleware.CtxMiddleware)

	routes.HandleFunc("/authorize", authorizeHandler(*p)).Methods("GET", "POST")
	routes.HandleFunc("/token", tokenHandler(*p)).Methods("POST")
	routes.HandleFunc("/introspect", introspectHandler(*p)).Methods("POST")
}

// authorizeHandler handles the authorization request and response.
func authorizeHandler(p oauth_provider.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ar, err := p.NewAuthorizeRequest(ctx, r)
		if err != nil {
			log.Printf("Error creating authorization request: %v", err)
			p.WriteAuthorizeError(ctx, w, ar, err)
			return
		}

		flowId := r.URL.Query().Get("flowId")
		if flowId == "" {
			// If flowId is not provided, user is redirected to the IdP login page.
			flowId, err := p.GetStorage().GetSessionStore().CreateAuthorizeRequestSession(ctx, ar)
			if err != nil {
				log.Println("Error storing request data", err)
			}

			w.Header().Set("Location", "http://localhost:8080/idp/login?flowId="+flowId)
			w.WriteHeader(http.StatusFound)
			return
		}

		// If flowId is provided, we assume the user is already authenticated.
		isAuthenticated, err := p.GetStorage().GetSessionStore().IsRequestSessionAuthenticated(ctx, flowId)
		if err != nil {
			log.Println("Error checking authentication status:", err)
			http.Error(w, "Invalid flowId", http.StatusBadRequest)
			return
		}

		if !isAuthenticated {
			log.Println("User is not authenticated")
			http.Error(w, "User is not authenticated", http.StatusUnauthorized)
			return
		}

		// User is authenticated..
		// Check requested scopes
		for _, scope := range ar.GetRequestedScopes() {
			ar.GrantScope(scope)
		}

		// user is authenticated, then...
		mySessionData := lib.NewSession(r.Form.Get("username"))

		response, err := p.NewAuthorizeResponse(ctx, ar, mySessionData)
		if err != nil {
			log.Printf("Error creating authorization response: %v", err)
			p.WriteAuthorizeError(ctx, w, ar, err)
			return
		}

		p.WriteAuthorizeResponse(ctx, w, ar, response)
		log.Println("Authorization response sent")
	}
}

// tokenHandler handles the token request and response.
func tokenHandler(p fosite.OAuth2Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println("token request received")
		// Create an empty session object that will be passed to storage implementation to populate (unmarshal) the session into.
		mySessionData := lib.NewSession("")

		accessRequest, err := p.NewAccessRequest(ctx, r, mySessionData)
		if err != nil {
			log.Println("Error creating access request from token handler:", err)
			p.WriteAccessError(ctx, w, accessRequest, err)
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

		response, err := p.NewAccessResponse(ctx, accessRequest)
		if err != nil {
			log.Println("Error creating access response from token handler:", err)
			p.WriteAccessError(ctx, w, accessRequest, err)
			return
		}

		// All done, send the response.
		p.WriteAccessResponse(ctx, w, accessRequest, response)
		log.Println("access response created from token handler")
	}
}

// introspectHandler handles the introspection request and response.
func introspectHandler(p fosite.OAuth2Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println("introspection request received")
		mySessionData := lib.NewSession("")
		responder, err := p.NewIntrospectionRequest(ctx, r, mySessionData)
		if err != nil {
			log.Println("Error creating introspection request:", err)
			p.WriteIntrospectionError(ctx, w, err)
			return
		}

		p.WriteIntrospectionResponse(ctx, w, responder)
		log.Println("introspection response sent")
	}
}
