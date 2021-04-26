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

func (db *database) CreateOrganization(
	ctx context.Context,
	org app.Organization,
) (int, error) {

	var id int
	err := db.GetContext(ctx, &id, `
		INSERT INTO organization (
			name,
			point_value
		) VALUES ($1, $2)
		RETURNING organization_id
	`, org.Name, org.PointValue)

	return id, errors.Wrap(err, "failed to insert organization")
}

func (db *database) UpdateOrganization(
	ctx context.Context,
	org app.Organization,
) error {

	_, err := db.ExecContext(ctx, `
		UPDATE organization SET
			name = $1,
			point_value = $2
		WHERE organization_id = $3
	`, org.Name, org.PointValue, org.ID)

	return errors.Wrap(err, "failed to update organization")
}

func (db *database) DeleteOrganization(ctx context.Context, orgID int) error {
	_, err := db.ExecContext(ctx, `
		DELETE FROM organization
		WHERE organization_id = $1
	`, orgID)
	return errors.Wrap(err, "failed to delete organization")
}
