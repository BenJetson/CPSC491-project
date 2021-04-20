package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/BenJetson/CPSC491-project/go/app"
)

type registrationRequest struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	ShouldNotify bool   `json:"should_notify"`
}

func (reg *registrationRequest) validateFields() (message string, err error) {
	if !validateEmail.MatchString(reg.Email) {
		err = errors.New("invalid email address")
		message = "Invalid email address."
		return
	}

	if len(reg.FirstName) < 1 {
		err = errors.New("first name cannot be blank")
		message = "First Name cannot be blank."
		return
	}

	if len(reg.LastName) < 1 {
		err = errors.New("last name cannot be blank")
		message = "Last Name cannot be blank."
		return
	}

	const minPassLength = 8
	if len(reg.Password) < minPassLength {
		err = errors.New("password does not meet length requirement")
		message = fmt.Sprintf("Password must be at least %d characters long.",
			minPassLength)
		return
	}

	return
}

func (svr *Server) handleRegistration(w http.ResponseWriter, r *http.Request) {
	s := getSessionFromContext(r.Context())
	if s != nil {
		svr.sendErrorResponse(w,
			errors.New("cannot register while logged in"),
			http.StatusUnauthorized, "Cannot register while logged in.")
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	var reg registrationRequest
	if err := d.Decode(&reg); err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to decode registration request"),
			http.StatusBadRequest, "Invalid registration format.")
		return
	} else if message, err := reg.validateFields(); err != nil {
		svr.sendErrorResponse(w, err, http.StatusBadRequest, message)
		return
	}

	hashedPass, err := app.NewPassword(reg.Password)
	if err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to hash password"),
			http.StatusInternalServerError, "")
		return
	}

	p := app.Person{
		FirstName: reg.FirstName,
		LastName:  reg.LastName,
		Email:     reg.Email,
		Password:  hashedPass,
		Role:      app.RoleDriver,
	}

	if _, err := svr.db.CreatePerson(r.Context(), p); err != nil {
		svr.sendErrorResponse(w,
			errors.Wrap(err, "failed to create registered person"),
			http.StatusInternalServerError, "")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
