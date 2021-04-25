package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BenJetson/CPSC491-project/go/app"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type applicationRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
}

func (svr *Server) handleSubmitApplication(
	w http.ResponseWriter,
	r *http.Request,
) {
	s := getSessionFromContext(r.Context())
	if s != nil {
		svr.sendErrorResponse(w,
			errors.New("cannot submit new application"),
			http.StatusUnauthorized, "Cannot submit new application.")
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var appReq applicationRequest
	if err := d.Decode(&appReq); err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to decode application request"),
			http.StatusBadRequest, "Invalid applicaion format.")
		return
	}

	svr.sendJSONResponse(w, appReq)
}

func (svr *Server) handleGetApplicationByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	pathParams := mux.Vars(r)

	appID, err := strconv.Atoi(pathParams["appID"])
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "appID must be an integer"),
			http.StatusBadRequest, "Application ID must be an integer.")
		return
	}
	a, err := svr.db.GetApplicationsForPerson(r.Context(), appID)
	if errors.Is(err, app.ErrNotFound) {
		svr.sendErrorResponse(w,
			errors.Wrapf(err, "no application with ID of %d", appID),
			http.StatusNotFound, "No such application.")
		return
	} else if err != nil {
		svr.sendErrorResponse(
			w,
			errors.Wrap(err, "failed to retrieve application"),
			http.StatusInternalServerError,
			"",
		)
		return
	}

	svr.sendJSONResponse(w, a)
}

func (svr *Server) handleGetApplicationsForOrganization(
	w http.ResponseWriter,
	r *http.Request,
) {

	pathParams := mux.Vars(r)

	appID, err := strconv.Atoi(pathParams["appID"])
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "appID must be an integer"),
			http.StatusBadRequest, "Application ID must be an integer.")
		return
	}
	o, err := svr.db.GetApplicationsForOrganization(r.Context(), appID)
	if errors.Is(err, app.ErrNotFound) {
		svr.sendErrorResponse(w,
			errors.Wrapf(err, "no application with ID of %d", appID),
			http.StatusNotFound, "No such application.")
		return
	} else if err != nil {
		svr.sendErrorResponse(
			w,
			errors.Wrap(err, "failed to retrieve application"),
			http.StatusInternalServerError,
			"",
		)
		return
	}

	svr.sendJSONResponse(w, o)
}

func (svr *Server) handleApproveApplication(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.WriteHeader(http.StatusNotImplemented) // TODO
}

func (svr *Server) handleGetOrganizations(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.WriteHeader(http.StatusNotImplemented) // TODO
}

func (svr *Server) handleGetMyApplications(
	w http.ResponseWriter,
	r *http.Request,
) {
	s := getSessionFromContext(r.Context())
	applications, err := svr.db.GetApplicationsForPerson(r.Context(), s.Person.ID)
	if err != nil {
		svr.sendErrorResponse(w, err, http.StatusInternalServerError, "")
		return
	}
	svr.sendJSONResponse(w, applications)
}
