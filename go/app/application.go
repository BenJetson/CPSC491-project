package app

import "gopkg.in/guregu/null.v4"

// Application represents a driver application to be sponsored by an
// organization.
type Application struct {
	ID             int
	ApplicantID    int
	OrganizationID int
	Comment        string
	Approved       null.Bool
	Reason         string
}
