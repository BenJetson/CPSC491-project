package db

import "github.com/BenJetson/CPSC491-project/go/app"

// GetPersonByID fetches a person given their ID number.
func (db *database) GetPersonByID(personID int) (app.Person, error) {
	return app.Person{}, nil // TODO
}

// GetPersonByEmail fetches a person given their email.
func (db *database) GetPersonByEmail(email string) (app.Person, error) {
	return app.Person{}, nil // TODO
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
