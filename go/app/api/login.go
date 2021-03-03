package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

const sessionCookieKey = "SESSION_TOKEN"

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (svr *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var credentials loginRequest
	if err := d.Decode(&credentials); err != nil {
		svr.sendErrorResponse(
			w,
			errors.Wrap(err, "failed to decode login request"),
			http.StatusBadRequest,
			"invalid login request format",
		)
		return
	}

	p, err := svr.db.GetPersonByEmail(credentials.Email)
	// FIXME need a way to detect not found errors and return 401.
	if err != nil {
		svr.sendErrorResponse(
			w,
			errors.Wrap(err, "failed to retrieve person"),
			http.StatusInternalServerError,
			"",
		)
		return
	}

	if !p.Password.Verify(credentials.Password) {
		svr.sendErrorResponse(
			w,
			errors.Wrap(err, "password did not match"),
			http.StatusUnauthorized,
			"Email address or password was incorrect.",
		)
		return
	}

	var s *app.Session
	if s, err = app.NewSession(p); err != nil {
		svr.sendErrorResponse(
			w,
			errors.Wrap(err, "failed to create new login session"),
			http.StatusInternalServerError,
			"",
		)
		return
	} else if err = svr.db.CreateSession(*s); err != nil {
		svr.sendErrorResponse(
			w,
			errors.Wrap(err, "failed to store login session"),
			http.StatusInternalServerError,
			"",
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	http.SetCookie(w, &http.Cookie{
		Name:  sessionCookieKey,
		Value: s.Token.String(),

		// Ensure that this cookie is only used on the same domain with the
		// same protocol.
		Domain:   svr.hostname(),
		SameSite: http.SameSiteStrictMode,
		Secure:   svr.useHTTPS(),

		// HttpOnly hides this cookie from JavaScript in browsers for security.
		HttpOnly: true,

		// Sessions expire on our app server-side, but let's ask the client to
		// ditch the cookie automatically as well.
		MaxAge: int(app.SessionLength / time.Second),
	})
}

func (svr *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	s := svr.GetSessionFromContext(r.Context())
	if s != nil {
		if err := svr.db.RevokeSession(s.Token); err != nil {
			svr.sendErrorResponse(
				w,
				errors.Wrap(err, "failed to revoke session"),
				http.StatusInternalServerError,
				"",
			)
			return
		}
	}

	// Destroy the session cookie on the client.
	w.WriteHeader(http.StatusOK)
	http.SetCookie(w, &http.Cookie{
		Name: sessionCookieKey,

		// Ensure that this cookie is only used on the same domain with the
		// same protocol.
		Domain:   svr.hostname(),
		SameSite: http.SameSiteStrictMode,
		Secure:   svr.useHTTPS(),

		// HttpOnly hides this cookie from JavaScript in browsers for security.
		HttpOnly: true,

		// A MaxAge less than zero will cause clients to destroy this cookie.
		MaxAge: -1,
	})
}
