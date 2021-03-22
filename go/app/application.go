package app

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// Application represents a driver application to be sponsored by an
// organization.
type Application struct {
	// ID uniquely identifies this application.
	ID int
	// ApplicantID is the person ID of the driver applying to be sponsored.
	ApplicantID int
	// OrganizationID is the ID of the organization the driver would like to be
	// sponsored by.
	OrganizationID int
	// Comment is a driver-supplied comment to go with their application.
	Comment string
	// Approved specifies whether or not the organization has approved this
	// application or not. Will be null if no decision has been made.
	Approved null.Bool
	// Reason is the reason the organization accepted/rejected this application.
	Reason string
	// CreatedAt is the timestamp of when the driver submitted this application.
	CreatedAt time.Time
	// ApprovedAt is the timestamp of when the organization approved or rejected
	// this application.
	ApprovedAt time.Time
}
