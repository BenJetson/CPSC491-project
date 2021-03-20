package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

type dbPerson struct {
	ID            int          `db:"person_id"`
	FirstName     string       `db:"first_name"`
	LastName      string       `db:"last_name"`
	Email         string       `db:"email"`
	Role          app.Role     `db:"role_id"`
	Password      app.Password `db:"pass_hash"`
	IsDeactivated bool         `db:"is_deactivated"`
	Affiliations  pq.Int64Array
}

func (p *dbPerson) toPerson() app.Person {
	out := app.Person{
		ID:            p.ID,
		FirstName:     p.FirstName,
		LastName:      p.LastName,
		Email:         p.Email,
		Role:          p.Role,
		Password:      p.Password,
		IsDeactivated: p.IsDeactivated,
		Affiliations:  make([]int, len(p.Affiliations)),
	}

	for i := range p.Affiliations {
		out.Affiliations[i] = int(p.Affiliations[i])
	}

	return out
}

// GetPersonByID fetches a person given their ID number.
func (db *database) GetPersonByID(
	ctx context.Context,
	personID int,
) (app.Person, error) {

	return app.Person{}, nil // TODO
}

// GetPersonByEmail fetches a person given their email.
func (db *database) GetPersonByEmail(
	ctx context.Context,
	email string,
) (app.Person, error) {

	var dbp dbPerson

	err := db.GetContext(ctx, &dbp, `
		SELECT
			p.person_id,
			p.first_name,
			p.last_name,
			p.email,
			p.role_id,
			p.pass_hash,
			p.is_deactivated,
			array_remove(array_agg(a.organization_id), NULL) as affiliations
		FROM person p
		LEFT JOIN affiliation a
			ON p.person_id = a.person_id
		WHERE email = $1
		GROUP BY p.person_id
	`, email)

	if errors.Is(err, sql.ErrNoRows) {
		return app.Person{}, errors.Wrapf(
			app.ErrNotFound,
			"no such person by email of '%s'", email,
		)
	}

	return dbp.toPerson(), errors.Wrap(err, "failed to get person")
}

// CreatePerson creates a new person given the details. Ignores the ID field.
func (db *database) CreatePerson(ctx context.Context, p app.Person) error {
	result, err := db.ExecContext(ctx, `
		INSERT INTO person (
			first_name,
			last_name,
			email,
			role_id,
			pass_hash
		) VALUES ($1, $2, $3, $4, $5)
	`, p.FirstName, p.LastName, p.Email, p.Role, p.Password)

	if err != nil {
		return errors.Wrap(err, "failed to insert person")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to check result of person insert")
	} else if n != 1 {
		return errors.Errorf(
			"insert person ought to affect 1 row, found: %d", n)
	}

	return nil
}

// UpdatePersonName updates a person's first and last name.
func (db *database) UpdatePersonName(
	ctx context.Context,
	personID int,
	firstName, lastName string,
) error {

	result, err := db.ExecContext(ctx, `
		UPDATE person SET
			first_name = $1,
			last_name = $2
		WHERE person_id = $3
	`, firstName, lastName, personID)

	if err != nil {
		return errors.Wrap(err, "failed to update person name")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to check result of person name update")
	} else if n != 1 {
		return errors.Wrapf(
			app.ErrNotFound,
			"no such person by id of %d", personID,
		)
	}

	return nil
}

// UpdatePersonRole updates a person's role.
func (db *database) UpdatePersonRole(
	ctx context.Context,
	personID int,
	roleType app.Role,
) error {

	result, err := db.ExecContext(ctx, `
		UPDATE person SET
			role_id = $1
		WHERE person_id = $2
	`, roleType, personID)

	if err != nil {
		return errors.Wrap(err, "failed to update person role")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to check result of person role update")
	} else if n != 1 {
		return errors.Wrapf(
			app.ErrNotFound,
			"no such person by id of %d", personID,
		)
	}

	return nil
}

// UpdatePersonPassword updates a person's password.
func (db *database) UpdatePersonPassword(
	ctx context.Context,
	personID int,
	newPass app.Password,
) error {

	result, err := db.ExecContext(ctx, `
		UPDATE person SET
			pass_hash = $1
		WHERE person_id = $2
	`, newPass, personID)

	if err != nil {
		return errors.Wrap(err, "failed to update person password")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err,
			"failed to check result of person password update")
	} else if n != 1 {
		return errors.Wrapf(
			app.ErrNotFound,
			"no such person by id of %d", personID,
		)
	}

	return nil
}

// ActivatePerson activates a person's account.
func (db *database) ActivatePerson(ctx context.Context, personID int) error {

	result, err := db.ExecContext(ctx, `
		UPDATE person SET
			is_deactivated = FALSE
		WHERE person_id = $1
	`, personID)

	if err != nil {
		return errors.Wrap(err, "failed to activate person")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to check result of person activation")
	} else if n != 1 {
		return errors.Wrapf(
			app.ErrNotFound,
			"no such person by id of %d", personID,
		)
	}

	return nil
}

// DeactivatePerson deactivates a person's account.
func (db *database) DeactivatePerson(ctx context.Context, personID int) error {
	result, err := db.ExecContext(ctx, `
		UPDATE person SET
			is_deactivated = TRUE
		WHERE person_id = $1
	`, personID)

	if err != nil {
		return errors.Wrap(err, "failed to deactivate person")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "failed to check result of person deactivation")
	} else if n != 1 {
		return errors.Wrapf(
			app.ErrNotFound,
			"no such person by id of %d", personID,
		)
	}

	return nil
}
