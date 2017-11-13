package users

import (
	"encoding/json"
	"errors"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Update : responds to PUT /users/:id: by updating an existing
// user
func Update(au models.User, user string, body []byte) (int, []byte) {
	var u models.User
	var existing models.User

	if err := u.Map(body); err != nil {
		h.L.Error(err.Error())
		return 400, []byte(err.Error())
	}

	// Check if authenticated user is admin or updating itself
	if !u.CanBeChangedBy(au) {
		err := errors.New("You're not allowed to perform this action, please contact your admin")
		h.L.Error(err.Error())
		return 403, []byte(err.Error())
	}

	// Check user exists
	if err := au.FindByUserName(user, &existing); err != nil {
		if err := au.FindByID(user, &existing); err != nil {
			h.L.Error(err.Error())
			return 404, []byte("Specified user not found")
		}
	}

	if existing.ID == 0 {
		err := errors.New("Specified user not found")
		h.L.Error(err.Error())
		return 404, []byte(err.Error())
	}

	u.Username = existing.Username

	if !au.IsAdmin() && existing.Username != au.Username {
		err := errors.New("You're not allowed to perform this action, please contact your admin")
		h.L.Error(err.Error())
		return 403, []byte(err.Error())
	}

	if !au.IsAdmin() && existing.IsAdmin() != u.IsAdmin() {
		err := errors.New("You're not allowed to perform this action, please contact your admin")
		h.L.Error(err.Error())
		return 403, []byte(err.Error())
	}

	if u.Password != nil {
		err := u.Validate()
		if err != nil {
			return 400, []byte(err.Error())
		}

		// Check the old password if it is present
		if u.OldPassword != nil && !existing.ValidPassword(*u.OldPassword) {
			err := errors.New("You're not allowed to perform this action, please contact your admin")
			h.L.Error(err.Error())
			return 403, []byte(err.Error())
		}
	}

	if err := u.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Error updating user")
	}

	if *u.MFA && !*existing.MFA {
		mfaSecret := u.MFASecret
		u.Redact()
		u.MFASecret = mfaSecret
	} else {
		u.Redact()
	}

	body, err := json.Marshal(u)
	if err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
