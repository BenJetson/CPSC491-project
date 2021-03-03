package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

			s, err := svr.db.GetSessionByToken(token)
			if err != nil {
				// FIXME need a way to detect not found errors and return 401.
				if err != nil {
					svr.sendErrorResponse(
						w,
						errors.Wrap(err, "failed to retrieve session"),
						http.StatusInternalServerError,
						"",
					)
					return
				}
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

func (svr *Server) GetSessionFromContext(ctx context.Context) *app.Session {
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
func (svr *Server) requireAuth(cfg authConfig) mux.MiddlewareFunc {
	// This will check the authConfig parameter values ONCE when the router is
	// first instantiated and the app will not be allowed to start if there is
	// a parameter mismatch.
	if (!cfg.requireRole && len(cfg.allowedRoles) > 0) ||
		(cfg.requireRole && len(cfg.allowedRoles) < 1) {
		panic("authConfig parameter mismatch")
	}

	return func(next http.Handler) http.Handler {
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
						errors.Errorf("endpoint does not allow role %v", s.Person.Role),
						http.StatusForbidden,
						"",
					)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
