package api

import (
	"errors"
	"net/http"
)

func (svr *Server) handleTODO(w http.ResponseWriter, _ *http.Request) {
	svr.sendErrorResponse(w, errors.New("not implemented"),
		http.StatusNotImplemented, "")
}
