package etsy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

// This is an assertion, which will cause the build to fail if the etsy.API type
// does not implement the app.CommerceVendor interface.
var _ app.CommerceVendor = (*Client)(nil)

// A Client can be used to make requests to the Etsy API, in a way that is
// compliant with the app.CommerceVendor interface.
type Client struct {
	apiKey     string
	httpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
}

// NewClient creates a new etsy.Client given the API key to use.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// NewClientFromEnv attempts to initialize an etsy.Client using an API key
// source from the environment.
func NewClientFromEnv() (*Client, error) {
	apiKey := os.Getenv("ETSY_API_KEY")
	if len(apiKey) < 1 {
		return nil, errors.New("must set ETSY_API_KEY")
	}

	return NewClient(apiKey), nil
}

const apiHost = "openapi.etsy.com"

func (c *Client) attemptReadEtsyErrorReason(b io.Reader) string {
	bodyBytes, err := ioutil.ReadAll(b)
	if err != nil {
		return "unknown"
	}
	return string(bodyBytes)
}

func (c *Client) injectAuth(params *url.Values) {
	params.Set("api_key", c.apiKey)
}

func (c *Client) injectEtsyQueryParams(
	q app.CommerceQuery,
	params *url.Values,
) {

	params.Set("keywords", q.Keywords)

	if q.Sort.Valid {
		params.Set("sort_on", string(q.Sort.By))
		params.Set("sort_order", string(q.Sort.Direction))
	}

	if q.Limit.Valid {
		params.Set("limit", fmt.Sprint(q.Limit.Int64))
	}

	if q.PageNo.Valid {
		params.Set("page", fmt.Sprint(q.PageNo.Int64))
	}
}

func (c *Client) injectCommonParams(params *url.Values) {
	// This requests the inclusion of the main product image.
	params.Set("includes", "MainImage")

	// This would request all images.
	// params.Set("Includes", "Images")
}

type searchResponse struct {
	Count   int           `json:"count"`
	Results []etsyProduct `json:"results"`
	Type    string        `json:"type"`

	// ... omitted fields ...

	// Params  struct {
	// 	Limit     null.String `json:"limit"`
	// 	Offset    int         `json:"offset"`
	// 	Page      null.Int    `json:"page"`
	// 	Keywords  string      `json:"keywords"`
	// 	SortOn    string      `json:"sort_on"`
	// 	SortOrder string      `json:"sort_order"`
	// 	ListingID null.Int   `json:"listing_id"`
	// } `json:"params"`
	// Pagination struct {
	// 	EffectiveLimit  int `json:"effective_limit"`
	// 	EffectiveOffset int `json:"effective_offset"`
	// 	NextOffset      int `json:"next_offset"`
	// 	EffectivePage   int `json:"effective_page"`
	// 	NextPage        int `json:"next_page"`
	// } `json:"pagination"`
}

// Search will query the Etsy catalog and return matching items.
func (c *Client) Search(
	ctx context.Context,
	q app.CommerceQuery,
) ([]app.CommerceProduct, error) {

	params := make(url.Values)
	c.injectAuth(&params)
	c.injectCommonParams(&params)
	c.injectEtsyQueryParams(q, &params)

	u := url.URL{
		Scheme:   "https",
		Host:     apiHost,
		Path:     "/v2/listings/active",
		RawQuery: params.Encode(),
	}

	r, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to craft listing request")
	}

	r = r.WithContext(ctx)

	res, err := c.httpClient.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "client failed to do listing request")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf(
			"received non-OK status for listing request (%d), reason: %s",
			res.StatusCode, c.attemptReadEtsyErrorReason(res.Body),
		)
	}

	var data searchResponse

	d := json.NewDecoder(res.Body)
	if err = d.Decode(&data); err != nil {
		return nil, errors.Wrap(err, "failed to decode Etsy products")
	}

	products := make([]app.CommerceProduct, len(data.Results))
	for idx, ep := range data.Results {
		products[idx], err = ep.toCommerceProduct()
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert Etsy product")
		}
	}

	return products, nil
}

// GetProductByID will fetch an Etsy product by its ID number.
func (c *Client) GetProductByID(
	ctx context.Context,
	productID int,
) (app.CommerceProduct, error) {

	params := make(url.Values)
	c.injectAuth(&params)
	c.injectCommonParams(&params)

	u := url.URL{
		Scheme:   "https",
		Host:     apiHost,
		Path:     fmt.Sprintf("/v2/listings/%d", productID),
		RawQuery: params.Encode(),
	}

	r, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return app.CommerceProduct{},
			errors.Wrap(err, "failed to craft listing request")
	}

	r = r.WithContext(ctx)

	res, err := c.httpClient.Do(r)
	if err != nil {
		return app.CommerceProduct{},
			errors.Wrap(err, "client failed to do listing request")
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return app.CommerceProduct{}, errors.Wrapf(app.ErrNotFound,
			"no Etsy listing exists with ID of %d", productID)
	} else if res.StatusCode != http.StatusOK {
		return app.CommerceProduct{}, errors.Errorf(
			"received non-OK status for listing by ID request (%d), reason: %s",
			res.StatusCode, c.attemptReadEtsyErrorReason(res.Body),
		)
	}

	var data searchResponse

	d := json.NewDecoder(res.Body)
	if err = d.Decode(&data); err != nil {
		return app.CommerceProduct{},
			errors.Wrap(err, "failed to decode Etsy products")
	}

	if len(data.Results) != 1 {
		return app.CommerceProduct{}, errors.Errorf(
			"expected exactly one product but found %d", len(data.Results))
	}

	p, err := data.Results[0].toCommerceProduct()
	if err != nil {
		return app.CommerceProduct{}, errors.Wrap(err,
			"failed to convert Etsy product")
	}

	return p, nil
}
