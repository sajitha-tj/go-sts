package client_repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type ClientStore struct {
	db *sql.DB
}

// NewClientStore initializes a new ClientStore with the given database connection.
func NewClientStore(db *sql.DB) *ClientStore {
	return &ClientStore{db: db}
}

// GetClient loads the client by its ID or returns an error if the client does not exist or another error occurred.
func (cs *ClientStore) GetClient(ctx context.Context, id string) (Client, error) {
	var client Client
	query := `
	SELECT id, secret, rotated_secrets, redirect_uris, grant_types, 
	response_types, scopes, public, audience 
	FROM clients WHERE id = $1`
	err := cs.db.QueryRowContext(ctx, query, id).Scan(
		&client.ClientID,
		&client.ClientSecret,
		&client.RotatedSecrets,
		&client.RedirectURIs,
		&client.GrantTypes,
		&client.ResponseTypes,
		&client.Scopes,
		&client.Public,
		&client.Audience,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Client{}, errors.New("client not found")
		}
		return Client{}, err
	}
	return client, nil
}

// ClientAssertionJWTValid returns an error if the JTI is known or the DB check failed and nil if the JTI is not known.
func (cs *ClientStore) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	query := "SELECT 1 FROM client_jtis WHERE jti = $1 AND expiry > $2"
	var exists int
	err := cs.db.QueryRowContext(ctx, query, jti, time.Now()).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil // JTI is not known
		}
		return err
	}
	return errors.New("JTI is already known")
}

// SetClientAssertionJWT marks a JTI as known for the given expiry time.
func (cs *ClientStore) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	// Clean up expired JTIs
	cleanupQuery := "DELETE FROM client_jtis WHERE expiry <= $1"
	_, err := cs.db.ExecContext(ctx, cleanupQuery, time.Now())
	if err != nil {
		return err
	}
	
	// Insert the new JTI
	insertQuery := "INSERT INTO client_jtis (jti, expiry) VALUES ($1, $2)"
	_, err = cs.db.ExecContext(ctx, insertQuery, jti, exp)
	return err
}


func (cs *ClientStore) CreateClient(ctx context.Context, client *Client) error {
	query := `
	INSERT INTO clients (id, secret, rotated_secrets, redirect_uris, grant_types, 
	response_types, scopes, public, audience) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := cs.db.ExecContext(ctx, query,
		client.ClientID,
		client.ClientSecret,
		client.RotatedSecrets,
		client.RedirectURIs,
		client.GrantTypes,
		client.ResponseTypes,
		client.Scopes,
		client.Public,
		client.Audience,
	)
	return err
}
