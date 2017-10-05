/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package envs

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// ForceDeletion : Deletes a service by name forcing it
func ForceDeletion(au models.User, name string) (int, []byte) {
	var e models.Env

	if st, res := h.IsAuthorizedToResource(&au, h.DeleteEnvForce, e.GetType(), name); st != 200 {
		return st, res
	}

	if err := e.DeleteByName(name); err != nil {
		h.L.Error(err.Error())
		return 500, []byte(err.Error())
	}

	return http.StatusOK, []byte(`{"status":"ok"}`)
}
