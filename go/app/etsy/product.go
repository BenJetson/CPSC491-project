package etsy

import (
	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

type etsyProduct struct {
	ID          int    `json:"listing_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`

	// ... omitted fields
	// Most of the data we get from Etsy we don't actually need.
}

func (ep *etsyProduct) toCommerceProduct() (p app.CommerceProduct, err error) {
	p.ID = ep.ID
	p.Title = ep.Title
	p.Description = ep.Description
	p.ImageURL = "" // TODO figure this out

	if p.Price, err = app.ParseMoneyFromString(ep.Price); err != nil {
		err = errors.Wrap(err, "could not parse Etsy product price")
		return
	}
	return
}
