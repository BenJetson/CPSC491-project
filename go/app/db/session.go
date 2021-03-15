package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

// GetSessionsForPerson fetches all sessions for a given person with matching ID.
func (db *database) GetSessionsForPerson(
	personID int,
	includeInvalid bool,
) ([]app.Session, error) {

	// WARNING: if you change the logic here, make sure it matches the IsValid
	// method logic of the app.Session type!

	var params []interface{}
	query := `
		SELECT
			s.session_id,
			s.token,
			s.created_at,
			s.expires_at,
			s.is_revoked,
			p.person_id,
			p.first_name,
			p.last_name,
			p.email,
			p.role_id,
			p.pass_hash,
			p.is_deactivated,
			array_remove(array_agg(a.organization_id), NULL) as affiliations
		FROM session s
		JOIN person p
			ON s.person_id = p.person_id
		LEFT JOIN affiliation a
			ON p.person_id = a.person_id
	`

	if !includeInvalid {
		now := time.Now()

		query += `
			WHERE
				s.is_revoked = FALSE
				AND p.is_deactivated = FALSE
				AND $1::timestamptz > s.created_at
				AND $1::timestamptz < s.expires_at
		`

		params = append(params, now)
	}

	query += `
		GROUP BY s.person_id
	`

	var ss []app.Session
	err := db.Select(ss, query, params...)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to select sessions")
	}

	return ss, nil
}

// GetSessionByToken fetches the session with matching token.
func (db *database) GetSessionByToken(token uuid.UUID) (app.Session, error) {
	var s app.Session

	err := db.Get(&s, `
		SELECT
			s.session_id,
			s.token,
			s.created_at,
			s.expires_at,
			s.is_revoked,
			p.person_id,
			p.first_name,
			p.last_name,
			p.email,
			p.role_id,
			p.pass_hash,
			p.is_deactivated,
			array_remove(array_agg(a.organization_id), NULL) as affiliations
		FROM session s
		JOIN person p
			ON s.person_id = p.person_id
		LEFT JOIN affiliation a
			ON p.person_id = a.person_id
		WHERE s.token = $1
		GROUP BY s.person_id
	`, token)

	if errors.Is(err, sql.ErrNoRows) {
		return app.Session{}, errors.Wrapf(
			app.ErrNotFound,
			"no such session by token '%s'", token,
		)
	}

	return s, errors.Wrap(err, "failed to get session")
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
		return errors.Errorf("insert session ought to affect 1 row, found: %d", n)
	}

	return nil
}

// RevokeSession revokes an existing session.
func (db *database) RevokeSession(sessionID int) error {
	_, err := db.Exec(`
		UPDATE session SET
			revoked = TRUE
		WHERE session_id = $1
	`, sessionID)

	return errors.Wrap(err, "failed to revoke session")
}

// RevokeSessionsForPersonExcept revokes all sessions for a person except the one
// with a matching token.
//
// This is useful in a password change scenario where you might want to invalidate
// all other login sessions except for the one where the user changed their
// password from.
func (db *database) RevokeSessionsForPersonExcept(personID int, sessionID int) error {
	_, err := db.Exec(`
		UPDATE session SET
			revoked = TRUE
		WHERE
			person_id = $1
			AND session_id != $2
	`, personID, sessionID)

	return errors.Wrap(err, "failed to revoke sessions for person")
}
