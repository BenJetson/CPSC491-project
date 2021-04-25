package app

// A Balance describes the number of points that a person has within a
// given Organization.
type Balance struct {
	PersonID          int    `db:"person_id" json:"person_id"`
	PersonFirstName   string `db:"first_name" json:"person_first_name"`
	PersonLastName    string `db:"last_name" json:"person_last_name"`
	OrganizationID    int    `db:"organization_id" json:"organization_id"`
	OrganizationTitle int    `db:"title" json:"organization_title"`
	Points            int    `db:"points" json:"balance"`
}
