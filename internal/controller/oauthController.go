package controller

import (
	"net/http"

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
	ctx := r.Context()
	c.oauthService.HandleAuthorizationRequest(ctx, w, r)
}

// TokenEndpointController handles the token endpoint for OAuth2.
func (c *OAuthController) TokenEndpointController(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	c.oauthService.HandleTokenRequest(ctx, w, r)
}
