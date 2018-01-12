/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package policies

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// List : responds to GET /policies/ with a list of all
// policies
func List(au models.User) (int, []byte) {
	var body []byte

	policies, err := au.GetPolicies()
	if err != nil {
		return 404, []byte(err.Error())
	}

	if body, err = json.Marshal(policies); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}
	return http.StatusOK, body
}
