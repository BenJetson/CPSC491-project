package db

import (
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
func (db *database) GetPersonByID(personID int) (app.Person, error) {
	return app.Person{}, nil // TODO
}

// GetPersonByEmail fetches a person given their email.
func (db *database) GetPersonByEmail(email string) (app.Person, error) {
	var dbp dbPerson

	err := db.Get(&dbp, `
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

	return dbp.toPerson(), errors.Wrap(err, "failed to get person")
}

// CreatePerson creates a new person given the details. Ignores the ID field.
func (db *database) CreatePerson(p app.Person) error {
	return nil // TODO
}

// UpdatePersonName updates a person's first and last name.
func (db *database) UpdatePersonName(personID int, firstName string, lastName string) error {
	return nil // TODO
}

// UpdatePersonRole updates a person's role.
func (db *database) UpdatePersonRole(personID int, roleType app.Role) error {
	return nil // TODO
}

// UpdatePersonPassword updates a person's password.
func (db *database) UpdatePersonPassword(personID int, p app.Password) error {
	return nil // TODO
}

// ActivatePerson activates a person's account.
func (db *database) ActivatePerson(personID int) error {
	return nil // TODO
}

// DeactivatePerson deactivates a person's account.
func (db *database) DeactivatePerson(personID int) error {
	return nil // TODO
}
