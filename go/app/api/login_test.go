package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/BenJetson/CPSC491-project/go/app"
	"github.com/BenJetson/CPSC491-project/go/app/mock"
)

type loginMockDB struct {
	*mock.DB

	personByEmail    app.Person
	personByEmailErr error

	sessionByToken    app.Session
	sessionByTokenErr error

	createSessionErr error

	revokeSessionErr error
}

func (db *loginMockDB) GetPersonByEmail(
	_ context.Context,
	_ string,
) (app.Person, error) {

	return db.personByEmail, db.personByEmailErr
}

func (db *loginMockDB) GetSessionByToken(
	_ context.Context,
	_ uuid.UUID,
) (app.Session, error) {

	return db.sessionByToken, db.sessionByTokenErr
}

func (db *loginMockDB) CreateSession(
	_ context.Context,
	_ app.Session,
) (int, error) {

	return 0, db.createSessionErr
}

func (db *loginMockDB) RevokeSession(_ context.Context, _ int) error {
	return db.revokeSessionErr
}

func testSessionTokenInject(_ *testing.T, r *http.Request, token uuid.UUID) {
	if token == uuid.Nil {
		return
	}

	c := http.Cookie{
		Name:  sessionCookieKey,
		Value: token.String(),
	}

	r.Header.Add("Cookie", c.String())
}

func TestHandleLogin(t *testing.T) {
	db := &loginMockDB{}
	api, _, _ := newTestAPI(t, db, nil)

	pass, err := app.NewPassword("zxcvbnJKL")
	require.NoError(t, err)

	p := app.Person{
		FirstName: "Billy Joe",
		LastName:  "Bob",
		Email:     "jack@box.net",
		Password:  pass,
	}

	testCases := []struct {
		alias               string
		body                string
		nilBody             bool
		dbPersonByEmail     app.Person
		dbPersonByEmailErr  error
		dbCreateSessionErr  error
		expectCode          int
		expectCookie        bool
		expectDestroyCookie bool
	}{
		{
			alias:      "NilBody",
			nilBody:    true,
			expectCode: http.StatusBadRequest,
		},
		{
			alias:      "BlankBody",
			body:       ``,
			expectCode: http.StatusBadRequest,
		},
		{
			alias:      "NotJSON",
			body:       `space trash <!--X @@@#!~==>()()@)@! aaaaaack.!`,
			expectCode: http.StatusBadRequest,
		},
		{
			alias: "GarbageJSON",
			body: `
				{
					"junk": true,
					"aack?": ["1", "23", null],
					"qwerty": "uiop",
					"asdf": null
				}
			`,
			expectCode: http.StatusBadRequest,
		},
		{
			alias: "NoSuchEmail",
			body: `
				{
					"email": "jack@box.net",
					"password": "p@$$w0rd=ye$"
				}
			`,
			dbPersonByEmailErr: errors.Wrap(app.ErrNotFound, "not so fast"),
			expectCode:         http.StatusUnauthorized,
		},
		{
			alias: "BlankEmail",
			body: `
				{
					"email": "",
					"password": "zxcvbnJKL"
				}
			`,
			expectCode: http.StatusBadRequest,
		},
		{
			alias: "WrongPassword",
			body: `
				{
					"email": "jack@box.net",
					"password": "p@$$w0rd=ye$"
				}
			`,
			dbPersonByEmail: p,
			expectCode:      http.StatusUnauthorized,
		},
		{
			alias: "BlankPassword",
			body: `
				{
					"email": "jack@box.net",
					"password": ""
				}
			`,
			expectCode: http.StatusBadRequest,
		},
		{
			alias: "GetPersonError",
			body: `
				{
					"email": "jack@box.net",
					"password": "zxcvbnJKL"
				}
			`,
			dbPersonByEmailErr: errors.New("whoops lost a few bits there"),
			expectCode:         http.StatusInternalServerError,
		},
		{
			alias: "CreateSessionError",
			body: `
				{
					"email": "jack@box.net",
					"password": "zxcvbnJKL"
				}
			`,
			dbPersonByEmail:    p,
			dbCreateSessionErr: errors.New("please hang up and try again"),
			expectCode:         http.StatusInternalServerError,
		},
		{
			alias: "Success",
			body: `
				{
					"email": "jack@box.net",
					"password": "zxcvbnJKL"
				}
			`,
			dbPersonByEmail: p,
			expectCode:      http.StatusOK,
			expectCookie:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.alias, func(t *testing.T) {
			db.personByEmail = tc.dbPersonByEmail
			db.personByEmailErr = tc.dbPersonByEmailErr
			db.createSessionErr = tc.dbCreateSessionErr

			var body io.Reader
			if !tc.nilBody {
				body = strings.NewReader(tc.body)
			}

			r := httptest.NewRequest("POST", "/login", body)
			w := httptest.NewRecorder()

			api.router.ServeHTTP(w, r)

			assert.Equal(t, tc.expectCode, w.Code)

			res := w.Result()
			defer res.Body.Close()

			jar := res.Cookies()

			var loginCookieIsSet, willDestroyLoginCookie bool
			for _, c := range jar {
				if c.Name == sessionCookieKey {
					loginCookieIsSet = true

					if c.MaxAge < 0 {
						willDestroyLoginCookie = true
					}

					break
				}
			}

			assert.Equal(
				t,
				tc.expectCookie || tc.expectDestroyCookie,
				loginCookieIsSet,
				"cookie presence expectation mismatch",
			)
			assert.Equal(t, tc.expectDestroyCookie, willDestroyLoginCookie,
				"cookie destruction expectation mismatch")
		})
	}
}

func TestHandleLogout(t *testing.T) {
	db := &loginMockDB{}
	api, _, _ := newTestAPI(t, db, nil)

	pass, err := app.NewPassword("zxcvbnJKL")
	require.NoError(t, err)

	p := app.Person{
		FirstName: "Billy Joe",
		LastName:  "Bob",
		Email:     "jack@box.net",
		Password:  pass,
	}

	s, err := app.NewSession(p)
	require.NoError(t, err)
	require.NotNil(t, s)
	require.NotEqual(t, uuid.Nil, s.Token)

	testCases := []struct {
		alias               string
		sessionToken        uuid.UUID
		dbSessionByToken    app.Session
		dbSessionByTokenErr error
		dbRevokeSessionErr  error
		expectCode          int
		expectCookie        bool
		expectDestroyCookie bool
	}{
		{
			alias:              "RevokeSessionError",
			sessionToken:       s.Token,
			dbSessionByToken:   *s,
			dbRevokeSessionErr: errors.New("please file cherwell ticket"),
			expectCode:         http.StatusInternalServerError,
		},
		// {
		// 	alias:               "SessionByTokenError",
		// 	sessionToken:        s.Token,
		// 	dbSessionByToken:    *s,
		// 	dbSessionByTokenErr: errors.New("aaaack! failed with fire"),
		// 	expectCode:          http.StatusInternalServerError,
		// },
		// {
		// 	alias:               "NoSessionButRevokeErr",
		// 	dbRevokeSessionErr:  errors.New("this should not matter"),
		// 	expectCode:          http.StatusOK,
		// 	expectCookie:        true,
		// 	expectDestroyCookie: true,
		// },
		// {
		// 	alias:               "NoSession",
		// 	expectCode:          http.StatusOK,
		// 	expectCookie:        true,
		// 	expectDestroyCookie: true,
		// },
		// {
		// 	alias:               "Success",
		// 	sessionToken:        s.Token,
		// 	dbSessionByToken:    *s,
		// 	expectCode:          http.StatusOK,
		// 	expectCookie:        true,
		// 	expectDestroyCookie: true,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.alias, func(t *testing.T) {
			db.revokeSessionErr = tc.dbRevokeSessionErr
			db.sessionByToken = tc.dbSessionByToken
			db.sessionByTokenErr = tc.dbSessionByTokenErr

			r := httptest.NewRequest("POST", "/logout", nil)
			w := httptest.NewRecorder()

			testSessionTokenInject(t, r, tc.sessionToken)

			api.router.ServeHTTP(w, r)

			assert.Equal(t, tc.expectCode, w.Code)

			res := w.Result()
			defer res.Body.Close()

			jar := res.Cookies()

			var loginCookieIsSet, willDestroyLoginCookie bool
			for _, c := range jar {
				if c.Name == sessionCookieKey {
					loginCookieIsSet = true

					if c.MaxAge < 0 {
						willDestroyLoginCookie = true
					}

					break
				}
			}

			assert.Equal(
				t,
				tc.expectCookie || tc.expectDestroyCookie,
				loginCookieIsSet,
				"cookie presence expectation mismatch",
			)
			assert.Equal(t, tc.expectDestroyCookie, willDestroyLoginCookie,
				"cookie destruction expectation mismatch")
		})
	}
}
