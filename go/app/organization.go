package app

type Organization struct {
	ID            int    `db:"organization_id"`
	Name          string `db:"name"`
	CoversionRate int    `db:"coversion_rate"`
}
