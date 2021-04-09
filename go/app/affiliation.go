package app

import (
	"gopkg.in/guregu/null.v4"
)

// An Affiliation describes a person's relationship with an organization.
type Affiliation struct {
	PersonID       int `db:"person_id"`
	OrganizationID int `db:"organization_id"`
	// Points is the quantity of points this person has with this Organization,
	// will be null for non-drivers.
	Points null.Int `db:"points"`
}
