package db

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/BenJetson/CPSC491-project/go/app"
)

func TestCreatePerson(t *testing.T) {
	db := newTestDB(t)
	defer db.cleanup(t)

	p, err := app.NewPassword("aoa")
	require.NoError(t, err)

	err = db.CreatePerson(app.Person{
		FirstName: "Ben",
		LastName:  "Godfrey",
		Email:     "bfgodfr@clemson.edu",
		Password:  p,
		Role:      app.RoleAdmin,
	})
	require.NoError(t, err)

	db.assertCount(t, "person", 1)

	// FIXME this is not a good test, write a new one!
	// TODO just an example, rewrite this!
}
