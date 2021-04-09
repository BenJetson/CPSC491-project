package app

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// SessionLength defines how long a login session for this app will last, at
// a MAXIMUM. Sessions may be revoked prior.
const SessionLength = 6 * time.Hour

// A Session represents an individual user session with our app.
type Session struct {
	// Person holds the user that this session belongs to.
	Person
	// ID is the session's identifying number. It is not secret.
	ID int `db:"session_id"`
	// Token is a unique, securely random generated UUIDv4 that also uniquely
	// identifies this session. It must be kept secret between our API server
	// and the web browser for the session.
	Token uuid.UUID `db:"token"`
	// CreatedAt is the timestamp that this session was started at.
	CreatedAt time.Time `db:"created_at"`
	// ExpiresAt is the timestamp when this session will expire permanently.
	ExpiresAt time.Time `db:"expires_at"`
	// IsRevoked is true when the session was manually revoked.
	IsRevoked bool `db:"is_revoked"`
}

// NewSession creates a new login session with a secure random token for a given
// person. It shall expire after SessionLength time has passed.
func NewSession(p Person) (*Session, error) {
	now := time.Now().UTC().Round(time.Second)

	token, err := uuid.NewRandom()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create random UUID")
	}

	return &Session{
		Token:     token,
		Person:    p,
		CreatedAt: now,
		ExpiresAt: now.Add(SessionLength),
	}, nil
}

// IsValid determines whether or not a session is still valid.
func (s *Session) IsValid() bool {
	now := time.Now().UTC().Round(time.Second)

	return !s.IsRevoked &&
		!s.Person.IsDeactivated &&
		(now.After(s.CreatedAt) || now.Equal(s.CreatedAt)) &&
		now.Before(s.ExpiresAt)
}
