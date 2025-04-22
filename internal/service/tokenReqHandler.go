package service

import (
	"context"
	"log"
	"net/http"

	"github.com/sajitha-tj/go-sts/internal/lib"
)

func (s *OAuthService) HandleTokenRequest(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	log.Println("token request received")
	// Create an empty session object that will be passed to storage implementation to populate (unmarshal) the session into.
	mySessionData := lib.NewSession("")

	accessRequest, err := s.oauth2Provider.NewAccessRequest(ctx, req, mySessionData)
	if err != nil {
		log.Println("Error creating access request from token handler:", err)
		s.oauth2Provider.WriteAccessError(ctx, w, accessRequest, err)
		return
	}

	if mySessionData.Username == "peter" {
		// do something...
	}
	// check scopes...
	for _, scope := range accessRequest.GetRequestedScopes() {
		log.Println("Requested scope:", scope, " granting..")
		accessRequest.GrantScope(scope)
	}

	response, err := s.oauth2Provider.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		log.Println("Error creating access response from token handler:", err)
		s.oauth2Provider.WriteAccessError(ctx, w, accessRequest, err)
		return
	}
	log.Println("response:", response.ToMap())

	// All done, send the response.
	s.oauth2Provider.WriteAccessResponse(ctx, w, accessRequest, response)
	log.Println("access response created from token handler")
}
