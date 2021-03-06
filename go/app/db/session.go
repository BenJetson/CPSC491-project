package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

type dbSession struct {
	dbPerson
	ID        int       `db:"session_id"`
	Token     uuid.UUID `db:"token"`
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
	IsRevoked bool      `db:"is_revoked"`
}

func (s *dbSession) toSession() app.Session {
	return app.Session{
		Person:    s.dbPerson.toPerson(),
		ID:        s.ID,
		Token:     s.Token,
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
		IsRevoked: s.IsRevoked,
	}
}

// GetSessionsForPerson fetches all sessions for a given person of matching ID.
func (db *database) GetSessionsForPerson(
	ctx context.Context,
	personID int,
	includeInvalid bool,
) ([]app.Session, error) {

	// WARNING: if you change the logic here, make sure it matches the IsValid
	// method logic of the app.Session type!

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
		WHERE s.person_id = $1
	`
	params := []interface{}{personID}

	if !includeInvalid {
		query += `
			AND s.is_revoked = FALSE
			AND p.is_deactivated = FALSE
		`

		now := time.Now().UTC().Round(time.Second)
		params = append(params, now)
		query += `
			AND $2::timestamptz >= s.created_at::timestamptz
			AND $2::timestamptz < s.expires_at::timestamptz
		`
	}

	query += `
		GROUP BY
			s.session_id,
			p.person_id,
			a.person_id
	`

	var dbss []dbSession
	err := db.SelectContext(ctx, &dbss, query, params...)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to select sessions")
	}

	ss := make([]app.Session, len(dbss))
	for idx, dbs := range dbss {
		ss[idx] = dbs.toSession()
	}

	return ss, nil
}

// GetSessionByToken fetches the session with matching token.
func (db *database) GetSessionByToken(
	ctx context.Context,
	token uuid.UUID,
) (app.Session, error) {

	var dbs dbSession

	err := db.GetContext(ctx, &dbs, `
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
		GROUP BY
			s.session_id,
			p.person_id,
			a.person_id
	`, token)

	if errors.Is(err, sql.ErrNoRows) {
		return app.Session{}, errors.Wrapf(
			app.ErrNotFound,
			"no such session by token '%s'", token,
		)
	}

	return dbs.toSession(), errors.Wrap(err, "failed to get session")
}

// CreateSession creates a new session, ignoring the ID field.
func (db *database) CreateSession(
	ctx context.Context,
	s app.Session,
) (int, error) {

	var id int
	err := db.GetContext(ctx, &id, `
		INSERT INTO session (
			token,
			person_id,
			created_at,
			expires_at
		) VALUES ($1, $2, $3, $4)
		RETURNING session_id
	`, s.Token, s.Person.ID, s.CreatedAt, s.ExpiresAt)

	return id, errors.Wrap(err, "failed to insert session")
}

// RevokeSession revokes an existing session.
func (db *database) RevokeSession(ctx context.Context, sessionID int) error {
	result, err := db.ExecContext(ctx, `
		UPDATE session SET
			is_revoked = TRUE
		WHERE session_id = $1
	`, sessionID)

	if err != nil {
		return errors.Wrap(err, "failed to revoke session")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to check result of revoke")
	} else if n != 1 {
		return errors.Wrapf(
			app.ErrNotFound,
			"no such session by id of %d", sessionID,
		)
	}

	return nil
}

// RevokeSessionsForPersonExcept revokes all sessions for a person except the
// one with a matching token.
//
// This is useful in a password change scenario where you might want to
// invalidate all other login sessions except for the one where the user changed
// their password from.
func (db *database) RevokeSessionsForPersonExcept(
	ctx context.Context,
	personID, sessionID int,
) error {

	_, err := db.ExecContext(ctx, `
		UPDATE session SET
			is_revoked = TRUE
		WHERE
			person_id = $1
			AND session_id != $2
	`, personID, sessionID)

	return errors.Wrap(err, "failed to revoke sessions for person")
}
