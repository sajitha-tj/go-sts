package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sajitha-tj/go-sts/internal/service/idp_service"
)

func IdPRoutes(router *mux.Router, path string, idpService *idp_service.IdPService) {
	routes := router.PathPrefix(path).Subrouter()

	routes.HandleFunc("/login", loginHandler(*idpService)).Methods("GET")
	routes.HandleFunc("/login/callback", loginCallbackHandler(*idpService)).Methods("POST")
}

func loginHandler(service idp_service.IdPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.HandleLogin(w, r)
	}
}

func loginCallbackHandler(service idp_service.IdPService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.HandleLoginCallback(w, r)
	}
}
