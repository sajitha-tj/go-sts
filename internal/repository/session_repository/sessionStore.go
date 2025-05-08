package session_repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ory/fosite"
)

const (
	AuthorizationCodeSessionType = "authorization_code"
	AccessTokenSessionType       = "access_token"
	RefreshTokenSessionType      = "refresh_token"
)

const (
	// Table names
	AuthorizationCodeSessionsTable = "authorization_code_sessions"
	AccessTokenSessionsTable       = "access_token_sessions"
	RefreshTokenSessionsTable      = "refresh_token_sessions"
	AuthorizeRequestTable          = "authorize_requests"
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
            INSERT INTO ` + AuthorizationCodeSessionsTable + ` (code, active, req_id, requested_at, request_data, client_id)
            VALUES ($1, $2, $3, $4, $5, $6)
        `
	case AccessTokenSessionType:
		query = `
            INSERT INTO ` + AccessTokenSessionsTable + ` (signature, active, req_id, requested_at, request_data, client_id)
            VALUES ($1, $2, $3, $4, $5, $6)
        `
	case RefreshTokenSessionType:
		query = `
            INSERT INTO ` + RefreshTokenSessionsTable + ` (signature, active, req_id, requested_at, request_data, client_id)
            VALUES ($1, $2, $3, $4, $5, $6)
        `
	default:
		return fosite.ErrInvalidRequest
	}

	requestData, e := serializeRequest(request.Sanitize([]string{}))
	if e != nil {
		return e
	}

	_, err := ss.db.ExecContext(
		ctx,
		query,
		payload,
		true, // Active
		request.GetID(),
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
			SELECT request_data
			FROM ` + AuthorizationCodeSessionsTable + `
			WHERE code = $1 AND active = true
		`
	case AccessTokenSessionType:
		query = `
			SELECT request_data
			FROM ` + AccessTokenSessionsTable + `
			WHERE signature = $1 AND active = true
		`
	case RefreshTokenSessionType:
		query = `
			SELECT request_data
			FROM ` + RefreshTokenSessionsTable + `
			WHERE signature = $1 AND active = true
		`
	default:
		return nil, fosite.ErrInvalidRequest
	}

	row := s.db.QueryRowContext(ctx, query, payload)
	var requestData string

	err := row.Scan(&requestData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	storedRequest := &StoredRequest{}
	if err := deserializeRequestData(requestData, storedRequest); err != nil {
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
			UPDATE ` + AuthorizationCodeSessionsTable + `
			SET active = false
			WHERE code = $1
		`
	case AccessTokenSessionType:
		query = `
			UPDATE ` + AccessTokenSessionsTable + `
			SET active = false
			WHERE signature = $1
		`
	case RefreshTokenSessionType:
		query = `
			UPDATE ` + RefreshTokenSessionsTable + `
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

func (s *SessionStore) GetAccessTokenSignatureFromReqId(ctx context.Context, requestID string) (string, error) {
	var query = `
		SELECT signature
		FROM ` + AccessTokenSessionsTable + `
		WHERE req_id = $1 AND active = true
	`
	row := s.db.QueryRowContext(ctx, query, requestID)
	var signature string
	err := row.Scan(&signature)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fosite.ErrNotFound
		}
		return "", err
	}
	return signature, nil
}

func (s *SessionStore) GetRefreshTokenSignatureFromReqId(ctx context.Context, requestID string) (string, error) {
	var query = `
		SELECT signature
		FROM ` + RefreshTokenSessionsTable + `
		WHERE req_id = $1 AND active = true
	`
	row := s.db.QueryRowContext(ctx, query, requestID)
	var signature string
	err := row.Scan(&signature)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fosite.ErrNotFound
		}
		return "", err
	}
	return signature, nil
}

func (s *SessionStore) CreateAuthorizeRequestSession(ctx context.Context, request fosite.Requester) (string, error) {
	query := `
		INSERT INTO ` + AuthorizeRequestTable + ` (req_id, request_data, authenticated, requested_at, exp_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (req_id) DO UPDATE SET
		request_data = EXCLUDED.request_data,
		requested_at = EXCLUDED.requested_at,
		exp_at = EXCLUDED.exp_at
	`
	requestData, err := serializeRequest(request)
	if err != nil {
		return "", err
	}
	_, err = s.db.ExecContext(
		ctx,
		query,
		request.GetID(),
		requestData,
		false, // Authenticated
		time.Now(),
		time.Now().Add(time.Minute*5),
	)

	if err != nil {
		return "", err
	}
	return request.GetID(), nil
}

func (s *SessionStore) GetAuthorizeRequestSession(ctx context.Context, requestID string) (fosite.AuthorizeRequester, error) {
	query := `
		SELECT request_data
		FROM ` + AuthorizeRequestTable + `
		WHERE req_id = $1 AND exp_at > $2
	`
	row := s.db.QueryRowContext(ctx, query, requestID, time.Now())
	var requestData string
	err := row.Scan(&requestData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fosite.ErrNotFound
		}
		return nil, err
	}

	storedRequest := &StoredRequest{}
	if err := deserializeRequestData(requestData, storedRequest); err != nil {
		return nil, err
	}

	request := &fosite.AuthorizeRequest{
		Request: fosite.Request{
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
		},
		ResponseTypes:        storedRequest.ResponseTypes,
		RedirectURI:          storedRequest.RedirectURI,
		State:                storedRequest.State,
		HandledResponseTypes: storedRequest.HandledResponseTypes,
		ResponseMode:         storedRequest.ResponseMode,
		DefaultResponseMode:  storedRequest.DefaultResponseMode,
	}
	return request, nil
}

func (s *SessionStore) AuthenticateAuthorizeRequestSession(ctx context.Context, requestID string) error {
	query := `
		UPDATE ` + AuthorizeRequestTable + `
		SET authenticated = true
		WHERE req_id = $1
	`
	_, err := s.db.ExecContext(ctx, query, requestID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionStore) IsRequestSessionAuthenticated(ctx context.Context, requestID string) (bool, error) {
	query := `
		SELECT authenticated
		FROM ` + AuthorizeRequestTable + `
		WHERE req_id = $1
	`
	row := s.db.QueryRowContext(ctx, query, requestID)
	var authenticated bool
	err := row.Scan(&authenticated)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fosite.ErrNotFound
		}
		return false, err
	}
	return authenticated, nil
}

// serializeRequest serializes the request data into a JSON string.
func serializeRequest(req fosite.Requester) (string, error) {
	reqData, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	return string(reqData), nil
}

// deserializeRequestData deserializes the JSON string into the StoredRequest struct.
func deserializeRequestData(reqData string, req *StoredRequest) error {
	return json.Unmarshal([]byte(reqData), &req)
}
