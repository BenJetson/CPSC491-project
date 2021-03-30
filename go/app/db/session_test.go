package db

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/BenJetson/CPSC491-project/go/app"
)

func TestGetSessionsForPerson(t *testing.T) {
	db := newTestDB(t)
	defer db.cleanup(t)

	ctx := context.Background()

	p1 := app.Person{
		ID:           1,
		FirstName:    "Ben",
		LastName:     "Godfrey",
		Email:        "bfgodfr@clemson.edu",
		Password:     `qwerty`,
		Role:         app.RoleAdmin,
		Affiliations: make([]int, 0),
	}
	_, err := db.CreatePerson(ctx, p1)
	require.NoError(t, err)

	p2 := app.Person{
		ID:           2,
		FirstName:    "Roger",
		LastName:     "Van Scoy",
		Email:        "vanscoy@clemson.edu",
		Password:     `zxcvbn`,
		Role:         app.RoleSponsor,
		Affiliations: make([]int, 0),
	}
	_, err = db.CreatePerson(ctx, p2)
	require.NoError(t, err)

	var ss []app.Session

	t.Run("NoSessions1IncludeInvalid", func(t *testing.T) {
		ss, err = db.GetSessionsForPerson(ctx, 1, true)
		require.NoError(t, err)
		require.Len(t, ss, 0)
	})

	t.Run("NoSessions2IncludeInvalid", func(t *testing.T) {
		ss, err = db.GetSessionsForPerson(ctx, 2, true)
		require.NoError(t, err)
		require.Len(t, ss, 0)
	})

	t.Run("NoSessions1ExcludeInvalid", func(t *testing.T) {
		ss, err = db.GetSessionsForPerson(ctx, 1, false)
		require.NoError(t, err)
		require.Len(t, ss, 0)
	})

	t.Run("NoSessions2ExcludeInvalid", func(t *testing.T) {
		ss, err = db.GetSessionsForPerson(ctx, 2, false)
		require.NoError(t, err)
		require.Len(t, ss, 0)
	})

	var s *app.Session
	var ss1, ss2 []app.Session
	for i := 0; i < 5; i++ {
		s, err = app.NewSession(p1)
		require.NoError(t, err)
		require.NotNil(t, s)

		_, err = db.CreateSession(ctx, *s)
		require.NoError(t, err)

		s.ID = i + 1

		ss1 = append(ss1, *s)
	}

	for i := 0; i < 3; i++ {
		s, err = app.NewSession(p2)
		require.NoError(t, err)
		require.NotNil(t, s)

		_, err = db.CreateSession(ctx, *s)
		require.NoError(t, err)

		s.ID = i + 6

		ss2 = append(ss2, *s)
	}

	sortSessions := func(ss []app.Session) {
		sort.Slice(ss, func(i, j int) bool { return ss[i].ID < ss[j].ID })
	}

	t.Run("GetFor1IncludeInvalid", func(t *testing.T) {
		ss, err = db.GetSessionsForPerson(ctx, 1, true)
		require.NoError(t, err)

		sortSessions(ss)
		assertEqualJSON(t, ss1, ss)
	})

	t.Run("GetFor2IncludeInvalid", func(t *testing.T) {
		ss, err = db.GetSessionsForPerson(ctx, 2, true)
		require.NoError(t, err)

		sortSessions(ss)
		assertEqualJSON(t, ss2, ss)
	})

	// Invalidate a few sessions
	err = db.RevokeSession(ctx, 3)
	require.NoError(t, err)
	err = db.RevokeSession(ctx, 4)
	require.NoError(t, err)
	err = db.RevokeSession(ctx, 6)
	require.NoError(t, err)

	var rss1, rss2 []app.Session
	rss1 = append(rss1, ss1[0:2]...)
	rss1 = append(rss1, ss1[4])
	rss2 = append(rss2, ss2[1:]...)

	t.Run("GetFor1ExcludeInvalid", func(t *testing.T) {
		ss, err = db.GetSessionsForPerson(ctx, 1, false)
		require.NoError(t, err)

		sortSessions(ss)
		assertEqualJSON(t, rss1, ss)
	})

	t.Run("GetFor2ExcludeInvalid", func(t *testing.T) {
		ss, err = db.GetSessionsForPerson(ctx, 2, false)
		require.NoError(t, err)

		sortSessions(ss)
		assertEqualJSON(t, rss2, ss)
	})

	t.Run("NoSessions", func(t *testing.T) {
		ss, err = db.GetSessionsForPerson(ctx, 38, false)
		require.NoError(t, err)
		require.Len(t, ss, 0)
	})
}
