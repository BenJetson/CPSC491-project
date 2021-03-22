package db

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BenJetson/CPSC491-project/go/app"
)

func TestCreatePerson(t *testing.T) {
	db := newTestDB(t)
	defer db.cleanup(t)

	pass1, err := app.NewPassword("aoa")
	require.NoError(t, err)

	pass2, err := app.NewPassword("asdf")
	require.NoError(t, err)

	ctx := context.Background()

	p1 := app.Person{
		FirstName: "Ben",
		LastName:  "Godfrey",
		Email:     "bfgodfr@clemson.edu",
		Password:  pass1,
		Role:      app.RoleAdmin,
	}

	p2 := app.Person{
		FirstName: "Roger",
		LastName:  "Van Scoy",
		Email:     "vanscoy@clemson.edu",
		Password:  pass2,
		Role:      app.RoleSponsor,
	}

	var id int

	t.Run("AddOne", func(t *testing.T) {
		id, err = db.CreatePerson(ctx, p1)
		require.NoError(t, err)
		assert.Equal(t, 1, id)

		db.assertCount(t, "person", 1)
		db.assertCountOf(t, "person", 1, `
			first_name = $1
			AND last_name = $2
			AND email = $3
			AND pass_hash = $4
			AND role_id = $5
		`, p1.FirstName, p1.LastName, p1.Email, p1.Password, p1.Role)
	})

	t.Run("AddAnother", func(t *testing.T) {
		id, err = db.CreatePerson(ctx, p2)
		require.NoError(t, err)
		assert.Equal(t, 2, id)

		db.assertCount(t, "person", 2)
		db.assertCountOf(t, "person", 1, `
			person_id = 1
			AND first_name = $1
			AND last_name = $2
			AND email = $3
			AND pass_hash = $4
			AND role_id = $5
		`, p1.FirstName, p1.LastName, p1.Email, p1.Password, p1.Role)
		db.assertCountOf(t, "person", 1, `
			person_id = 2
			AND first_name = $1
			AND last_name = $2
			AND email = $3
			AND pass_hash = $4
			AND role_id = $5
		`, p2.FirstName, p2.LastName, p2.Email, p2.Password, p2.Role)
	})
}

func TestGetPersonByID(t *testing.T) {
	db := newTestDB(t)
	defer db.cleanup(t)

	pass1, err := app.NewPassword("aoa")
	require.NoError(t, err)

	pass2, err := app.NewPassword("asdf")
	require.NoError(t, err)

	ctx := context.Background()

	p1 := app.Person{
		ID:           1,
		FirstName:    "Ben",
		LastName:     "Godfrey",
		Email:        "bfgodfr@clemson.edu",
		Password:     pass1,
		Role:         app.RoleAdmin,
		Affiliations: make([]int, 0),
	}
	_, err = db.CreatePerson(ctx, p1)
	require.NoError(t, err)

	p2 := app.Person{
		ID:           2,
		FirstName:    "Roger",
		LastName:     "Van Scoy",
		Email:        "vanscoy@clemson.edu",
		Password:     pass2,
		Role:         app.RoleSponsor,
		Affiliations: make([]int, 0),
	}
	_, err = db.CreatePerson(ctx, p2)
	require.NoError(t, err)

	var actual app.Person

	t.Run("FetchOne", func(t *testing.T) {
		actual, err = db.GetPersonByID(ctx, 1)
		require.NoError(t, err)
		assert.Equal(t, p1, actual)
	})

	t.Run("FetchTwo", func(t *testing.T) {
		actual, err = db.GetPersonByID(ctx, 2)
		require.NoError(t, err)
		assert.Equal(t, p2, actual)
	})

	t.Run("NoSuchPerson", func(t *testing.T) {
		_, err = db.GetPersonByID(ctx, 7)
		require.Error(t, err)
		assert.True(t, errors.Is(err, app.ErrNotFound))
	})
}

func TestGetPersonByEmail(t *testing.T) {
	db := newTestDB(t)
	defer db.cleanup(t)

	pass1, err := app.NewPassword("aoa")
	require.NoError(t, err)

	pass2, err := app.NewPassword("asdf")
	require.NoError(t, err)

	ctx := context.Background()

	p1 := app.Person{
		ID:           1,
		FirstName:    "Ben",
		LastName:     "Godfrey",
		Email:        "bfgodfr@clemson.edu",
		Password:     pass1,
		Role:         app.RoleAdmin,
		Affiliations: make([]int, 0),
	}
	_, err = db.CreatePerson(ctx, p1)
	require.NoError(t, err)

	p2 := app.Person{
		ID:           2,
		FirstName:    "Roger",
		LastName:     "Van Scoy",
		Email:        "vanscoy@clemson.edu",
		Password:     pass2,
		Role:         app.RoleSponsor,
		Affiliations: make([]int, 0),
	}
	_, err = db.CreatePerson(ctx, p2)
	require.NoError(t, err)

	var actual app.Person

	t.Run("FetchOne", func(t *testing.T) {
		actual, err = db.GetPersonByEmail(ctx, "bfgodfr@clemson.edu")
		require.NoError(t, err)
		assert.Equal(t, p1, actual)
	})

	t.Run("FetchTwo", func(t *testing.T) {
		actual, err = db.GetPersonByEmail(ctx, "vanscoy@clemson.edu")
		require.NoError(t, err)
		assert.Equal(t, p2, actual)
	})

	t.Run("NoSuchPerson", func(t *testing.T) {
		_, err = db.GetPersonByEmail(ctx, "aack@void.net")
		require.Error(t, err)
		assert.True(t, errors.Is(err, app.ErrNotFound))
	})
}
