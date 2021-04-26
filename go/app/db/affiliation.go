package db

import (
	"context"

	"gopkg.in/guregu/null.v4"

	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

func (db *database) AddPersonAffiliation(
	ctx context.Context,
	personID, orgID int,
	role app.Role,
) error {

	var points null.Int
	if role == app.RoleDriver {
		points = null.IntFrom(0)
	}

	_, err := db.ExecContext(ctx, `
		INSERT INTO affiliation (
			person_id,
			organization_id,
			points
		) VALUES ($1, $2, $3)
		ON CONFLICT (person_id, organization_id)
		DO UPDATE SET points = $3
	`, personID, orgID, points)

	return errors.Wrap(err, "failed to add affiliation")
}

func (db *database) RemovePersonAffiliation(
	ctx context.Context,
	personID, orgID int,
) error {

	_, err := db.ExecContext(ctx, `
		DELETE FROM affiliation
		WHERE
			person_id = $1
			AND organization_id = $2
	`, personID, orgID)

	return errors.Wrap(err, "failed to delete affiliation")
}

func (db *database) SetPointsForAffiliation(
	ctx context.Context,
	personID, orgID int,
	points null.Int,
) error {

	_, err := db.ExecContext(ctx, `
		UPDATE affiliation SET
			points = $1
		WHERE
			person_id = $2
			AND organization
	`, points, personID, orgID)

	return errors.Wrap(err, "failed to update points for affiliation")
}

func (db *database) GetBalancesForPerson(
	ctx context.Context,
	personID int,
) ([]app.Balance, error) {

	var bs []app.Balance

	err := db.SelectContext(ctx, &bs, `
		SELECT
			p.person_id,
			p.first_name,
			p.last_name,
			o.organization_id,
			o.name,
			a.points
		FROM affiliation a
		JOIN person p
			ON a.person_id = p.person_id
		JOIN organization o
			ON a.organization_id = o.organization_id
		WHERE a.person_id = $1
	`, personID)

	return bs, errors.Wrap(err, "failed to retrieve balances")
}
