package db

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

func (db *database) GetAllOrganizations(
	ctx context.Context,
) ([]app.Organization, error) {

	var orgs []app.Organization

	err := db.SelectContext(ctx, &orgs, `
		SELECT
			organization_id,
			name,
			point_value
		FROM organization
		ORDER BY name ASC
	`)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, errors.Wrap(err, "failed to fetch orgs")
	}

	return orgs, nil
}

func (db *database) GetOrganizationByID(
	ctx context.Context,
	orgID int,
) (app.Organization, error) {

	var org app.Organization

	err := db.GetContext(ctx, &org, `
		SELECT
			organization_id,
			name,
			point_value
		FROM organization
		WHERE organization_id = $1
	`, orgID)

	if errors.Is(err, sql.ErrNoRows) {
		return app.Organization{}, errors.Wrapf(
			app.ErrNotFound,
			"no such organization by id of '%d'", orgID,
		)
	}

	return org, errors.Wrap(err, "failed to select organizations")
}
