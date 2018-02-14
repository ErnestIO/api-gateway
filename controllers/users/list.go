/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package users

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// List : responds to GET /users/ with a list of authorized users
func List(au models.User) (int, []byte) {
	var users []models.User
	var user models.User

	if au.IsAdmin() {
		if err := au.FindAll(&users); err != nil {
			return 500, models.NewJSONError("Internal server error")
		}
	} else {
		if err := au.FindByUserName(au.Username, &user); err != nil {
			return 500, models.NewJSONError("Internal server error")
		}
		users = append(users, user)
	}

	for i := 0; i < len(users); i++ {
		users[i].Redact()
		users[i].Improve()
	}

	body, err := json.Marshal(users)
	if err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, body
}
