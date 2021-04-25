package api

import "github.com/pkg/errors"

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
