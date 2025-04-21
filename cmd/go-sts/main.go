package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/sajitha-tj/go-sts/internal/controller"
	"github.com/sajitha-tj/go-sts/internal/service"
	"github.com/sajitha-tj/go-sts/internal/storage"
)

func main() {
	storage := storage.NewStorage()
	defer storage.Close()
	oauthService := service.NewOauthService(storage)
	oauthController := controller.NewOAuthController(oauthService)

	// Set up the router
	r := mux.NewRouter()
	r.HandleFunc("/authorize", oauthController.AuthorizeEndpointController).Methods("POST")
	r.HandleFunc("/token", oauthController.TokenEndpointController).Methods("POST")

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
