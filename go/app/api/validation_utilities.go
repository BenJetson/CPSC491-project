package api

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

// validateEmail is a regular expression validator for email addresses.
//
// Source: W3C speficiation for email adresses.
// URL: https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html
var validateEmail = regexp.MustCompile(
	`^[a-zA-Z0-9.!#$%&'*+\/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]` +
		`{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}` +
		`[a-zA-Z0-9])?)*$`,
)

func validatePassword(pass string) (message string, err error) {
	const minPassLength = 8

	if len(pass) < minPassLength {
		err = errors.New("password does not meet length requirement")
		message = fmt.Sprintf("Password must be at least %d characters long.",
			minPassLength)
		return
	}

	return
}
