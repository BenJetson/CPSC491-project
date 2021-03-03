package app

import "github.com/google/uuid"

// DataStore is the common interface for durable data storage.
type DataStore interface {
	PersonStore
	SessionStore
}

type PersonStore interface {
	GetPersonByID(personID int) (Person, error)
	GetPersonByEmail(email string) (Person, error)

	CreatePerson(p Person) error

	UpdatePersonName(personID int, firstName string, lastName string) error
	UpdatePersonRole(personID int, roleType Role) error
	UpdatePersonPassword(personID int, p Password) error
	ActivatePerson(personID int) error
	DeactivatePerson(personID int) error
}

type SessionStore interface {
	GetSessionsForPerson(personID int) ([]Session, error)
	GetSessionByToken(token uuid.UUID) (Session, error)

	CreateSession(s Session) error

	RevokeSession(token uuid.UUID) error
	RevokeSessionsForPersonExcept(personID int, token uuid.UUID) error
}
