package session_repository

import "time"

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
