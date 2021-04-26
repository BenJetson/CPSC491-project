package app

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

// Application represents a driver application to be sponsored by an
// organization.
type Application struct {
	// ID uniquely identifies this application.
	ID int `db:"application_id" json:"id"`
	// ApplicantID is the person ID of the driver applying to be sponsored.
	ApplicantID int `db:"applicant_id" json:"applicant_id"`

	// OrganizationID is the ID of the organization the driver would like to be
	// sponsored by.
	OrganizationID    int    `db:"organization_id" json:"organization_id"`
	OrganizationTitle string `db:"name" json:"organization_name"`
	// Comment is a driver-supplied comment to go with their application.
	Comment string `db:"comment" json:"comment"`
	// Approved specifies whether or not the organization has approved this
	// application or not. Will be null if no decision has been made.
	Approved null.Bool `db:"approved" json:"approved"`
	// Reason is the reason the organization accepted/rejected this application.
	Reason null.String `db:"reason" json:"reason"`
	// CreatedAt is the timestamp of when the driver submitted this application.
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	// ApprovedAt is the timestamp of when the organization approved or rejected
	// this application.
	ApprovedAt null.Time `db:"approved_at" json:"approved_at"`
}
