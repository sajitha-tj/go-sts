package storage

import (
	"database/sql"
	"log"

	"github.com/sajitha-tj/go-sts/config"
	"github.com/sajitha-tj/go-sts/internal/repository"
	"github.com/sajitha-tj/go-sts/setup"
)

type Storage struct {
	dbConnector	 *sql.DB
	clientStore	 *repository.ClientStore
	sessionStore *repository.SessionStore
}

// NewStorage initializes a new Storage instance with a database connection and stores.
// Storage instance is responsible for managing the database connection and providing access to the client, user and session stores.
func NewStorage() *Storage {
	dbUser := config.GetConfigInstance().DB_USER
	dbPassword := config.GetConfigInstance().DB_PASSWORD
	dbName := config.GetConfigInstance().DB_NAME

	db, err := sql.Open("postgres", "user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " sslmode=disable")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Initialize the temporary database
	if err := setup.NewTestDB(db).Initialize(); err != nil {
		log.Fatalf("Error initializing temporary database: %v", err)
	}

	clientStore := repository.NewClientStore(db)
	sessionStore := repository.NewSessionStore(db)

	return &Storage{
		dbConnector: db,
		clientStore: clientStore,
		sessionStore: sessionStore,
	}
}

// Close closes the database connection.
// It should be called when the application is shutting down.
func (s *Storage) Close() {
	if err := s.dbConnector.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
	log.Println("Database connection closed successfully")
}