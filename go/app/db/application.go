package db

import (
	"context"

	"github.com/BenJetson/CPSC491-project/go/app"
)

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
