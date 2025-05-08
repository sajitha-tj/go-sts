package routes

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/ory/fosite"
	"github.com/sajitha-tj/go-sts/internal/service/oauth_provider"
)

type UserClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AcceptLoginRequestPayload struct {
	FlowId     string      `json:"flowId"`
	Success    bool        `json:"success"`
	UserClaims *UserClaims `json:"userClaims"`
}

type LoginAcceptedResponse struct {
	ResponseType fosite.Arguments `json:"response_type"`
	RedirectURI  url.URL          `json:"redirect_uri"`
	ClientID     string           `json:"client_id"`
	Scope        fosite.Arguments `json:"scope"`
	State        string           `json:"state"`
}

func AuthenticationRoutes(router *mux.Router, path string, p *oauth_provider.Provider) {
	routes := router.PathPrefix(path).Subrouter()

	routes.HandleFunc("/login/accept", acceptLoginHandler(*p)).Methods("PUT")
}

func acceptLoginHandler(p oauth_provider.Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Parse the request body
		var payload AcceptLoginRequestPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Error decoding JSON payload", http.StatusBadRequest)
			return
		}
		flowId := payload.FlowId
		if flowId == "" {
			http.Error(w, "flowId is required", http.StatusBadRequest)
			return
		}

		if err := p.GetStorage().GetSessionStore().AuthenticateAuthorizeRequestSession(ctx, flowId); err != nil {
			http.Error(w, "Error authenticating user", http.StatusInternalServerError)
			return
		}

		ar, err := p.GetStorage().GetSessionStore().GetAuthorizeRequestSession(ctx, flowId)
		if err != nil {
			http.Error(w, "Error getting authorize request session", http.StatusInternalServerError)
			return
		}

		// Create the response
		response := LoginAcceptedResponse{
			ResponseType: ar.GetResponseTypes(),
			RedirectURI:  *ar.GetRedirectURI(),
			ClientID:     ar.GetClient().GetID(),
			Scope:        ar.GetRequestedScopes(),
			State:        ar.GetState(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			return
		}
	}
}
