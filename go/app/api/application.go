package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

type applicationApprovalRequest struct {
	IsApproved    bool   `json:"is_approved"`
	Reason        string `json:"reason"`
	ApplicationID int    `json:"application_id"`
}

type applicationSubmissionRequest struct {
	OrganizationID int    `json:"organization_id"`
	Comment        string `json:"comment"`
}

func (svr *Server) handleSubmitApplication(
	w http.ResponseWriter,
	r *http.Request,
) {

	s := getSessionFromContext(r.Context())
	if s == nil {
		svr.sendErrorResponse(w,
			errors.New("missing session for application"),
			http.StatusInternalServerError, "")
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var appReq applicationSubmissionRequest
	if err := d.Decode(&appReq); err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to decode application request"),
			http.StatusBadRequest, "Invalid applicaion format.")
		return
	}

	_, err := svr.db.CreateApplication(r.Context(), app.Application{
		ApplicantID:    s.Person.ID,
		OrganizationID: appReq.OrganizationID,
		Comment:        appReq.Comment,
	})

	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to create application"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

	a, err := svr.db.GetApplicationByID(r.Context(), appID)
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

	orgID, err := getOrganizationIDOfSponsor(r)
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "cannot determine sponsor organization identity"),
			http.StatusInternalServerError, "")
		return
	}

	apps, err := svr.db.GetApplicationsForOrganization(r.Context(), orgID)
	if err != nil {
		svr.sendErrorResponse(
			w,
			errors.Wrap(err, "failed to retrieve application"),
			http.StatusInternalServerError,
			"",
		)
		return
	}

	if apps == nil {
		apps = make([]app.Application, 0)
	}

	svr.sendJSONResponse(w, apps)
}

func (svr *Server) handleApproveApplication(
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

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var data applicationApprovalRequest
	if err = d.Decode(&data); err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "received bad json data"),
			http.StatusBadRequest, "Bad JSON data.")
		return
	}

	app, err := svr.db.GetApplicationByID(r.Context(), appID)
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "could not get application"),
			http.StatusInternalServerError, "")
		return
	}

	if svr.requireOrganization(orgConfig{orgID: app.OrganizationID}, w, r) {
		return
	}

	err = svr.db.UpdateApplicationApproval(r.Context(),
		data.ApplicationID, data.IsApproved, data.Reason)
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to approve app"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (svr *Server) handleGetMyApplications(
	w http.ResponseWriter,
	r *http.Request,
) {

	s := getSessionFromContext(r.Context())
	apps, err := svr.db.GetApplicationsForPerson(r.Context(), s.Person.ID)
	if err != nil {
		svr.sendErrorResponse(w, err, http.StatusInternalServerError, "")
		return
	}

	if apps == nil {
		apps = make([]app.Application, 0)
	}

	svr.sendJSONResponse(w, apps)
}
