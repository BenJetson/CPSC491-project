package app

type Person struct {
	ID            int      `db:"person_id"`
	FirstName     string   `db:"first_name"`
	LastName      string   `db:"last_name"`
	Email         string   `db:"email"`
	Role          Role     `db:"role_id"`
	Password      Password `db:"pass_hash"`
	IsDeactivated bool     `db:"is_deactivated"`
}
