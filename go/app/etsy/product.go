package etsy

import (
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v4"

	"github.com/BenJetson/CPSC491-project/go/app"
)

type etsyProductImage struct {
	ListingImageID int64  `json:"listing_image_id"`
	URL75X75       string `json:"url_75x75"`
	URL170X135     string `json:"url_170x135"`
	URL570Xn       string `json:"url_570xN"`
	URLFullxfull   string `json:"url_fullxfull"`
}

type etsyProduct struct {
	ID          int    `json:"listing_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       string `json:"price"`

	MainImage etsyProductImage `json:"MainImage"`

	// ... omitted fields
	// Most of the data we get from Etsy we don't actually need.
}

func (ep *etsyProduct) toCommerceProduct() (p app.CommerceProduct, err error) {
	p.ID = ep.ID
	p.Title = ep.Title
	p.Description = ep.Description

	if len(ep.MainImage.URL170X135) > 0 {
		p.ImageURL = null.StringFrom(ep.MainImage.URL170X135)
	}

	if p.Price, err = app.ParseMoneyFromString(ep.Price); err != nil {
		err = errors.Wrap(err, "could not parse Etsy product price")
		return
	}
	return
}
