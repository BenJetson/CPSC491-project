package app

import (
	"context"

	"gopkg.in/guregu/null.v4"
)

type CommerceSortBy string

const (
	CommerceSortByCreated CommerceSortBy = "created"
	CommerceSortByPrice   CommerceSortBy = "price"
	CommerceSortByRating  CommerceSortBy = "score"
)

type CommerceSortDirection string

const (
	CommerceSortDirectionAscending  CommerceSortDirection = "up"
	CommerceSortDirectionDescending CommerceSortDirection = "up"
)

type CommerceSort struct {
	By        CommerceSortBy
	Direction CommerceSortDirection
	Valid     bool
}

type CommerceQuery struct {
	Keywords string
	Sort     CommerceSort
	Limit    null.Int
	PageNo   null.Int
}

type CommerceProduct struct {
	ID          int
	Title       string
	Description string
	ImageURL    string
	Price       Money
}

type CommerceVendor interface {
	Search(ctx context.Context, q CommerceQuery) ([]CommerceProduct, error)
	GetProductByID(ctx context.Context, productID int) (CommerceProduct, error)
}
