package service

import (
	"context"
	"log"
	"net/http"

	"github.com/sajitha-tj/go-sts/internal/lib"
)

func (s *OAuthService) HandleAuthorizationRequest(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	log.Println("Handling authorization request")
	ar, err := s.oauth2Provider.NewAuthorizeRequest(ctx, req)
	if err != nil {
		log.Printf("Error creating authorization request: %v", err)
		s.oauth2Provider.WriteAuthorizeError(ctx, w, ar, err)
		return
	}

	// Check if the user is authenticated (if !authenticated: return)
	if !handleAuthentication(w, req, AuthenticationData{
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
	mySessionData := lib.NewSession(req.Form.Get("username"))

	response, err := s.oauth2Provider.NewAuthorizeResponse(ctx, ar, mySessionData)
	if err != nil {
		log.Printf("Error creating authorization response: %v", err)
		s.oauth2Provider.WriteAuthorizeError(ctx, w, ar, err)
		return
	}

	s.oauth2Provider.WriteAuthorizeResponse(ctx, w, ar, response)
	log.Println("Authorization response sent")
}
