package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

func (svr *Server) handleAdminGetAllUsers(
	w http.ResponseWriter,
	r *http.Request,
) {

	people, err := svr.db.GetAllPeople(r.Context())
	if err != nil {
		svr.sendErrorResponse(w, err, http.StatusInternalServerError, "")
		return
	}

	svr.sendJSONResponse(w, people)
}

func (svr *Server) handleAdminGetUserByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	pathParams := mux.Vars(r)

	userID, err := strconv.Atoi(pathParams["userID"])
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "userID must be an integer"),
			http.StatusBadRequest, "User ID must be an integer.")
		return
	}

	p, err := svr.db.GetPersonByID(r.Context(), userID)
	if errors.Is(err, app.ErrNotFound) {
		svr.sendErrorResponse(w,
			errors.Wrapf(err, "no person with ID of %d", userID),
			http.StatusNotFound, "No such user.")
		return
	} else if err != nil {
		svr.sendErrorResponse(
			w,
			errors.Wrap(err, "failed to retrieve person"),
			http.StatusInternalServerError,
			"",
		)
		return
	}

	svr.sendJSONResponse(w, p)
}

func (svr *Server) handleAdminUpdateUserName(
	w http.ResponseWriter,
	r *http.Request,
) {

	pathParams := mux.Vars(r)

	userID, err := strconv.Atoi(pathParams["userID"])
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "userID must be an integer"),
			http.StatusBadRequest, "User ID must be an integer.")
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var data nameChangeRequest
	var message string
	if err = d.Decode(&data); err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "received bad json data"),
			http.StatusBadRequest, "Bad JSON data.")
		return
	} else if message, err = data.validateFields(); err != nil {
		svr.sendErrorResponse(w, err, http.StatusBadRequest, message)
		return
	}

	err = svr.db.UpdatePersonName(r.Context(),
		userID, data.FirstName, data.LastName)
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "failed to update name"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (svr *Server) handleAdminUpdateUserEmail(
	w http.ResponseWriter,
	r *http.Request,
) {

	pathParams := mux.Vars(r)

	userID, err := strconv.Atoi(pathParams["userID"])
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "userID must be an integer"),
			http.StatusBadRequest, "User ID must be an integer.")
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var data emailChangeRequest
	var message string
	if err = d.Decode(&data); err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "received bad json data"),
			http.StatusBadRequest, "Bad JSON data.")
		return
	} else if message, err = data.validateFields(); err != nil {
		svr.sendErrorResponse(w, err, http.StatusBadRequest, message)
		return
	}

	err = svr.db.UpdatePersonEmail(r.Context(), userID, data.Email)
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "failed to update email"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (svr *Server) handleAdminUpdateUserPassword(
	w http.ResponseWriter,
	r *http.Request,
) {

	pathParams := mux.Vars(r)

	userID, err := strconv.Atoi(pathParams["userID"])
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "userID must be an integer"),
			http.StatusBadRequest, "User ID must be an integer.")
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var data passwordChangeRequest
	var message string
	if err = d.Decode(&data); err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "received bad json data"),
			http.StatusBadRequest, "Bad JSON data.")
		return
	} else if message, err = data.validateFields(false); err != nil {
		svr.sendErrorResponse(w, err, http.StatusBadRequest, message)
		return
	}

	hashedPass, err := app.NewPassword(data.NewPassword)
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to hash password"),
			http.StatusInternalServerError, "")
		return
	}

	err = svr.db.UpdatePersonPassword(r.Context(), userID, hashedPass)
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "failed to update password"),
			http.StatusInternalServerError, "")
		return
	}

	s := getSessionFromContext(r.Context())
	err = svr.db.RevokeSessionsForPersonExcept(r.Context(), userID, s.ID)
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "failed to revoke sessions"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (svr *Server) handleAdminDeactivateUser(
	w http.ResponseWriter,
	r *http.Request,
) {

	pathParams := mux.Vars(r)

	userID, err := strconv.Atoi(pathParams["userID"])
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "userID must be an integer"),
			http.StatusBadRequest, "User ID must be an integer.")
		return
	}

	err = svr.db.DeactivatePerson(r.Context(), userID)
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "failed to deactivate user"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (svr *Server) handleAdminActivateUser(
	w http.ResponseWriter,
	r *http.Request,
) {

	pathParams := mux.Vars(r)

	userID, err := strconv.Atoi(pathParams["userID"])
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "userID must be an integer"),
			http.StatusBadRequest, "User ID must be an integer.")
		return
	}

	err = svr.db.ActivatePerson(r.Context(), userID)
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "failed to activate user"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
