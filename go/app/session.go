package app

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const SessionLength = 6 * time.Hour

type Session struct {
	ID        int       `db:"session_id"`
	Token     uuid.UUID `db:"token"`
	Person    Person
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
	IsRevoked bool      `db:"is_revoked"`
}

func NewSession(p Person) (*Session, error) {
	now := time.Now()

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

func (s *Session) IsValid() bool {
	now := time.Now()

	return !s.IsRevoked &&
		!s.Person.IsDeactivated &&
		now.After(s.CreatedAt) &&
		now.Before(s.ExpiresAt)
}
