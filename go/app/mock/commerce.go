package mock

import (
	"context"

	"github.com/BenJetson/CPSC491-project/go/app"
)

// This is an assertion, which will cause the build to fail if the etsy.API type
// does not implement the app.CommerceVendor interface.
var _ app.CommerceVendor = (*CommerceVendor)(nil)

// A CommerceVendor mocks a commerce vendor API.
type CommerceVendor struct{}

// Search will query the vendor's catalog and return matching items.
func (c *CommerceVendor) Search(
	ctx context.Context,
	q app.CommerceQuery,
) ([]app.CommerceProduct, error) {

	return nil, nil
}

// GetProductByID will fetch a product by its ID number from the vendor's
// catalog.
func (c *CommerceVendor) GetProductByID(
	ctx context.Context,
	productID int,
) (app.CommerceProduct, error) {

	return app.CommerceProduct{}, nil
}
