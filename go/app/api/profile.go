package api

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

func (svr *Server) getMyProfileUserID(
	w http.ResponseWriter,
	r *http.Request,
) (session app.Session, userID int, ok bool) {
	s := getSessionFromContext(r.Context())
	if s == nil {
		svr.sendErrorResponse(w, errors.New("missing session for profile"),
			http.StatusInternalServerError, "")

		return
	}

	session = *s
	userID = s.Person.ID
	ok = true

	return
}

type nameChangeRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (r *nameChangeRequest) validateFields() (message string, err error) {
	defer func() {
		if message != "" {
			err = errors.New(message)
		}
	}()

	if len(r.FirstName) < 1 {
		message = "First name cannot be blank."
		return
	}

	if len(r.LastName) < 1 {
		message = "Last name cannot be blank."
		return
	}

	return
}

func (svr *Server) handleMyProfileUpdateName(
	w http.ResponseWriter,
	r *http.Request,
) {

	_, userID, ok := svr.getMyProfileUserID(w, r)
	if !ok {
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var data nameChangeRequest
	var message string
	if err := d.Decode(&data); err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "received bad json data"),
			http.StatusBadRequest, "Bad JSON data.")
		return
	} else if message, err = data.validateFields(); err != nil {
		svr.sendErrorResponse(w, err, http.StatusBadRequest, message)
		return
	}

	err := svr.db.UpdatePersonName(r.Context(),
		userID, data.FirstName, data.LastName)
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "failed to update name"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type emailChangeRequest struct {
	Email string `json:"email"`
}

func (r *emailChangeRequest) validateFields() (message string, err error) {
	defer func() {
		if message != "" {
			err = errors.New(message)
		}
	}()

	if !validateEmail.MatchString(r.Email) {
		message = "Invalid email address."
		return
	}

	return
}

func (svr *Server) handleMyProfileUpdateEmail(
	w http.ResponseWriter,
	r *http.Request,
) {

	_, userID, ok := svr.getMyProfileUserID(w, r)
	if !ok {
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var data emailChangeRequest
	var message string
	if err := d.Decode(&data); err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "received bad json data"),
			http.StatusBadRequest, "Bad JSON data.")
		return
	} else if message, err = data.validateFields(); err != nil {
		svr.sendErrorResponse(w, err, http.StatusBadRequest, message)
		return
	}

	err := svr.db.UpdatePersonEmail(r.Context(), userID, data.Email)
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "failed to update email"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type passwordChangeRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (r *passwordChangeRequest) validateFields(
	requireCurrent bool,
) (message string, err error) {

	defer func() {
		if message != "" {
			err = errors.New(message)
		}
	}()

	if requireCurrent && len(r.CurrentPassword) < 1 {
		message = "Current Password cannot be blank."
		return
	}

	if len(r.NewPassword) < 1 {
		message = "New Password cannot be blank."
		return
	}

	if message, err = validatePassword(r.NewPassword); err != nil {
		return
	}

	return
}

func (svr *Server) handleMyProfileUpdatePassword(
	w http.ResponseWriter,
	r *http.Request,
) {

	s, userID, ok := svr.getMyProfileUserID(w, r)
	if !ok {
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var data passwordChangeRequest
	var message string
	if err := d.Decode(&data); err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "received bad json data"),
			http.StatusBadRequest, "Bad JSON data.")
		return
	} else if message, err = data.validateFields(false); err != nil {
		svr.sendErrorResponse(w, err, http.StatusBadRequest, message)
		return
	}

	if !s.Person.Password.Verify(data.CurrentPassword) {
		svr.sendErrorResponse(
			w,
			errors.New("password did not match"),
			http.StatusForbidden,
			"Current password was incorrect.",
		)
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

	err = svr.db.RevokeSessionsForPersonExcept(r.Context(), userID, s.ID)
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "failed to revoke sessions"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (svr *Server) handleMyProfileDeactivate(
	w http.ResponseWriter,
	r *http.Request,
) {

	_, userID, ok := svr.getMyProfileUserID(w, r)
	if !ok {
		return
	}

	err := svr.db.DeactivatePerson(r.Context(), userID)
	if err != nil {
		svr.sendErrorResponse(w, errors.Wrap(err, "failed to deactivate user"),
			http.StatusInternalServerError, "")
		return
	}

	// Force the frontend to log the user out.
	w.WriteHeader(http.StatusUnauthorized)
}
