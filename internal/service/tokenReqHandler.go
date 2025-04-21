package service

import (
	"context"
	"log"
	"net/http"

	"github.com/ory/fosite"
)

func (s *OAuthService) HandleTokenRequest(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	log.Println("token request received")
	// Create an empty session object that will be passed to storage implementation to populate (unmarshal) the session into.
	mySessionData := new(fosite.DefaultSession)

	accessRequest, err := s.oauth2Provider.NewAccessRequest(ctx, req, mySessionData)
	if err != nil {
		log.Println("Error creating access request from token handler:", err)
		s.oauth2Provider.WriteAccessError(ctx, w, accessRequest, err)
		return
	}

	if mySessionData.Username == "peter" {
		// do something...
	}
	
	response, err := s.oauth2Provider.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		log.Println("Error creating access response from token handler:", err)
		s.oauth2Provider.WriteAccessError(ctx, w, accessRequest, err)
		return
	}

	// All done, send the response.
	s.oauth2Provider.WriteAccessResponse(ctx, w, accessRequest, response)
	log.Println("access response created from token handler")
}
