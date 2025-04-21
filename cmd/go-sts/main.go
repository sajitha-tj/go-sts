package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/sajitha-tj/go-sts/config"
	"github.com/sajitha-tj/go-sts/internal/controller"
	"github.com/sajitha-tj/go-sts/internal/service"
	"github.com/sajitha-tj/go-sts/internal/storage"
)

func main() {
	config := config.GetConfigInstance()

	storage := storage.NewStorage()
	defer storage.Close()
	oauthService := service.NewOauthService(storage)
	oauthController := controller.NewOAuthController(oauthService)

	// Set up the router
	r := mux.NewRouter()
	r.HandleFunc("/authorize", oauthController.AuthorizeEndpointController).Methods("POST")
	r.HandleFunc("/token", oauthController.TokenEndpointController).Methods("POST")

	// Start the server
	log.Println("Starting server on :", config.PORT)
	if err := http.ListenAndServe(":" + config.PORT, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
