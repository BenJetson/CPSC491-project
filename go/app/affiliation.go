package app

import (
	"gopkg.in/guregu/null.v4"
)

type Affiliation struct {
	PersonID       int      `db:"person_id"`
	OrganizationID int      `db:"organization_id"`
	Points         null.Int `db:"points"`
}
