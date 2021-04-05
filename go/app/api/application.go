package api

import "net/http"

func (svr *Server) handleSubmitApplication(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.WriteHeader(http.StatusNotImplemented) // TODO
}

func (svr *Server) handleGetApplicationByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.WriteHeader(http.StatusNotImplemented) // TODO
}

func (svr *Server) handleGetApplicationsForPerson(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.WriteHeader(http.StatusNotImplemented) // TODO
}

func (svr *Server) handleGetApplicationsForOrganization(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.WriteHeader(http.StatusNotImplemented) // TODO
}

func (svr *Server) handleApproveApplication(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.WriteHeader(http.StatusNotImplemented) // TODO
}
