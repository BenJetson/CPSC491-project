package db

import (
	"context"

	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

// GetApplicationByID fetches an application by its ID number.
func (db *database) GetApplicationByID(
	ctx context.Context,
	appID int,
) (app.Application, error) {

	var a app.Application

	err := db.GetContext(ctx, &a, `
		SELECT
			a.application_id,
			a.applicant_id,
			a.organization_id,
			a.comment,
			a.approved,
			a.reason,
			a.created_at,
			a.approved_at
		FROM application a
		WHERE application_id = $1
	`, appID)

	return a, errors.Wrap(err, "failed to get application by ID")
}

// GetApplicationsForPerson fetches all applications submitted by a person.
func (db *database) GetApplicationsForPerson(
	ctx context.Context,
	personID int,
) ([]app.Application, error) {

	var apps []app.Application

	err := db.SelectContext(ctx, &apps, `
		SELECT
			a.application_id,
			a.applicant_id,
			a.organization_id,
			a.comment,
			a.approved,
			a.reason,
			a.created_at,
			a.approved_at
		FROM application a
		WHERE applicant_id = $1
	`, personID)

	return apps, errors.Wrap(err, "failed to select application for person")
}

// GetApplicationsForOrganization fetches all applications submitted for
// an organization.
func (db *database) GetApplicationsForOrganization(
	ctx context.Context,
	orgID int,
) ([]app.Application, error) {

	var apps []app.Application

	err := db.SelectContext(ctx, &apps, `
		SELECT
			a.application_id
			a.applicant_id
			a.organization_id
			a.comment
			a.approved
			a.reason
			a.created_at
			a.approved_at
		FROM application a
		WHERE organization_id = $1
	`, orgID)

	return apps, errors.Wrap(err,
		"failed to select application for organization")
}

// CreateApplication creates a new application in the database.
func (db *database) CreateApplication(
	ctx context.Context,
	a app.Application,
) (int, error) {

	var id int
	err := db.GetContext(ctx, &id, `
		INSERT INTO application (
			applicant_id,
			organization_id,
			comment,
		) VALUES ($1, $2, $3)
		RETURNING application_id
	`, a.ApplicantID, a.OrganizationID, a.Comment)

	return id, errors.Wrap(err, "failed to insert application")
}

// UpdateApplicationApproval sets the application approval status in the
// database.
func (db *database) UpdateApplicationApproval(
	ctx context.Context,
	appID int,
	status bool,
	reason string,
) error {

	result, err := db.ExecContext(ctx, `
		UPDATE application SET
			approved = $1,
			reason = $2
		WHERE application_id = $3
	`, status, reason, appID)

	if err != nil {
		return errors.Wrap(err, "failed to update application approval")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err,
			"failed to check result of application approval update")
	} else if n != 1 {
		return errors.Wrapf(
			app.ErrNotFound,
			"no such application by id of %d", appID,
		)
	}

	return nil
}
