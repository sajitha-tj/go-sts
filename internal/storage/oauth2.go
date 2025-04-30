package storage

import (
	"context"

	"github.com/ory/fosite"
	"github.com/sajitha-tj/go-sts/internal/repository/session_repository"
)

// Implement AuthorizeCodeStorage methods
func (s *Storage) CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) error {
	return s.sessionStore.CreateSession(ctx, code, session_repository.AuthorizationCodeSessionType, request)
}

func (s *Storage) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (fosite.Requester, error) {
	return s.sessionStore.GetSession(ctx, code, session_repository.AuthorizationCodeSessionType, session)
}

func (s *Storage) InvalidateAuthorizeCodeSession(ctx context.Context, code string) error {
	return s.sessionStore.InvalidateSession(ctx, code, session_repository.AuthorizationCodeSessionType)
}

// Implement AccessTokenStorage methods
func (s *Storage) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) error {
	// Implement logic to store the access token session
	return s.sessionStore.CreateSession(ctx, signature, session_repository.AccessTokenSessionType, request)
}

func (s *Storage) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	// Implement logic to retrieve the access token session
	return s.sessionStore.GetSession(ctx, signature, session_repository.AccessTokenSessionType, session)
}

func (s *Storage) DeleteAccessTokenSession(ctx context.Context, signature string) error {
	// Implement logic to delete the access token session
	return s.sessionStore.InvalidateSession(ctx, signature, session_repository.AccessTokenSessionType)
}

// Implement RefreshTokenStorage methods
func (s *Storage) CreateRefreshTokenSession(ctx context.Context, signature string, accessSignature string, request fosite.Requester) error {
	// Implement logic to store the refresh token session
	return s.sessionStore.CreateSession(ctx, signature, session_repository.RefreshTokenSessionType, request)
}

func (s *Storage) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	// Implement logic to retrieve the refresh token session
	return s.sessionStore.GetSession(ctx, signature, session_repository.RefreshTokenSessionType, session)
}

func (s *Storage) DeleteRefreshTokenSession(ctx context.Context, signature string) error {
	// Implement logic to delete the refresh token session
	return s.sessionStore.InvalidateSession(ctx, signature, session_repository.RefreshTokenSessionType)
}

func (s *Storage) RotateRefreshToken(ctx context.Context, requestID string, refreshTokenSignature string) error {
	// Implement logic to rotate the refresh token
	if err := s.RevokeRefreshToken(ctx, requestID); err != nil {
		return err
	}
	return s.RevokeAccessToken(ctx, requestID)
}

func (s *Storage) RevokeRefreshToken(ctx context.Context, requestID string) error {
	// Implement logic to revoke the refresh token
	signature, err := s.sessionStore.GetRefreshTokenSignatureFromReqId(ctx, requestID)
	if err != nil {
		return err
	}
	return s.DeleteRefreshTokenSession(ctx, signature)
}

func (s *Storage) RevokeAccessToken(ctx context.Context, requestID string) error {
	// Implement logic to revoke the access token
	signature, err := s.sessionStore.GetAccessTokenSignatureFromReqId(ctx, requestID)
	if err != nil {
		return err
	}
	return s.DeleteAccessTokenSession(ctx, signature)
}
