/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package envs

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Update : Not implemented
func Update(au models.User, name string, body []byte) (int, []byte) {
	var err error
	var resp []byte
	var e models.Env
	var input models.Env

	if st, res := h.IsAuthorizedToResource(&au, h.UpdateEnv, input.GetType(), name); st != 200 {
		return st, res
	}

	if err := json.Unmarshal(body, &input); err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	// Get existing environment
	if err := e.FindByName(name); err != nil {
		return 404, []byte(err.Error())
	}

	e.Options = input.Options
	e.Credentials = input.Credentials

	if err := e.Save(); err != nil {
		return 500, []byte(err.Error())
	}

	resp, err = json.Marshal(e)
	if err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	return http.StatusOK, resp
}
