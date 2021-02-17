package app

// Bar stores arbitrary data, not very useful outside of the demo.
type Bar struct {
	ID   int    `json:"id" db:"bar_id"`
	Name string `json:"name" db:"name"`
}
