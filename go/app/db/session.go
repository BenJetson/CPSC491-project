package db

import (
	"github.com/google/uuid"

	"github.com/BenJetson/CPSC491-project/go/app"
)

// GetSessionsForPerson fetches all sessions for a given person with matching ID.
func (db *database) GetSessionsForPerson(personID int) ([]app.Session, error) {
	return nil, nil // TODO
}

// GetSessionByToken fetches the session with matching token.
func (db *database) GetSessionByToken(token uuid.UUID) (app.Session, error) {
	return app.Session{}, nil // TODO
}

// CreateSession creates a new session, ignoring the ID field.
func (db *database) CreateSession(s app.Session) error {
	return nil // TODO
}

// RevokeSession revokes an existing session.
func (db *database) RevokeSession(token uuid.UUID) error {
	return nil // TODO
}

// RevokeSessionsForPersonExcept revokes all sessions for a person except the one
// with a matching token.
//
// This is useful in a password change scenario where you might want to invalidate
// all other login sessions except for the one where the user changed their
// password from.
func (db *database) RevokeSessionsForPersonExcept(personID int, token uuid.UUID) error {
	return nil // TODO
}
