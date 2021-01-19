package demo

type Bar struct {
	ID   int    `json:"id" db:"bar_id"`
	Name string `json:"name" db:"name"`
}
