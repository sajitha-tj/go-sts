package storage

import (
	"database/sql"
	"log"

	"github.com/sajitha-tj/go-sts/internal/repository/client_repository"
	"github.com/sajitha-tj/go-sts/internal/repository/session_repository"
	"github.com/sajitha-tj/go-sts/internal/repository/user_repository"
	"github.com/sajitha-tj/go-sts/setup"
)

type Storage struct {
	dbConnector  *sql.DB
	clientStore  *client_repository.ClientStore
	sessionStore *session_repository.SessionStore
	userStore    *user_repository.UserStore
}

// NewStorage initializes a new Storage instance with a database connection and stores.
// Storage instance is responsible for managing the database connection and providing access to the client, user and session stores.
func NewStorage(db *sql.DB) *Storage {
	// Initialize the temporary database
	if err := setup.NewTestDB(db).Initialize(); err != nil {
		log.Fatalf("Error initializing temporary database: %v", err)
	}

	clientStore := client_repository.NewClientStore(db)
	sessionStore := session_repository.NewSessionStore(db)
	userStore := user_repository.NewUserStore(db)

	return &Storage{
		dbConnector:  db,
		clientStore:  clientStore,
		sessionStore: sessionStore,
		userStore:    userStore,
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

func (s *Storage) GetClientStore() *client_repository.ClientStore {
	return s.clientStore
}

func (s *Storage) GetSessionStore() *session_repository.SessionStore {
	return s.sessionStore
}

func (s *Storage) GetUserStore() *user_repository.UserStore {
	return s.userStore
}
