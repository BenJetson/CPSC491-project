package app

// Points describes a particular amount of points, the value of each point,
// and the organization those points are associated with.
type Points struct {
	// Amount describes the quantity of points.
	Amount int `json:"amount"`
	// OrganizationID is the identifier of the organization that the points
	// originated from.
	OrganizationID int `json:"organization_id"`
	// PointValue describes the amount of Money that each point is worth.
	PointValue Money `json:"point_value"`
}
