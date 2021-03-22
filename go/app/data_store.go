package app

import (
	"context"

	"github.com/google/uuid"
	// "gopkg.in/guregu/null.v4"
)

// DataStore is the common interface for durable data storage.
type DataStore interface {
	PersonStore
	SessionStore
}

// PersonStore defines methods for working with app.Person objects in the
// database.
type PersonStore interface {
	GetPersonByID(
		ctx context.Context,
		personID int,
	) (Person, error)
	GetPersonByEmail(ctx context.Context, email string) (Person, error)

	CreatePerson(ctx context.Context, p Person) error

	UpdatePersonName(
		ctx context.Context,
		personID int,
		firstName, lastName string,
	) error
	UpdatePersonRole(ctx context.Context, personID int, roleType Role) error
	UpdatePersonPassword(ctx context.Context, personID int, p Password) error
	ActivatePerson(ctx context.Context, personID int) error
	DeactivatePerson(ctx context.Context, personID int) error
}

// SessionStore defines methods for working with app.Session objects in the
// database.
type SessionStore interface {
	GetSessionsForPerson(
		ctx context.Context,
		personID int,
		includeInvalid bool,
	) ([]Session, error)
	GetSessionByToken(ctx context.Context, token uuid.UUID) (Session, error)

	CreateSession(ctx context.Context, s Session) error

	RevokeSession(ctx context.Context, sessionID int) error
	RevokeSessionsForPersonExcept(
		ctx context.Context,
		personID, sessionID int,
	) error
}
