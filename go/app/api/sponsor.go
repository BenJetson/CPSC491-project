package api

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

func (svr *Server) handleSponsorGetOwnOrganization(
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

	org, err := svr.db.GetOrganizationByID(r.Context(), orgID)
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to fetch sponsor organization"),
			http.StatusInternalServerError, "")
		return
	}

	svr.sendJSONResponse(w, org)
}

type organizationRequest struct {
	Name       string    `json:"name"`
	PointValue app.Money `json:"point_value"`
}

func (r *organizationRequest) validateFields() (message string, err error) {
	defer func() {
		if message != "" {
			err = errors.New(message)
		}
	}()

	if len(r.Name) < 1 {
		message = "Name cannot be blank."
		return
	}

	if r.PointValue < 1 {
		message = "Point Value cannot be less than one cent."
		return
	}

	return
}

func (svr *Server) handleSponsorUpdateOwnOrganization(
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

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var data organizationRequest
	var message string
	if err = d.Decode(&data); err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "received bad json data"),
			http.StatusBadRequest, "Bad JSON data.")
		return
	} else if message, err = data.validateFields(); err != nil {
		svr.sendErrorResponse(w, err, http.StatusBadRequest, message)
		return
	}

	err = svr.db.UpdateOrganization(r.Context(), app.Organization{
		ID:         orgID,
		Name:       data.Name,
		PointValue: data.PointValue,
	})
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to update sponsor organization"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
