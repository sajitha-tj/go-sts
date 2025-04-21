package service

import (
	"context"
	"log"
	"net/http"

	"github.com/ory/fosite"
)

func (s *OAuthService) HandleAuthorizationRequest(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	log.Println("Handling authorization request")
	ar, err := s.oauth2Provider.NewAuthorizeRequest(ctx, req)
	if err != nil {
		log.Printf("Error creating authorization request: %v", err)
		s.oauth2Provider.WriteAuthorizeError(ctx, w, ar, err)
		return
	}

	// check if the user is logged in and gives his consent
	// for now, checking if username is 'peter'
	if req.Form.Get("username") != "peter" || req.Form.Get("password") != "secret" {
		log.Println("User not authenticated")
		err := fosite.ErrInvalidClient.WithDescription("Invalid credentials")
		s.oauth2Provider.WriteAuthorizeError(ctx, w, ar, err)
	}
	// check scopes...

	// user is authenticated, then...
	mySessionData := &fosite.DefaultSession{
		Username: req.Form.Get("username"),
	}

	response, err := s.oauth2Provider.NewAuthorizeResponse(ctx, ar, mySessionData)
	if err != nil {
		log.Printf("Error creating authorization response: %v", err)
		s.oauth2Provider.WriteAuthorizeError(ctx, w, ar, err)
		return
	}

	s.oauth2Provider.WriteAuthorizeResponse(ctx, w, ar, response)
	log.Println("Authorization response sent")
}
