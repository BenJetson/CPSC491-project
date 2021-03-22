package app

// A Role specifies what role a user has in our app and therefore what
// permissions they might have.
type Role int

const (
	// RoleAdmin users can access all parts of the system, even content that is
	// not owned by their user.
	RoleAdmin Role = 1
	// RoleSponsor can access sponsor control panels for their organization as
	// well as all driver information for their affiliated drivers.
	RoleSponsor Role = 2
	// RoleUser can access the application system only, since they are not
	// affiliated with any organization.
	RoleUser Role = 3
	// RoleDriver can access the ordering system and view their point balance,
	// as they are affiliated with a sponsor organization.
	RoleDriver Role = 4
)
