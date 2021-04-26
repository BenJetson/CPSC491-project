package db

import (
	"context"
	"database/sql"
	"strings"

	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

func (db *database) GetProductsForOrganization(
	ctx context.Context,
	orgID int,
) ([]app.CatalogProduct, error) {

	org, err := db.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve org for catalog")
	}

	var ps []app.Product

	err = db.SelectContext(ctx, &ps, `
		SELECT
			product_id,
			vendor_id,
			organization_id,
			title,
			description,
			image_url,
			price,
			is_available
		FROM product
		WHERE
			organization_id = $1
			AND is_available = TRUE
		ORDER BY title ASC
	`, orgID)

	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve products")
	}

	cps := make([]app.CatalogProduct, len(ps))

	for idx, p := range ps {
		cps[idx] = p.ToCatalogProduct(org)
	}

	return cps, nil
}

func (db *database) SearchProductCatalog(
	ctx context.Context,
	orgID int,
	keywords string,
) ([]app.CatalogProduct, error) {

	org, err := db.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve org for catalog")
	}

	var ps []app.Product

	if len(keywords) > 100 {
		keywords = keywords[:100]
	}

	// TODO this should be secure because we are passing it to the database
	// still as a parameter. However, should still research this.
	fullTextSearchQuery := strings.Join(strings.Split(keywords, " "), " & ")

	err = db.SelectContext(ctx, &ps, `
		SELECT
			product_id,
			vendor_id,
			organization_id,
			title,
			description,
			image_url,
			price,
			is_available
		FROM product
		WHERE
			organization_id = $1
			AND is_available = TRUE
			AND searchable @@ to_tsquery($2)
		ORDER BY
			ts_rank(searchable, to_tsquery($2)) DESC,
			title ASC
	`, orgID, fullTextSearchQuery)

	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve products")
	}

	cps := make([]app.CatalogProduct, len(ps))

	for idx, p := range ps {
		cps[idx] = p.ToCatalogProduct(org)
	}

	return cps, nil
}

func (db *database) GetProductByID(
	ctx context.Context,
	productID, orgID int,
) (app.CatalogProduct, error) {

	org, err := db.GetOrganizationByID(ctx, orgID)
	if err != nil {
		return app.CatalogProduct{}, errors.Wrap(err,
			"failed to retrieve org for catalog")
	}

	var p app.Product

	err = db.GetContext(ctx, &p, `
		SELECT
			product_id,
			vendor_id,
			organization_id,
			title,
			description,
			image_url,
			price,
			is_available
		FROM product
		WHERE
			product_id = $1
			AND organization_id = $2
			AND is_available = TRUE
		ORDER BY title ASC
	`, productID, orgID)

	if errors.Is(err, sql.ErrNoRows) {
		return app.CatalogProduct{}, errors.Wrapf(
			app.ErrNotFound,
			"no such product by id of '%d' with organization id of '%d'",
			productID, orgID,
		)
	} else if err != nil {
		return app.CatalogProduct{}, errors.Wrap(err,
			"failed to retrieve product")
	}

	return p.ToCatalogProduct(org), nil
}

func (db *database) AddProduct(
	ctx context.Context,
	p app.Product,
) (int, error) {

	var id int
	err := db.GetContext(ctx, &id, `
		INSERT INTO product (
			vendor_id,
			organization_id,
			title,
			description,
			image_url,
			price
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (vendor_id, organization_id)
		DO UPDATE SET
			title = $3,
			description = $4,
			image_url = $5,
			price = $6,
			is_available = TRUE
		RETURNING product_id `,

		p.VendorID,       // $1
		p.OrganizationID, // $2
		p.Title,          // $3
		p.Description,    // $4
		p.ImageURL,       // $5
		p.Price,          // $6
	)

	return id, errors.Wrap(err, "failed to insert product")
}

func (db *database) MakeProductUnavailable(
	ctx context.Context,
	productID, orgID int,
) error {

	result, err := db.ExecContext(ctx, `
		UPDATE product SET
			is_available = FALSE
		WHERE
			product_id = $1
			AND organization_id = $2
	`, productID, orgID)

	if err != nil {
		return errors.Wrap(err, "failed to flag product as unavailable")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err,
			"failed to check result of product availability update")
	} else if n != 1 {
		return errors.Wrapf(
			app.ErrNotFound,
			"no such product by id of %d in organization %d", productID, orgID,
		)
	}

	return nil
}
