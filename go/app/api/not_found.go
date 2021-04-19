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
