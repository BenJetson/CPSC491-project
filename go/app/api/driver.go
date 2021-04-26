package api

import (
	"net/http"

	"github.com/pkg/errors"
)

func (svr *Server) handleDriverGetBalances(
	w http.ResponseWriter,
	r *http.Request,
) {

	s := getSessionFromContext(r.Context())
	if s == nil {
		svr.sendErrorResponse(w, errors.New("missing session for balance"),
			http.StatusInternalServerError, "")
		return
	}

	bs, err := svr.db.GetBalancesForPerson(r.Context(), s.Person.ID)
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "balance fetch failed"),
			http.StatusInternalServerError, "")
		return
	}

	svr.sendJSONResponse(w, bs)
}
