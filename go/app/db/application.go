package db

import (
	"context"

	"github.com/BenJetson/CPSC491-project/go/app"
)

// GetApplicationByID fetches an application by its ID number.
func (db *database) GetApplicationByID(
	ctx context.Context,
	appID int,
) (app.Application, error) {

	return app.Application{}, nil // TODO
}

// GetApplicationsForPerson fetches all applications submitted by a person.
func (db *database) GetApplicationsForPerson(
	ctx context.Context,
	personID int,
) ([]app.Application, error) {

	return nil, nil // TODO
}

// GetApplicationsForOrganization fetches all applications submitted for
// an organization.
func (db *database) GetApplicationsForOrganization(
	ctx context.Context,
	orgID int,
) ([]app.Application, error) {

	return nil, nil // TODO
}

// CreateApplication creates a new application in the database.
func (db *database) CreateApplication(
	ctx context.Context,
	a app.Application,
) (int, error) {

	return 0, nil // TODO
}

// UpdateApplicationApproval sets the application approval status in the
// database.
func (db *database) UpdateApplicationApproval(
	ctx context.Context,
	appID int,
	status bool,
	reason string,
) error {

	return nil // TODO
}
