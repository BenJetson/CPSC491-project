package app

import (
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
	GetPersonByID(personID int) (Person, error)
	GetPersonByEmail(email string) (Person, error)

	CreatePerson(p Person) error

	UpdatePersonName(personID int, firstName, lastName string) error
	UpdatePersonRole(personID int, roleType Role) error
	UpdatePersonPassword(personID int, p Password) error
	ActivatePerson(personID int) error
	DeactivatePerson(personID int) error
}

// SessionStore defines methods for working with app.Session objects in the
// database.
type SessionStore interface {
	GetSessionsForPerson(personID int, includeInvalid bool) ([]Session, error)
	GetSessionByToken(token uuid.UUID) (Session, error)

	CreateSession(s Session) error

	RevokeSession(sessionID int) error
	RevokeSessionsForPersonExcept(personID, sessionID int) error
}
