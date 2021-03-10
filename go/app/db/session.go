package db

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

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
	result, err := db.Exec(`
		INSERT INTO session (
			token,
			person_id,
			created_at,
			expires_at
		) VALUES ($1, $2, $3, $4)
	`, s.Token, s.Person.ID, s.CreatedAt, s.ExpiresAt)

	if err != nil {
		return errors.Wrap(err, "failed to insert session")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to check result of session insert")
	} else if n != 1 {
		return errors.Errorf(" insert session ought to affect 1 row, found: %d", n)
	}

	return nil
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
