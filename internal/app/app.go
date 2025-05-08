package app

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/sajitha-tj/go-sts/internal/configs"
	"github.com/sajitha-tj/go-sts/internal/db"
	"github.com/sajitha-tj/go-sts/internal/routes"
	"github.com/sajitha-tj/go-sts/internal/service/dcr_service"
	"github.com/sajitha-tj/go-sts/internal/service/idp_service"
	"github.com/sajitha-tj/go-sts/internal/service/oauth_provider"
	"github.com/sajitha-tj/go-sts/internal/storage"
	"github.com/sajitha-tj/go-sts/setup"
)

type AppDependencies struct {
	// authService
	oauthProvider oauth_provider.Provider
	dcrService    dcr_service.DcrService
	idpService    idp_service.IdPService
}

func CreateServer() (*mux.Router, error) {
	deps, err := initializeAppDependencies()
	if err != nil {
		return nil, err
	}

	r := mux.NewRouter()

	routes.OAuthRoutes(r, "/", &deps.oauthProvider)
	routes.IdPRoutes(r, "/idp", &deps.idpService)
	routes.DcrRoutes(r, "/dcr", &deps.dcrService)
	routes.AuthenticationRoutes(r, "/auth", &deps.oauthProvider)

	return r, nil
}

func initializeAppDependencies() (*AppDependencies, error) {
	config := configs.GetConfig()

	db, err := db.New(&config.Database)
	if err != nil {
		return nil, err
	}

	// Initialize the temporary database
	if err := setup.NewTestDB(db).Initialize(); err != nil {
		log.Fatalf("Error initializing temporary database: %v", err)
	}

	storage := storage.NewStorage(db)

	oauthProvider := oauth_provider.NewOauthProvider(storage)
	dcrService := dcr_service.NewDcrService(storage.GetClientStore())
	idpService := idp_service.NewIdPService(storage.GetUserStore())

	return &AppDependencies{
		oauthProvider: oauthProvider,
		dcrService:    *dcrService,
		idpService:    *idpService,
	}, nil
}
