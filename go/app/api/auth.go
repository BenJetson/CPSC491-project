package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

type contextKey string

const contextKeySession contextKey = "session"

func (svr *Server) authContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(sessionCookieKey)
		if err == nil {
			token, err := uuid.Parse(c.Value)
			if err != nil {
				svr.sendErrorResponse(
					w,
					errors.Wrap(err, "bad login token"),
					http.StatusUnauthorized,
					"",
				)
				return
			}

			s, err := svr.db.GetSessionByToken(r.Context(), token)
			if errors.Is(err, app.ErrNotFound) {
				svr.logger.
					WithField("token", token).
					Info("received unknown session token; destroying cookie")
				svr.destroySessionCookie(w)

				// Attach an invalid session so it is not attached to context.
				s = app.Session{IsRevoked: true}
			} else if err != nil {
				svr.sendErrorResponse(
					w,
					errors.Wrap(err, "failed to retrieve session"),
					http.StatusInternalServerError,
					"",
				)
				return
			}

			if s.IsValid() {
				// Attach session to context; attach new context to request.
				ctx := context.WithValue(r.Context(), contextKeySession, s)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}

// getSessionFromContext retrieves the session object from the request context
// if one is available, and returns nil otherwise.
//
// This only works if authContextMiddleware has already run for this request.
func getSessionFromContext(ctx context.Context) *app.Session {
	s, ok := ctx.Value(contextKeySession).(app.Session)
	if !ok {
		return nil
	}
	return &s
}

// An authConfig specifies what authentication parameters are required for an
// endpoint.
// nolint: unused // FIXME: remove this once it gets used
type authConfig struct {
	// requireRole determines whether or not required roles are enforced.
	requireRole bool
	// allowedRoles is the list of roles allowed to access this endpoint.
	// Must be blank when !requireRole and have at least one value otherwise.
	allowedRoles []app.Role
}

// requireAuth is a middleware that may be applied to a route or subrouter that
// will enforce authentication requirements.
// nolint: unused // FIXME: remove this once it gets used
func (svr *Server) requireAuth(
	cfg authConfig,
	handler http.HandlerFunc,
) http.HandlerFunc {

	// This will check the authConfig parameter values ONCE when the router is
	// first instantiated and the app will not be allowed to start if there is
	// a parameter mismatch.
	if (!cfg.requireRole && len(cfg.allowedRoles) > 0) ||
		(cfg.requireRole && len(cfg.allowedRoles) < 1) {
		panic("authConfig parameter mismatch")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := getSessionFromContext(r.Context())
		if s == nil {
			svr.sendErrorResponse(
				w,
				errors.New("endpoint requires auth but nil session"),
				http.StatusUnauthorized,
				"",
			)
			return
		}

		if cfg.requireRole {
			hasRole := false
			for _, r := range cfg.allowedRoles {
				if r == s.Person.Role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				svr.sendErrorResponse(
					w,
					errors.Errorf(
						"endpoint does not allow role %v",
						s.Person.Role,
					),
					http.StatusForbidden,
					"",
				)
				return
			}
		}

		// User passed the auth check. Call the handler.
		handler(w, r)
	})
}

// nolint: unused // TODO remove this when it gets used
type identityConfig struct {
	// personID is the identifier for the target Person.
	personID int
	// adminOverride should be true when admins may override this identity
	// requirement.
	adminOverride bool
	// sponsorOverride should be true when sponsors OF THE TARGET PERSON
	sponsorOverride bool
}

// requireIdentity may be used within a handler to guard an operation when only
// a particular person may perform that operation.
//
// Suggested usage:
//
//     if svr.requireIdentity(identityConfig{personID: x}, w, r) {
// 	   	   return
//     }
//
// Upon failure of this check, this method will write an appropriate
// authorization or internal server error to the ResponseWriter for you.
//
// nolint: unused, gocyclo // should keep all identity code here for security.
func (svr *Server) requireIdentity(
	cfg identityConfig,
	w http.ResponseWriter,
	r *http.Request,
) bool {

	s := getSessionFromContext(r.Context())
	if s == nil {
		// FAIL: No login session was detected.
		svr.sendErrorResponse(
			w,
			errors.New("endpoint requires auth but nil session"),
			http.StatusUnauthorized,
			"",
		)
		return true
	}

	if s.Person.ID == cfg.personID {
		// PASS: Current user's Person ID matches required identity.
		return false
	}

	if cfg.adminOverride && s.Person.Role == app.RoleAdmin {
		// PASS: Current user's Person has role of admin, and admins are
		// allowed to override this identity requirement.
		return false
	}

	if cfg.sponsorOverride && s.Person.Role == app.RoleSponsor {
		// We must determine which organizations the user in the config is
		// sponsored by, then check to see if the current user is a sponsor
		// for that organization.

		p, err := svr.db.GetPersonByID(r.Context(), cfg.personID)
		if err != nil {
			// FAIL: could not verify identity due to database problem.
			svr.sendErrorResponse(
				w,
				errors.Wrap(err, "failed to retrieve session"),
				http.StatusInternalServerError,
				"",
			)
			return true
		}

		if p.Role == app.RoleDriver {
			// Only drivers have sponsors, so this is the only time when a
			// sponsor override may apply.

			if len(s.Person.Affiliations) != 1 {
				// FAIL: sponsor should only have one organization.
				svr.sendErrorResponse(
					w,
					errors.Errorf(
						"sponsor (person id %d) has multiple affiliations (%d)",
						s.Person.ID,
						len(s.Person.Affiliations),
					),
					http.StatusInternalServerError,
					"",
				)
				return true
			}

			target := s.Person.Affiliations[0]
			for _, orgID := range p.Affiliations {
				if target == orgID {
					// PASS: Current user is a sponsor and the required identity
					// is affiliated with the organization they have a sponsor
					// role within.
					return false
				}
			}
		}
	}

	// FINAL FAIL: none of the conditions were met.
	svr.sendErrorResponse(
		w,
		errors.Errorf("user does not match required identity: %+v", cfg),
		http.StatusForbidden,
		"",
	)
	return true
}

// nolint: unused // TODO remove this when it gets used
type orgConfig struct {
	// orgID is the target organization ID that the current user must be
	// affiliated with for this check to pass.
	orgID int
	// adminOverride should be true when administrators may override this
	// organization requirement.
	adminOverride bool
}

// requireOrganization may be used within a handler to guard a particular
// operation with an organization affiliation requirement when only an affiliate
// of that organization may perform that operation.
//
// Suggested usage:
//
//     if svr.requireOrganization(orgConfig{orgID: x}, w, r) {
// 	   	   return
//     }
//
// Upon failure of this check, this method will write an appropriate
// authorization or internal server error to the ResponseWriter for you.
//
// nolint: unused // TODO remove this when it gets used
func (svr *Server) requireOrganization(
	cfg orgConfig,
	w http.ResponseWriter,
	r *http.Request,
) bool {

	s := getSessionFromContext(r.Context())
	if s == nil {
		// FAIL: No login session was detected.
		svr.sendErrorResponse(
			w,
			errors.New("endpoint requires auth but nil session"),
			http.StatusUnauthorized,
			"",
		)
		return true
	}

	if cfg.adminOverride && s.Person.Role == app.RoleAdmin {
		// PASS: Current user's Person has role of admin, and admins are
		// allowed to override this organization requirement.
		return false
	}

	for _, orgID := range s.Person.Affiliations {
		if orgID == cfg.orgID {
			// PASS: Current user's Person is affiliated with the target
			// organization.
			return true
		}
	}

	// FINAL FAIL: none of the conditions were met.
	svr.sendErrorResponse(
		w,
		errors.Errorf("user does not match required organization: %+v", cfg),
		http.StatusForbidden,
		"",
	)
	return false
}
