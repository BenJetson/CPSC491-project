package api

import (
	"net/http"

	"github.com/pkg/errors"
)

func (svr *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	svr.sendErrorResponse(w,
		errors.Errorf("requested route of %s does not exist", r.URL.Path),
		http.StatusNotFound, "")
}

func (svr *Server) handleMethodNotAllowed(
	w http.ResponseWriter,
	r *http.Request,
) {

	svr.sendErrorResponse(w,
		errors.Errorf("wrong method (%s) for route %s",
			r.Method, r.URL.Path),
		http.StatusMethodNotAllowed,
		"Cannot %s this endpoint.", r.Method,
	)
}
