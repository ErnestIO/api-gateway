/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package notifications

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Update : ...
func Update(au models.User, id string, body []byte) (int, []byte) {
	var err error
	var d models.Notification
	var existing models.Notification

	if au.Admin == false {
		return 403, []byte("You should provide admin credentials to perform this action")
	}

	if d.Map(body) != nil {
		return http.StatusBadRequest, []byte("Invalid input")
	}

	if err = existing.FindByID(id, &existing); err != nil {
		return 500, []byte("Internal server error")
	}

	existing.Config = d.Config

	if err = existing.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	if body, err = json.Marshal(d); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
