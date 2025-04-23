package controller

import (
	"context"
	"log"
	"net/http"

	"github.com/sajitha-tj/go-sts/config"
	"github.com/sajitha-tj/go-sts/internal/lib"
	"github.com/sajitha-tj/go-sts/internal/service"
)

type OAuthController struct {
	oauthService *service.OAuthService
}

func NewOAuthController(oauthService *service.OAuthService) *OAuthController {
	return &OAuthController{
		oauthService: oauthService,
	}
}

// AuthorizeEndpointController handles the authorization endpoint for OAuth2.
func (c *OAuthController) AuthorizeEndpointController(w http.ResponseWriter, r *http.Request) {
	c.handleRequest(w, r, c.oauthService.HandleAuthorizationRequest)
}

// TokenEndpointController handles the token endpoint for OAuth2.
func (c *OAuthController) TokenEndpointController(w http.ResponseWriter, r *http.Request) {
	c.handleRequest(w, r, c.oauthService.HandleTokenRequest)
}

func (c *OAuthController) handleRequest(w http.ResponseWriter, r *http.Request, handler func(context.Context, http.ResponseWriter, *http.Request)) {
	ctx := r.Context()

	// Add issuerId to context
	issuerId, err := lib.GetIssuerId(r.Host)
	if err != nil {
		log.Println("Error getting issuerId:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ctx = context.WithValue(ctx, config.CTX_KEY_ISSUER, issuerId)

	r = r.WithContext(ctx)
	handler(ctx, w, r)
}
