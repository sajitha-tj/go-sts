package app

import (
	"github.com/gorilla/mux"
	"github.com/ory/fosite"
	"github.com/sajitha-tj/go-sts/internal/configs"
	"github.com/sajitha-tj/go-sts/internal/db"
	"github.com/sajitha-tj/go-sts/internal/middleware"
	"github.com/sajitha-tj/go-sts/internal/routes"
	"github.com/sajitha-tj/go-sts/internal/service/authentication_service"
	"github.com/sajitha-tj/go-sts/internal/service/oauth_provider"
	"github.com/sajitha-tj/go-sts/internal/storage"
)

type AppDependencies struct {
	// authService
	oauthProvider         fosite.OAuth2Provider
	authenticationService authentication_service.AuthenticationService
}

func CreateServer() (*mux.Router, error) {
	deps, err := initializeAppDependencies()
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()

	r.Use(middleware.CtxMiddleware)

	routes.OAuthRoutes(r, "/", &deps.authenticationService, &deps.oauthProvider)

	return r, nil
}

func initializeAppDependencies() (*AppDependencies, error) {
	config := configs.GetConfig()

	db, err := db.New(&config.Database)
	if err != nil {
		return nil, err
	}

	storage := storage.NewStorage(db)

	oauthProvider := oauth_provider.NewOauthProvider(storage)
	authenticationService := authentication_service.NewAuthenticationService(storage.GetUserStore())

	return &AppDependencies{
		oauthProvider:         oauthProvider,
		authenticationService: *authenticationService,
	}, nil
}
