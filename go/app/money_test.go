package app

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoneyCompare(t *testing.T) {
	a := MustMakeMoneyFromComponents(1, 75)
	b := MustMakeMoneyFromComponents(2, 23)
	c := MustMakeMoneyFromComponents(1, 75)

	testCases := []struct {
		alias          string
		actual, expect bool
	}{
		{
			alias:  "Equal",
			actual: a == c,
			expect: true,
		},
		{
			alias:  "BadEqual",
			actual: a == b,
			expect: false,
		},
		{
			alias:  "NotEqual",
			actual: a != b,
			expect: true,
		},
		{
			alias:  "BadNotEqual",
			actual: a != c,
			expect: false,
		},
		{
			alias:  "Less",
			actual: a < b,
			expect: true,
		},
		{
			alias:  "BadLess",
			actual: b < a,
			expect: false,
		},
		{
			alias:  "More",
			actual: b > a,
			expect: true,
		},
		{
			alias:  "BadMore",
			actual: a > b,
			expect: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.alias, func(t *testing.T) {
			assert.Equal(t, tc.expect, tc.actual)
		})
	}
}

func TestMakeMoneyFromComponents(t *testing.T) {
	testCases := []struct {
		alias          string
		dollars, cents int
		expect         Money
		expectErr      bool
	}{
		{
			alias:   "FiveDollarFootLong",
			dollars: 5,
			cents:   0,
			expect:  500,
		},
		{
			alias:   "Triple",
			dollars: 3,
			cents:   33,
			expect:  333,
		},
		{
			alias:   "JustAFewDollars",
			dollars: 4727189292,
			cents:   20,
			expect:  472718929220,
		},
		{
			alias:     "TooManyCents",
			dollars:   7,
			cents:     100,
			expectErr: true,
		},
		{
			alias:     "TooFewCents",
			dollars:   -4,
			cents:     -100,
			expectErr: true,
		},
		{
			alias:     "TooManyDollars",
			dollars:   92233720368547759, // (maxint / 100)  + 1
			cents:     0,
			expectErr: true,
		},
		{
			alias:     "TooFewDollars",
			dollars:   -92233720368547759, // (minint / 100) - 1
			cents:     0,
			expectErr: true,
		},
		{
			alias:     "DollarSignMismatch",
			dollars:   67,
			cents:     -92,
			expectErr: true,
		},
		{
			alias:     "CentSignMismatch",
			dollars:   -43,
			cents:     46,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.alias, func(t *testing.T) {
			actual, err := MakeMoneyFromComponents(tc.dollars, tc.cents)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expect, actual)
			}
		})
	}
}

func TestMustMakeMoneyFromComponents(t *testing.T) {
	t.Run("ShouldPass", func(t *testing.T) {
		require.NotPanics(t, func() {
			m := MustMakeMoneyFromComponents(30, 22)
			assert.Equal(t, Money(3022), m)
		})
	})

	t.Run("ShouldPanic", func(t *testing.T) {
		require.Panics(t, func() {
			MustMakeMoneyFromComponents(-17, 9)
		})
	})
}

func TestMoneyString(t *testing.T) {
	testCases := []struct {
		alias  string
		m      Money
		expect string
	}{
		{
			alias:  "Zero",
			m:      MustMakeMoneyFromComponents(0, 0),
			expect: "$0.00",
		},
		{
			alias:  "OneCent",
			m:      MustMakeMoneyFromComponents(0, 1),
			expect: "$0.01",
		},
		{
			alias:  "TenCents",
			m:      MustMakeMoneyFromComponents(0, 10),
			expect: "$0.10",
		},
		{
			alias:  "OneDollar",
			m:      MustMakeMoneyFromComponents(1, 0),
			expect: "$1.00",
		},
		{
			alias:  "DollarsAndCents",
			m:      MustMakeMoneyFromComponents(5, 7),
			expect: "$5.07",
		},
		{
			alias:  "BenTen",
			m:      MustMakeMoneyFromComponents(100, 10),
			expect: "$100.10",
		},
		{
			alias:  "TenBen",
			m:      MustMakeMoneyFromComponents(1000, 0),
			expect: "$1,000.00",
		},
		{
			alias: "EnrichmentCenterFailure",
			m:     MustMakeMoneyFromComponents(9999999999999999, 99),
			// i1.theportalwiki.net/img/1/1d/Announcer_openingcourtesy01.wav
			expect: "$9,999,999,999,999,999.99",
		},
		{
			alias:  "NegativeCents",
			m:      MustMakeMoneyFromComponents(0, -7),
			expect: "-$0.07",
		},
		{
			alias:  "NegativeDollars",
			m:      MustMakeMoneyFromComponents(-4273, 0),
			expect: "-$4,273.00",
		},
		{
			alias:  "NegativeDouble",
			m:      MustMakeMoneyFromComponents(-38291, -86),
			expect: "-$38,291.86",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.alias, func(t *testing.T) {
			actual := tc.m.String()
			assert.Equal(t, tc.expect, actual, "string out should match expect")

			actual = fmt.Sprintf("%v", tc.m)
			assert.Equal(t, tc.expect, actual, "format out should match expect")
		})
	}
}

func TestMakeMoneyFromFloat(t *testing.T) {
	testCases := []struct {
		alias  string
		f      float64
		expect Money
	}{
		{
			alias:  "Zero",
			f:      0.0,
			expect: MustMakeMoneyFromComponents(0, 0),
		},
		{
			alias:  "NoFractionalCents",
			f:      37.195,
			expect: MustMakeMoneyFromComponents(37, 19),
		},
		{
			alias:  "JustCents",
			f:      0.01,
			expect: MustMakeMoneyFromComponents(0, 1),
		},
		{
			alias:  "Negative",
			f:      -5.32,
			expect: MustMakeMoneyFromComponents(-5, -32),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.alias, func(t *testing.T) {
			actual := MakeMoneyFromFloat(tc.f)
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func TestMoneyRE(t *testing.T) {
	testCases := []struct {
		str    string
		expect bool
	}{
		{"$0.00", true},
		{"$0.01", true},
		{"$0", true},
		{"$1.21", true},
		{"$1", true},
		{"$100.32", true},
		{"$100", true},
		{"$1,100.32", true},
		{"$22,100.32", true},
		{"$333,100.32", true},
		{"$1,100,100.32", true},
		{"$22,100,100.32", true},
		{"$333,100,100.32", true},
		{"$333,100,100", true},
		{"$1100.32", true},
		{"$22100.32", true},
		{"$333100.32", true},
		{"$1100100.32", true},
		{"$22100100.32", true},
		{"$333100100.32", true},
		{"$333100100", true},
		{"0.00", true},
		{"0.01", true},
		{"0", true},
		{"1.21", true},
		{"1", true},
		{"100.32", true},
		{"100", true},
		{"1,100.32", true},
		{"22,100.32", true},
		{"333,100.32", true},
		{"1,100,100.32", true},
		{"22,100,100.32", true},
		{"333,100,100.32", true},
		{"333,100,100", true},
		{"1100.32", true},
		{"22100.32", true},
		{"333100.32", true},
		{"1100100.32", true},
		{"22100100.32", true},
		{"333100100.32", true},
		{"333100100", true},
		{"-$0.00", true},
		{"-$0.01", true},
		{"-$0", true},
		{"-$1.21", true},
		{"-$1", true},
		{"-$100.32", true},
		{"-$100", true},
		{"-$1,100.32", true},
		{"-$22,100.32", true},
		{"-$333,100.32", true},
		{"-$1,100,100.32", true},
		{"-$22,100,100.32", true},
		{"-$333,100,100.32", true},
		{"-$333,100,100", true},
		{"-$1100.32", true},
		{"-$22100.32", true},
		{"-$333100.32", true},
		{"-$1100100.32", true},
		{"-$22100100.32", true},
		{"-$333100100.32", true},
		{"-$333100100", true},
		{"-0.00", true},
		{"-0.01", true},
		{"-0", true},
		{"-1.21", true},
		{"-1", true},
		{"-100.32", true},
		{"-100", true},
		{"-1,100.32", true},
		{"-22,100.32", true},
		{"-333,100.32", true},
		{"-1,100,100.32", true},
		{"-22,100,100.32", true},
		{"-333,100,100.32", true},
		{"-333,100,100", true},
		{"-1100.32", true},
		{"-22100.32", true},
		{"-333100.32", true},
		{"-1100100.32", true},
		{"-22100100.32", true},
		{"-333100100.32", true},
		{"-333100100", true},
		// examples from the ParseMoneyFromString docs
		{"-$40,324,921.76", true},
		{"$503", true},
		{"504.32", true},
		{"-30", true},
		{"$3202.32", true},
		{"8,108.31", true},
		{"244,422", true},
		// some nonexamples
		{"400.5", false},        // partial cents
		{"40,42.41", false},     // wrong placement of comma
		{"$-493.21", false},     // negative after dollar sign
		{"$3$33.23", false},     // repeat pattern not allowed
		{"$53.595", false},      // fractional cents not allowed
		{"$1000,000.43", false}, // if using commas must put all commas
		{"$$400", false},        // double dollar sign
		{"$.32", false},         // must specify zero dollars
		{"$garbage", false},     // garbage; self explanatory
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", idx), func(t *testing.T) {
			actual := moneyRE.Match([]byte(tc.str))
			assert.Equalf(t, tc.expect, actual,
				"string '%s' ought to return %t", tc.str, tc.expect)
		})
	}
}

func TestParseMoneyFromString(t *testing.T) {
	testCases := []struct {
		str       string
		expect    Money
		expectErr bool
	}{
		{
			str:    "$0.00",
			expect: Money(0),
		},
		{
			str:    "$0.01",
			expect: Money(1),
		},
		{
			str:    "$0",
			expect: Money(0),
		},
		{
			str:    "$1.21",
			expect: Money(121),
		},
		{
			str:    "$1",
			expect: Money(100),
		},
		{
			str:    "$100.32",
			expect: Money(10032),
		},
		{
			str:    "$100",
			expect: Money(10000),
		},
		{
			str:    "$1,100.32",
			expect: Money(110032),
		},
		{
			str:    "$22,100.32",
			expect: Money(2210032),
		},
		{
			str:    "$333,100.32",
			expect: Money(33310032),
		},
		{
			str:    "$1,100,100.32",
			expect: Money(110010032),
		},
		{
			str:    "$22,100,100.32",
			expect: Money(2210010032),
		},
		{
			str:    "$333,100,100.32",
			expect: Money(33310010032),
		},
		{
			str:    "$333,100,100",
			expect: Money(33310010000),
		},
		{
			str:    "$1100.32",
			expect: Money(110032),
		},
		{
			str:    "$22100.32",
			expect: Money(2210032),
		},
		{
			str:    "$333100.32",
			expect: Money(33310032),
		},
		{
			str:    "$1100100.32",
			expect: Money(110010032),
		},
		{
			str:    "$22100100.32",
			expect: Money(2210010032),
		},
		{
			str:    "$333100100.32",
			expect: Money(33310010032),
		},
		{
			str:    "$333100100",
			expect: Money(33310010000),
		},
		{
			str:    "0.00",
			expect: Money(0),
		},
		{
			str:    "0.01",
			expect: Money(1),
		},
		{
			str:    "0",
			expect: Money(0),
		},
		{
			str:    "1.21",
			expect: Money(121),
		},
		{
			str:    "1",
			expect: Money(100),
		},
		{
			str:    "100.32",
			expect: Money(10032),
		},
		{
			str:    "100",
			expect: Money(10000),
		},
		{
			str:    "1,100.32",
			expect: Money(110032),
		},
		{
			str:    "22,100.32",
			expect: Money(2210032),
		},
		{
			str:    "333,100.32",
			expect: Money(33310032),
		},
		{
			str:    "1,100,100.32",
			expect: Money(110010032),
		},
		{
			str:    "22,100,100.32",
			expect: Money(2210010032),
		},
		{
			str:    "333,100,100.32",
			expect: Money(33310010032),
		},
		{
			str:    "333,100,100",
			expect: Money(33310010000),
		},
		{
			str:    "1100.32",
			expect: Money(110032),
		},
		{
			str:    "22100.32",
			expect: Money(2210032),
		},
		{
			str:    "333100.32",
			expect: Money(33310032),
		},
		{
			str:    "1100100.32",
			expect: Money(110010032),
		},
		{
			str:    "22100100.32",
			expect: Money(2210010032),
		},
		{
			str:    "333100100.32",
			expect: Money(33310010032),
		},
		{
			str:    "333100100",
			expect: Money(33310010000),
		},
		{
			str:    "-$0.00",
			expect: Money(0),
		},
		{
			str:    "-$0.01",
			expect: Money(-1),
		},
		{
			str:    "-$0",
			expect: Money(0),
		},
		{
			str:    "-$1.21",
			expect: Money(-121),
		},
		{
			str:    "-$1",
			expect: Money(-100),
		},
		{
			str:    "-$100.32",
			expect: Money(-10032),
		},
		{
			str:    "-$100",
			expect: Money(-10000),
		},
		{
			str:    "-$1,100.32",
			expect: Money(-110032),
		},
		{
			str:    "-$22,100.32",
			expect: Money(-2210032),
		},
		{
			str:    "-$333,100.32",
			expect: Money(-33310032),
		},
		{
			str:    "-$1,100,100.32",
			expect: Money(-110010032),
		},
		{
			str:    "-$22,100,100.32",
			expect: Money(-2210010032),
		},
		{
			str:    "-$333,100,100.32",
			expect: Money(-33310010032),
		},
		{
			str:    "-$333,100,100",
			expect: Money(-33310010000),
		},
		{
			str:    "-$1100.32",
			expect: Money(-110032),
		},
		{
			str:    "-$22100.32",
			expect: Money(-2210032),
		},
		{
			str:    "-$333100.32",
			expect: Money(-33310032),
		},
		{
			str:    "-$1100100.32",
			expect: Money(-110010032),
		},
		{
			str:    "-$22100100.32",
			expect: Money(-2210010032),
		},
		{
			str:    "-$333100100.32",
			expect: Money(-33310010032),
		},
		{
			str:    "-$333100100",
			expect: Money(-33310010000),
		},
		{
			str:    "-0.00",
			expect: Money(0),
		},
		{
			str:    "-0.01",
			expect: Money(-1),
		},
		{
			str:    "-0",
			expect: Money(0),
		},
		{
			str:    "-1.21",
			expect: Money(-121),
		},
		{
			str:    "-1",
			expect: Money(-100),
		},
		{
			str:    "-100.32",
			expect: Money(-10032),
		},
		{
			str:    "-100",
			expect: Money(-10000),
		},
		{
			str:    "-1,100.32",
			expect: Money(-110032),
		},
		{
			str:    "-22,100.32",
			expect: Money(-2210032),
		},
		{
			str:    "-333,100.32",
			expect: Money(-33310032),
		},
		{
			str:    "-1,100,100.32",
			expect: Money(-110010032),
		},
		{
			str:    "-22,100,100.32",
			expect: Money(-2210010032),
		},
		{
			str:    "-333,100,100.32",
			expect: Money(-33310010032),
		},
		{
			str:    "-333,100,100",
			expect: Money(-33310010000),
		},
		{
			str:    "-1100.32",
			expect: Money(-110032),
		},
		{
			str:    "-22100.32",
			expect: Money(-2210032),
		},
		{
			str:    "-333100.32",
			expect: Money(-33310032),
		},
		{
			str:    "-1100100.32",
			expect: Money(-110010032),
		},
		{
			str:    "-22100100.32",
			expect: Money(-2210010032),
		},
		{
			str:    "-333100100.32",
			expect: Money(-33310010032),
		},
		{
			str:    "-333100100",
			expect: Money(-33310010000),
		},
		{
			str:    "-$40,324,921.76",
			expect: Money(-4032492176),
		},
		{
			str:    "$503",
			expect: Money(50300),
		},
		{
			str:    "504.32",
			expect: Money(50432),
		},
		{
			str:    "-30",
			expect: Money(-3000),
		},
		{
			str:    "$3202.32",
			expect: Money(320232),
		},
		{
			str:    "8,108.31",
			expect: Money(810831),
		},
		{
			str:    "244,422",
			expect: Money(24442200),
		},
		{
			str:       "400.5",
			expectErr: true,
		},
		{
			str:       "40,42.41",
			expectErr: true,
		},
		{
			str:       "$-493.21",
			expectErr: true,
		},
		{
			str:       "$3$33.23",
			expectErr: true,
		},
		{
			str:       "$53.595",
			expectErr: true,
		},
		{
			str:       "$1000,000.43",
			expectErr: true,
		},
		{
			str:       "$$400",
			expectErr: true,
		},
		{
			str:       "$.32",
			expectErr: true,
		},
		{
			str:       "$garbage",
			expectErr: true,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", idx), func(t *testing.T) {
			actual, err := ParseMoneyFromString(tc.str)

			if tc.expectErr {
				require.Errorf(t, err,
					"string '%s' ought to produce an error", tc.str)
			} else {
				require.NoError(t, err,
					"string '%s' ought not produce an error", tc.str)
				assert.Equalf(t, tc.expect, actual,
					"string '%s' ought to return %t", tc.str, tc.expect)
			}
		})
	}
}

func TestConvertToPoints(t *testing.T) {
	orgs := map[string]Organization{
		"CheapCo": {
			ID:   1,
			Name: "Cheap Co.",
			// A point is worth one cent, how cheap! xD
			PointValue: MustMakeMoneyFromComponents(0, 1),
		},
		"BigBucks": {
			ID:   2,
			Name: "Big Bucks, Inc.",
			// A point is worth $2.25, seemingly generous.
			PointValue: MustMakeMoneyFromComponents(2, 25),
		},
		"ACME": {
			ID:   3,
			Name: "ACME LLC",
			// A point is worth $1.33, a strange amount for a mysterious company.
			PointValue: MustMakeMoneyFromComponents(1, 33),
		},
	}

	costs := map[string]Money{
		"A": MustMakeMoneyFromComponents(5, 0),
		"B": MustMakeMoneyFromComponents(200, 67),
		"C": MustMakeMoneyFromComponents(0, 75),
		"D": MustMakeMoneyFromComponents(34, 88),
	}

	testCases := []struct {
		cost         string
		org          string
		expectAmount int
	}{
		{
			cost:         "A",
			org:          "CheapCo",
			expectAmount: 500,
		},
		{
			cost:         "B",
			org:          "CheapCo",
			expectAmount: 20067,
		},
		{
			cost:         "C",
			org:          "CheapCo",
			expectAmount: 75,
		},
		{
			cost:         "D",
			org:          "CheapCo",
			expectAmount: 3488,
		},
		{
			cost:         "A",
			org:          "BigBucks",
			expectAmount: 3,
		},
		{
			cost:         "B",
			org:          "BigBucks",
			expectAmount: 90,
		},
		{
			cost:         "C",
			org:          "BigBucks",
			expectAmount: 1,
		},
		{
			cost:         "D",
			org:          "BigBucks",
			expectAmount: 16,
		},
		{
			cost:         "A",
			org:          "ACME",
			expectAmount: 4,
		},
		{
			cost:         "B",
			org:          "ACME",
			expectAmount: 151,
		},
		{
			cost:         "C",
			org:          "ACME",
			expectAmount: 1,
		},
		{
			cost:         "D",
			org:          "ACME",
			expectAmount: 27,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%s", tc.org, tc.cost), func(t *testing.T) {
			cost := costs[tc.cost]
			org := orgs[tc.org]

			actual := cost.ConvertToPoints(org)

			assert.Equal(t, tc.expectAmount, actual.Amount)
			assert.Equal(t, org.ID, actual.OrganizationID)
			assert.Equal(t, org.PointValue, actual.PointValue)
		})
	}
}
