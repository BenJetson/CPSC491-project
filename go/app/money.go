package app

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/pkg/errors"
)

// Money provides an abstraction layer around monetary values stored as an
// integer count of cents.
type Money int

const (
	// CentsPerDollar describes the number of cents in a US Dollar.
	CentsPerDollar int = 100

	// MaxDollarAmount is the maximum dollar amount a Money can store.
	MaxDollarAmount int = math.MaxInt64 / CentsPerDollar
	// MinDollarAmount is the minimum dollar amount a Money can store.
	MinDollarAmount int = math.MinInt64 / CentsPerDollar

	// MaxCentAmount is the maximum number of *fractional* cents in a US Dollar.
	MaxCentAmount int = 99
	// MinCentAmount is the minimum number of *fractional* cents in a US Dollar.
	MinCentAmount int = -99
)

// MakeMoneyFromFloat makes a Money from a floating point integer. The decimal
// will be truncated at two places (floored).
func MakeMoneyFromFloat(f float64) Money {
	// Multiply by cents per dollar constant to shift the decimal two places
	// to the left. Then, cast to Money (an alias of int) which will truncate
	// remaining decimal values, an effective floor operation.

	return Money(f * float64(CentsPerDollar))
}

// MakeMoneyFromComponents makes a Money from its dollar and cent component
// values. Both component values must be within their respective ranges and
// have matching signs (unless either is zero).
func MakeMoneyFromComponents(dollars, cents int) (Money, error) {
	var m Money

	if dollars > MaxDollarAmount || dollars < MinDollarAmount {
		return m, fmt.Errorf(
			"dollars component outside of valid range [%d, %d]",
			MinDollarAmount, MaxDollarAmount,
		)
	} else if cents > MaxCentAmount || cents < MinCentAmount {
		return m, errors.Errorf(
			"cents component outside of valid range [%d, %d]",
			MinCentAmount, MaxCentAmount,
		)
	}

	dNeg, cNeg := dollars < 0, cents < 0
	dZero, cZero := dollars == 0, cents == 0

	if dNeg != cNeg && !dZero && !cZero {
		return m, errors.New("component sign mismatch - dollars and cents " +
			"must have same sign unless one is zero")
	}

	total := cents
	total += dollars * CentsPerDollar

	return Money(total), nil
}

// MustMakeMoneyFromComponents runs MakeMoneyFromComponents and follows the
// same constraints. HOWEVER, when any precondition fails and an error occurs,
// this function shall panic.
//
// It is intended for this function only to be used by tests and constants
// where error checking often is impractical. All other use cases should use
// MakeMoneyFromComponents and check the error.
func MustMakeMoneyFromComponents(dollars, cents int) Money {
	m, err := MakeMoneyFromComponents(dollars, cents)
	if err != nil {
		panic(err)
	}
	return m
}

// moneyRE is used to check money strings for validity, per the rules of
// ParseMoneyFromString.
var moneyRE = regexp.MustCompile(
	`^-?\$?[0-9]{1,3}((,[0-9]{3})*|([0-9]{3})*)(\.[0-9]{2})?$`)

// ParseMoneyFromString attempts to parse Money from a string.
//
// This string must follow these rules:
// 	- if there is a negative sign, it must be the first character
// 	- the dollar sign is optional
// 	- the commas separating thousands in the dollar amount are optional
// 	- the cents are optional, but if present must have two decimal places
//
// Examples of accepted strings:
// 	- "-$40,324,921.76"
// 	- "$503"
// 	- "504.32"
// 	- "-30"
// 	- "$3202.32"
// 	- "8,108.31"
// 	- "244,422"
func ParseMoneyFromString(s string) (Money, error) {
	var m Money

	// Check to make sure the given string matches the accepted format.
	if !moneyRE.Match([]byte(s)) {
		return m, errors.New("money string does not match expected format")
	}

	// Discard useless constructs that markup the string for humans.
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")

	// All that should remain is a floating point number; attempt to parse.
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// This should be impossible; regex match should guarantee a valid
		// float after the deletions, but just to be safe.
		return m, errors.Wrap(err, "failed to parse money string as float")
	}

	// Return Money value of the parsed floating point number.
	return MakeMoneyFromFloat(f), nil
}

// Components breaks a Money down into its dollar and cent components.
func (m Money) Components() (dollars, cents int) {
	dollars = int(m) / CentsPerDollar
	cents = int(m) % CentsPerDollar
	return
}

// String returns a string representation of this Money value. Will have a
// leading dollar sign and use thousands separators for human readability.
//
// This returned value is compatible with ParseMoneyFromString.
func (m Money) String() string {
	dollars, cents := m.Components()

	sign := ""

	if dollars < 0 {
		sign = "-"
		dollars *= -1
	}

	if cents < 0 {
		sign = "-"
		cents *= -1
	}

	return fmt.Sprintf(
		"%s$%s.%02d",
		sign,
		humanize.Comma(int64(dollars)),
		cents,
	)
}

// ConvertToPoints converts this Money value into points of the given
// organization based on the point value.
func (m Money) ConvertToPoints(org Organization) Points {
	// Pad the money value so that it rounds up to the next whole point.
	// Items must cost at least one point, and this guarantees that also.
	padding := m % org.PointValue
	if padding != 0 {
		padding = org.PointValue - padding
	}

	// The amount is then the money value divided by the point value.
	amt := (m + padding) / org.PointValue

	return Points{
		Amount:         int(amt),
		PointValue:     org.PointValue,
		OrganizationID: org.ID,
	}
}
