package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sajitha-tj/go-sts/internal/service/dcr_service"
)

func DcrRoutes(router *mux.Router, path string, service *dcr_service.DcrService) {
	routes := router.PathPrefix(path).Subrouter()
	routes.HandleFunc("/register", registerHandler(*service)).Methods("POST")
}

func registerHandler(service dcr_service.DcrService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		service.RegisterClient(w, r, ctx)
	}
}
