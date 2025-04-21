package storage

import (
	"context"
	"time"

	"github.com/ory/fosite"
)

func (s Storage) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	return s.clientStore.GetClient(ctx, id)
}

func (s Storage) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	return s.clientStore.ClientAssertionJWTValid(ctx, jti)
}

func (s Storage) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	return s.clientStore.SetClientAssertionJWT(ctx, jti, exp)
}
