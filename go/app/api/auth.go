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

type authConfig struct {
	requireUser  bool
	requireRole  bool
	allowedRoles []int
}

func requireAuth(cfg authConfig) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		})
	}
}
