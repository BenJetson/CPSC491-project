package app

// A Person represents a user of our app.
type Person struct {
	// ID is the uniquely identifying number for this person's account.
	ID int `db:"person_id" json:"-"`
	// FirstName is the person's first name.
	FirstName string `db:"first_name" json:"first_name"`
	// LastName is the person's last name.
	LastName string `db:"last_name" json:"last_name"`
	// Email is the person's email address. May be changed.
	Email string `db:"email" json:"email"`
	// Role is the person's current role. May be changed.
	Role Role `db:"role_id" json:"role_id"`
	// Password is the person's current hashed password.
	Password Password `db:"pass_hash" json:"-"`
	// IsDeactivated is true when a person's account is deactivated and
	// therefore cannot be authenticated against.
	IsDeactivated bool `db:"is_deactivated" json:"-"`
	// Affiliations is a list of organization IDs that this user is
	// associated with.
	Affiliations []int `json:"affiliations"`
}
