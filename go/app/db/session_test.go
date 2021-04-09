package db

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func TestGetSessionByToken(t *testing.T) {
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

	t.Run("NoSuchSession", func(t *testing.T) {
		_, err = db.GetSessionByToken(ctx, uuid.Nil)
		require.Error(t, err)
		assert.True(t, errors.Is(err, app.ErrNotFound))
	})

	var actual app.Session

	// nolint: scopelint // allow for test purposes.
	for groupidx, ss := range [][]app.Session{ss1, ss2} {
		for idx, expect := range ss {
			t.Run(fmt.Sprintf("ss%d-%d", groupidx+1, idx), func(t *testing.T) {
				actual, err = db.GetSessionByToken(ctx, expect.Token)
				require.NoError(t, err)
				assertEqualJSON(t, expect, actual)
			})
		}
	}
}

func TestCreateSession(t *testing.T) {
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

	var s *app.Session
	var ss []app.Session

	assertSessionPresence := func(t *testing.T) {
		db.assertCount(t, "session", len(ss))

		for _, s := range ss {
			db.assertCountOf(t, "session", 1, `
				token = $1
				AND person_id = $2
				AND created_at = $3::timestamptz
				AND expires_at = $4::timestamptz
				AND session_id = $5
			`, s.Token, s.Person.ID, s.CreatedAt, s.ExpiresAt, s.ID)
		}
	}

	t.Run("NoSuchPerson", func(t *testing.T) {
		s, err = app.NewSession(app.Person{})
		require.NoError(t, err)
		require.NotNil(t, s)

		_, err = db.CreateSession(ctx, *s)
		require.Error(t, err)

		assertSessionPresence(t)
	})

	for i := 0; i < 20; i++ {
		i := i // pin this value for scopelint

		t.Run(fmt.Sprintf("total%d", i+1), func(t *testing.T) {
			p := p1
			if i%2 == 0 {
				p = p2
			}

			s, err = app.NewSession(p)
			require.NoError(t, err)
			require.NotNil(t, s)

			s.ID, err = db.CreateSession(ctx, *s)
			require.NoError(t, err)
			assert.NotZero(t, s.ID)

			ss = append(ss, *s)

			assertSessionPresence(t)
		})
	}
}

func TestRevokeSession(t *testing.T) {
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

	var s *app.Session
	var ss []app.Session

	assertSessionPresence := func(t *testing.T) {
		db.assertCount(t, "session", len(ss))

		for _, s := range ss {
			db.assertCountOf(t, "session", 1, `
				token = $1
				AND person_id = $2
				AND created_at = $3::timestamptz
				AND expires_at = $4::timestamptz
				AND session_id = $5
			`, s.Token, s.Person.ID, s.CreatedAt, s.ExpiresAt, s.ID)
		}
	}

	t.Run("NoSuchSession", func(t *testing.T) {
		err = db.RevokeSession(ctx, 881)
		require.Error(t, err)
		assert.True(t, errors.Is(err, app.ErrNotFound))

		assertSessionPresence(t)
	})

	// Create sessions
	for i := 0; i < 18; i++ {
		p := p1
		if i%2 == 0 {
			p = p2
		}

		s, err = app.NewSession(p)
		require.NoError(t, err)
		require.NotNil(t, s)

		s.ID, err = db.CreateSession(ctx, *s)
		require.NoError(t, err)
		assert.NotZero(t, s.ID)

		ss = append(ss, *s)

		assertSessionPresence(t)
	}

	for i := 0; i < 18; i += 3 {
		i := i // pin this value for scopelint

		t.Run(fmt.Sprintf("revoke#%d", i+1), func(t *testing.T) {
			err = db.RevokeSession(ctx, ss[i].ID)
			ss[i].IsRevoked = true

			assertSessionPresence(t)
		})
	}
}

func TestRevokeSessionsForPersonExcept(t *testing.T) {
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

	var s *app.Session
	var ss1, ss2 []app.Session

	assertSessionPresence := func(t *testing.T) {
		db.assertCount(t, "session", len(ss1)+len(ss2))

		for _, s := range append(ss1, ss2...) {
			db.assertCountOf(t, "session", 1, `
				token = $1
				AND person_id = $2
				AND created_at = $3::timestamptz
				AND expires_at = $4::timestamptz
				AND session_id = $5
			`, s.Token, s.Person.ID, s.CreatedAt, s.ExpiresAt, s.ID)
		}
	}

	// Create sessions
	for i := 0; i < 57; i++ {
		p := p1
		if i%2 == 0 {
			p = p2
		}

		s, err = app.NewSession(p)
		require.NoError(t, err)
		require.NotNil(t, s)

		s.ID, err = db.CreateSession(ctx, *s)
		require.NoError(t, err)
		assert.NotZero(t, s.ID)

		if i%2 == 0 {
			ss2 = append(ss2, *s)
		} else {
			ss1 = append(ss1, *s)
		}

		assertSessionPresence(t)
	}

	t.Run("NoSuchPersonOrSession", func(t *testing.T) {
		err = db.RevokeSessionsForPersonExcept(ctx, 92, 482)
		require.NoError(t, err, "this should not be an error condition, "+
			"since person may have had zero sessions")

		assertSessionPresence(t)
	})

	t.Run("RevokePerson1", func(t *testing.T) {
		for i := range ss1 {
			if i != 9 {
				ss1[i].IsRevoked = true
			}
		}

		err = db.RevokeSessionsForPersonExcept(ctx, p1.ID, ss1[9].ID)
		require.NoError(t, err)

		assertSessionPresence(t)
	})

	t.Run("RevokePerson2", func(t *testing.T) {
		for i := range ss2 {
			if i != 18 {
				ss2[i].IsRevoked = true
			}
		}

		err = db.RevokeSessionsForPersonExcept(ctx, p2.ID, ss2[18].ID)
		require.NoError(t, err)

		assertSessionPresence(t)
	})
}
