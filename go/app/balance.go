package app

type Balance struct {
	PersonID          int    `db:"person_id" json:"person_id"`
	PersonFirstName   string `db:"first_name" json:"person_first_name"`
	PersonLastName    string `db:"last_name" json:"person_last_name"`
	OrganizationID    int    `db:"organization_id" json:"organization_id"`
	OrganizationTitle int    `db:"title" json:"organization_title"`
	Balance           int    `db:"points" json:"balance"`
}
