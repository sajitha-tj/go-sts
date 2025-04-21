package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ory/fosite"
)

const (
	AuthorizationCodeSessionType = "authorization_code"
	AccessTokenSessionType       = "access_token"
	RefreshTokenSessionType      = "refresh_token"
)

type AuthorizationCodeSession struct {
	ID          string
	Active      bool
	Code        string
	RequestedAt time.Time
	ClientID    string
}

type TokenSession struct {
	ID          string
	Active      bool
	Signature   string
	RequestedAt time.Time
	ClientID    string
}

type SessionStore struct {
	db *sql.DB
}

func NewSessionStore(db *sql.DB) *SessionStore {
	return &SessionStore{db: db}
}

func (ss *SessionStore) CreateSession(ctx context.Context, payload string, sessionType string, request fosite.Requester) error {
	var query string
	switch sessionType {
	case AuthorizationCodeSessionType:
		query = `
            INSERT INTO authorization_code_sessions (id, active, code, requested_at, session_id, client_id)
            VALUES ($1, $2, $3, $4, $5, $6)
        `
	case AccessTokenSessionType:
		query = `
            INSERT INTO access_token_sessions (id, active, signature, requested_at, session_id, client_id)
            VALUES ($1, $2, $3, $4, $5, $6)
        `
	case RefreshTokenSessionType:
		query = `
            INSERT INTO refresh_token_sessions (id, active, signature, requested_at, session_id, client_id)
            VALUES ($1, $2, $3, $4, $5, $6)
        `
	default:
		return fosite.ErrInvalidRequest
	}

	_, err := ss.db.ExecContext(
		ctx,
		query,
		request.GetID(),
		true, // Active
		payload,
		time.Now(),
		request.GetSession().GetUsername(),
		request.GetClient().GetID(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionStore) GetSession(ctx context.Context, payload string, sessionType string, session fosite.Session) (fosite.Requester, error) {
	var query string
	switch sessionType {
	case AuthorizationCodeSessionType:
		query = `
			SELECT id, active, code, requested_at, client_id
			FROM authorization_code_sessions
			WHERE code = $1 AND active = true
		`
	case AccessTokenSessionType:
		query = `
			SELECT id, active, signature, requested_at, client_id
			FROM access_token_sessions
			WHERE signature = $1 AND active = true
		`
	case RefreshTokenSessionType:
		query = `
			SELECT id, active, signature, requested_at, client_id
			FROM refresh_token_sessions
			WHERE signature = $1 AND active = true
		`
	default:
		return nil, fosite.ErrInvalidRequest
	}

	row := s.db.QueryRowContext(ctx, query, payload)
	var id, clientID string
	var active bool
	var requestedAt time.Time

	err := row.Scan(&id, &active, &payload, &requestedAt, &clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	request := fosite.NewRequest()
	request.ID = id
	request.Client = &fosite.DefaultClient{ID: clientID}
	request.Session = session
	request.RequestedAt = requestedAt

	return request, nil
}

func (s *SessionStore) InvalidateSession(ctx context.Context, payload string, sessionType string) error {
	var query string
	switch sessionType {
	case AuthorizationCodeSessionType:
		query = `
			UPDATE authorization_code_sessions
			SET active = false
			WHERE code = $1
		`
	case AccessTokenSessionType:
		query = `
			UPDATE access_token_sessions
			SET active = false
			WHERE signature = $1
		`
	case RefreshTokenSessionType:
		query = `
			UPDATE refresh_token_sessions
			SET active = false
			WHERE signature = $1
		`
	default:
		return fosite.ErrInvalidRequest
	}

	_, err := s.db.ExecContext(ctx, query, payload)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionStore) RotateRefreshToken(ctx context.Context, requestID string, refreshTokenSignature string) error {
	// Delete the current refresh token
	deleteQuery := `
		DELETE FROM refresh_token_sessions
		WHERE signature = $1
	`
	_, err := s.db.ExecContext(ctx, deleteQuery, refreshTokenSignature)
	if err != nil {
		return err
	}

	// Generate a new refresh token
	insertQuery := `
		INSERT INTO refresh_token_sessions (id, active, signature, requested_at, session_id, client_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err = s.db.ExecContext(
		ctx,
		insertQuery,
		requestID,
		true,                  // Active
		refreshTokenSignature, // Use the same signature or generate a new one as needed
		time.Now(),
		requestID, // Assuming session_id is the same as requestID
		"",        // Assuming client_id is not required here, replace with actual value if needed
	)
	if err != nil {
		return err
	}

	return nil
}
