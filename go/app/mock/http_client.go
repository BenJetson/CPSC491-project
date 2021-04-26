package mock

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// HTTPClient mocks an http.Client by answering any call to the Do method with
// the given fields.
type HTTPClient struct {
	Response string
	Code     int
	Err      error
}

// Do accepts an http.Request (ignores its value) and crafts an appropriate
// response based on the fields set.
//
// If the Err field is not nil, the response will be nil and vice versa.
func (c *HTTPClient) Do(_ *http.Request) (res *http.Response, err error) {
	if c.Err != nil {
		err = c.Err
		return
	}

	res = &http.Response{
		Body:       ioutil.NopCloser(strings.NewReader(c.Response)),
		StatusCode: c.Code,
	}
	return
}

// Reset clears all fields of the mock HTTPClient.
func (c *HTTPClient) Reset() { *c = HTTPClient{} }
