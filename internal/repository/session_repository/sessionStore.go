package session_repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/ory/fosite"
)

const (
	AuthorizationCodeSessionType = "authorization_code"
	AccessTokenSessionType       = "access_token"
	RefreshTokenSessionType      = "refresh_token"
)

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
            INSERT INTO authorization_code_sessions (id, active, code, requested_at, request_data, client_id)
            VALUES ($1, $2, $3, $4, $5, $6)
        `
	case AccessTokenSessionType:
		query = `
            INSERT INTO access_token_sessions (id, active, signature, requested_at, request_data, client_id)
            VALUES ($1, $2, $3, $4, $5, $6)
        `
	case RefreshTokenSessionType:
		query = `
            INSERT INTO refresh_token_sessions (id, active, signature, requested_at, request_data, client_id)
            VALUES ($1, $2, $3, $4, $5, $6)
        `
	default:
		return fosite.ErrInvalidRequest
	}

	requestData, e := getSerializedRequest(request)
	if e != nil {
		log.Println("Error serializing session data:", e)
		return e
	}

	_, err := ss.db.ExecContext(
		ctx,
		query,
		request.GetID(),
		true, // Active
		payload,
		time.Now(),
		requestData,
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
			SELECT id, active, code, requested_at, request_data, client_id
			FROM authorization_code_sessions
			WHERE code = $1 AND active = true
		`
	case AccessTokenSessionType:
		query = `
			SELECT id, active, signature, requested_at, request_data, client_id
			FROM access_token_sessions
			WHERE signature = $1 AND active = true
		`
	case RefreshTokenSessionType:
		query = `
			SELECT id, active, signature, requested_at, request_data, client_id
			FROM refresh_token_sessions
			WHERE signature = $1 AND active = true
		`
	default:
		return nil, fosite.ErrInvalidRequest
	}

	row := s.db.QueryRowContext(ctx, query, payload)
	var id, clientID, requestData string
	var active bool
	var requestedAt time.Time

	err := row.Scan(&id, &active, &payload, &requestedAt, &requestData, &clientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	storedRequest := &StoredRequest{}
	if err := deserializeRequestData(requestData, storedRequest); err != nil {
		log.Println("Error deserializing session data:", err)
		return nil, err
	}

	request := &fosite.Request{
		ID:                storedRequest.ID,
		RequestedAt:       storedRequest.RequestedAt,
		Client:            storedRequest.Client,
		RequestedScope:    storedRequest.RequestedScope,
		GrantedScope:      storedRequest.GrantedScope,
		Form:              storedRequest.Form,
		Session:           &storedRequest.Session,
		RequestedAudience: storedRequest.RequestedAudience,
		GrantedAudience:   storedRequest.GrantedAudience,
		Lang:              storedRequest.Lang,
	}
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
	// Implement logic to rotate the refresh token
	return nil
}

func getSerializedRequest(req any) (string, error) {
	reqData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	return string(reqData), nil
}

func deserializeRequestData(reqData string, req any) error {
	return json.Unmarshal([]byte(reqData), &req)
}
