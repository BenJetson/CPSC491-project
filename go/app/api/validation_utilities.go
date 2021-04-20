package api

import "regexp"

// validateEmail is a regular expression validator for email addresses.
//
// Source: W3C speficiation for email adresses.
// URL: https://www.w3.org/TR/2016/REC-html51-20161101/sec-forms.html
var validateEmail = regexp.MustCompile(
	`^[a-zA-Z0-9.!#$%&'*+\/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]` +
		`{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}` +
		`[a-zA-Z0-9])?)*$`,
)
