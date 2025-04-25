package main

import (
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/sajitha-tj/go-sts/internal/configs"
	"github.com/sajitha-tj/go-sts/internal/app"
)

func main() {
	// Configs
	if err := configs.LoadConfigs(); err != nil {
		log.Fatal("Error loading configs:", err)
	}

	app, err := app.MakeAPIServer()
	if err != nil {
		log.Fatal("Error creating API server:", err)
	}

	port := configs.GetConfig().Server.Port

	// Start the server
	log.Println("Starting server on :", port)
	if err := http.ListenAndServe(":"+port, app); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
