package app

import "gopkg.in/guregu/null.v4"

// A Product represents an entry in the product table.
type Product struct {
	ID             int         `db:"product_id"`
	VendorID       int         `db:"vendor_id"`
	OrganizationID int         `db:"organization_id"`
	Title          string      `db:"title"`
	Description    string      `db:"description"`
	ImageURL       null.String `db:"image_url"`
	Price          Money       `db:"price"`
	IsAvailable    bool        `db:"is_available"`
}

// ToCatalogProduct converts this Product to a CatalogProduct.
func (p *Product) ToCatalogProduct(org Organization) CatalogProduct {
	return CatalogProduct{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		ImageURL:    p.ImageURL,
		Points:      p.Price.ConvertToPoints(org),
	}
}

// A CatalogProduct is a product affiliated with an organization, with its cost
// measured in Points.
type CatalogProduct struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	ImageURL    null.String `json:"image_url"`
	Points      Points      `json:"points"`
}
