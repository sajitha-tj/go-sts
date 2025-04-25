package app

import (
	"github.com/gorilla/mux"
	"github.com/ory/fosite"
	"github.com/sajitha-tj/go-sts/internal/middleware"
	"github.com/sajitha-tj/go-sts/internal/routes"
	"github.com/sajitha-tj/go-sts/internal/service/oauth_provider"
	"github.com/sajitha-tj/go-sts/internal/storage"
)

type AppDependencies struct {
	// authService
	oauthProvider fosite.OAuth2Provider
	storage       *storage.Storage
}

func MakeAPIServer() (*mux.Router, error) {
	deps, err := initializeAppDependencies()
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()

	r.Use(middleware.CtxMiddleware)

	routes.OAuthRoutes(r, "/", &deps.oauthProvider)

	return r, nil
}

func initializeAppDependencies() (*AppDependencies, error) {
	// Initialize dependencies here
	// GetConfigs() etc..
	storage := storage.NewStorage()
	oauthProvider := oauth_provider.NewOauthProvider(storage)
	return &AppDependencies{
		oauthProvider: oauthProvider,
		storage:       storage,
	}, nil
}
