package app

// An Organization contains information about a particular sponsor organization.
type Organization struct {
	// ID uniquely identifies this organization.
	ID int `db:"organization_id"`
	// Name is a human-readable alias for this organization.
	Name string `db:"name"`
	// PointValue describes the ratio between points and real dollars.
	// Each point is worth a PointValue amount of Money.
	PointValue Money `db:"point_value"`
}
