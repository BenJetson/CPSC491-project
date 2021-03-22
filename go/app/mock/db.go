package mock

import (
	"context"

	"github.com/google/uuid"

	"github.com/BenJetson/CPSC491-project/go/app"
)

// This is an assertion, which will cause the build to fail if the mock.DB type
// does not implement the app.DataStore interface.
var _ app.DataStore = (*DB)(nil)

// DB is a mock of interface app.DataStore. It shall implement all of its
// methods, each of which does nothing.
//
// This can be useful for tests, where this type may be embedded in a more
// specific mock with overrided methods.
type DB struct{}

//
//
// PersonStore methods
//
//

// GetPersonByID mocks fetching a person by their ID.
func (db *DB) GetPersonByID(
	ctx context.Context,
	personID int,
) (app.Person, error) {

	return app.Person{}, nil
}

// GetPersonByEmail mocks fetching a person by their email address.
func (db *DB) GetPersonByEmail(
	ctx context.Context,
	email string,
) (app.Person, error) {

	return app.Person{}, nil
}

// CreatePerson mocks creating a new person.
func (db *DB) CreatePerson(ctx context.Context, p app.Person) (int, error) {
	return 0, nil
}

// UpdatePersonName mocks updating a person's name.
func (db *DB) UpdatePersonName(
	ctx context.Context,
	personID int,
	firstName, lastName string,
) error {

	return nil
}

// UpdatePersonRole mocks updating a person's role.
func (db *DB) UpdatePersonRole(
	ctx context.Context,
	personID int,
	roleType app.Role,
) error {

	return nil
}

// UpdatePersonPassword mocks updating a person's password.
func (db *DB) UpdatePersonPassword(
	ctx context.Context,
	personID int,
	p app.Password,
) error {

	return nil
}

// ActivatePerson mocks activating a person's account.
func (db *DB) ActivatePerson(
	ctx context.Context,
	personID int,
) error {

	return nil
}

// DeactivatePerson mocks deactivating a person's account.
func (db *DB) DeactivatePerson(
	ctx context.Context,
	personID int,
) error {

	return nil
}

//
//
// SessionStore methods
//
//

// GetSessionsForPerson mocks fetching all sessions for a person.
func (db *DB) GetSessionsForPerson(ctx context.Context,
	personID int,
	includeInvalid bool,
) ([]app.Session, error) {

	return nil, nil
}

// GetSessionByToken mocks fetching a session by its token.
func (db *DB) GetSessionByToken(
	ctx context.Context,
	token uuid.UUID,
) (app.Session, error) {

	return app.Session{}, nil
}

// CreateSession mocks creating a new session.
func (db *DB) CreateSession(ctx context.Context, s app.Session) (int, error) {
	return 0, nil
}

// RevokeSession mocks revoking a session.
func (db *DB) RevokeSession(ctx context.Context, sessionID int) error {
	return nil
}

// RevokeSessionsForPersonExcept mocks revoking all sessions for a person
// except for the session with matching token.
func (db *DB) RevokeSessionsForPersonExcept(
	ctx context.Context,
	personID, sessionID int,
) error {

	return nil
}

//
//
// ApplicationStore methods
//
//

// CreateApplication mocks creating an application in the database.
func (db *DB) CreateApplication(
	ctx context.Context,
	a app.Application,
) (int, error) {

	return 0, nil
}

// UpdateApplicationApproval mocks setting application approval.
func (db *DB) UpdateApplicationApproval(
	ctx context.Context,
	appID int,
	status bool,
	reason string,
) error {

	return nil
}
