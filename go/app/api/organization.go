package api

import (
	"net/http"

	"github.com/pkg/errors"
)

func (svr *Server) handleGetAllOrganizations(
	w http.ResponseWriter,
	r *http.Request,
) {

	orgs, err := svr.db.GetAllOrganizations(r.Context())
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to get organizations"),
			http.StatusInternalServerError, "")
		return
	}

	svr.sendJSONResponse(w, orgs)
}
