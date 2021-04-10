package app

import (
	"context"

	"gopkg.in/guregu/null.v4"
)

// CommerceSortBy is a pseudo-enumeration of fields on which products may be
// sorted in search.
type CommerceSortBy string

// CommerceSortBy options control the field by which search results are sorted.
const (
	CommerceSortByCreated CommerceSortBy = "created"
	CommerceSortByPrice   CommerceSortBy = "price"
	CommerceSortByRating  CommerceSortBy = "score"
)

// CommerceSortDirection is a pseudo-enumeration of directions for sorting
// search results.
type CommerceSortDirection string

// CommerceSearchDirection options control the direction search results
// are sorted.
const (
	CommerceSortDirectionAscending  CommerceSortDirection = "up"
	CommerceSortDirectionDescending CommerceSortDirection = "up"
)

// CommerceSort controls sorting of search results.
type CommerceSort struct {
	// By controls which field results are sorted by.
	By CommerceSortBy
	// Direction controls which direction results are sorted.
	Direction CommerceSortDirection
	// Valid controls whether or not results are sorted.
	Valid bool
}

// A CommerceQuery desceribes a search query for products.
type CommerceQuery struct {
	Keywords string
	Sort     CommerceSort
	Limit    null.Int
	PageNo   null.Int
}

// A CommerceProduct describes a product of a third-party vendor.
type CommerceProduct struct {
	ID          int
	Title       string
	Description string
	ImageURL    null.String
	Price       Money
}

// CommerceVendor describes a common interface for dealing with third-party
// eCommerce vendors.
type CommerceVendor interface {
	Search(ctx context.Context, q CommerceQuery) ([]CommerceProduct, error)
	GetProductByID(ctx context.Context, productID int) (CommerceProduct, error)
}
