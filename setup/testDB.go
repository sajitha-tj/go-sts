package setup

// This is a temporary database setup for testing purposes.
// It creates a set of tables and populates them with test data.
// Data are stored in a PostgreSQL database, and the database connection is passed as a parameter.

import (
	"database/sql"
	"log"
)

const (
	// Table names
	AuthorizationCodeSessionsTable = "authorization_code_sessions"
	AccessTokenSessionsTable       = "access_token_sessions"
	RefreshTokenSessionsTable      = "refresh_token_sessions"
	UsersTable                     = "users"
	ClientsTable                   = "clients"
	AuthorizeRequestTable          = "authorize_requests"
)

type TestDB struct {
	db *sql.DB
}

func NewTestDB(db *sql.DB) *TestDB {
	return &TestDB{db: db}
}

func (t *TestDB) Initialize() error {
	if err := t.dropTables(); err != nil {
		return err
	}
	if err := t.createTables(); err != nil {
		return err
	}
	if err := t.populateTables(); err != nil {
		return err
	}
	log.Println("Database initialized with test data.")
	return nil
}

func (t *TestDB) dropTables() error {
	queries := []string{
		"DROP TABLE IF EXISTS " + AuthorizationCodeSessionsTable,
		"DROP TABLE IF EXISTS " + AccessTokenSessionsTable,
		"DROP TABLE IF EXISTS " + RefreshTokenSessionsTable,
		"DROP TABLE IF EXISTS " + UsersTable,
		"DROP TABLE IF EXISTS " + ClientsTable,
		"DROP TABLE IF EXISTS " + AuthorizeRequestTable,
	}

	for _, query := range queries {
		if _, err := t.db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func (t *TestDB) createTables() error {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS ` + AuthorizationCodeSessionsTable + ` (
			code TEXT PRIMARY KEY,
			active BOOLEAN,
			req_id TEXT,
			requested_at TIMESTAMP,
			request_data JSONB,
			client_id TEXT
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS ` + AccessTokenSessionsTable + ` (
			signature TEXT PRIMARY KEY,
			active BOOLEAN,
			req_id TEXT,
			requested_at TIMESTAMP,
			request_data JSONB,
			client_id TEXT
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS ` + RefreshTokenSessionsTable + ` (
			signature TEXT PRIMARY KEY,
			active BOOLEAN,
			req_id TEXT,
			requested_at TIMESTAMP,
			request_data JSONB,
			client_id TEXT
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS ` + UsersTable + ` (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE,
			password TEXT,
			email TEXT UNIQUE,
			created_at TIMESTAMP
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS ` + ClientsTable + ` (
			id TEXT PRIMARY KEY,
			secret TEXT,
			rotated_secrets JSONB,
			redirect_uris JSONB,
			grant_types JSONB,
			response_types JSONB,
			scopes JSONB,
			public BOOLEAN,
			audience JSONB
		)
		`,
		`
		CREATE TABLE IF NOT EXISTS ` + AuthorizeRequestTable + ` (
			req_id TEXT PRIMARY KEY,
			request_data JSONB,
			authenticated BOOLEAN DEFAULT FALSE,
			requested_at TIMESTAMP,
			exp_at TIMESTAMP
		)
		`,
	}

	for _, query := range queries {
		if _, err := t.db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}

func (t *TestDB) populateTables() error {
	queries := []string{
		`
		INSERT INTO ` + UsersTable + ` (id, username, password, email, created_at)
		VALUES ('user_1', 'peter', 'secret', 'peter@mail.com', CURRENT_TIMESTAMP)
		ON CONFLICT DO NOTHING
		`,
		`
		INSERT INTO ` + ClientsTable + ` (id, secret, rotated_secrets, redirect_uris, grant_types, response_types, scopes, public, audience)
		VALUES (
			'my-client',
			'$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO', -- Hashed secret: foobar
			'["$2y$10$X51gLxUQJ.hGw1epgHTE5u0bt64xM0COU7K9iAp.OFg8p2pUd.1zC"]', -- Rotated secrets: ["foobaz"]
			'["http://localhost:3846/callback"]', -- Redirect URIs as JSON array
			'["implicit", "refresh_token", "authorization_code", "password", "client_credentials"]', -- Grant types as JSON array
			'["id_token", "code", "token", "id_token token", "code id_token", "code token", "code id_token token"]', -- Response types as JSON array
			'["fosite", "openid", "photos", "offline", "offline_access"]', -- Scopes as JSON array
			false, -- Public
			'["example_audience"]' -- Audience as JSON array
		)
		ON CONFLICT DO NOTHING
		`,
	}

	for _, query := range queries {
		if _, err := t.db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}
