package api

import (
	"net/http"
	"runtime/debug"

	"github.com/pkg/errors"
)

func (svr *Server) panicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			cause := recover()
			if cause == nil {
				return
			}

			err := errors.Errorf("panic caused by: %s", cause)
			if cause, isError := cause.(error); isError {
				err = errors.Wrap(cause, "error caused panic")
			}

			svr.logger.
				WithError(err).
				WithField("stack", string(debug.Stack())).
				Error("panic when handling api request; recovering")

			svr.sendErrorResponse(w, err, http.StatusInternalServerError, "")
		}()

		next.ServeHTTP(w, r)
	})
}
