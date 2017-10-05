/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package envs

import (
	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Reset : Respons to POST /services/:service/reset/ and updates the
// service status to errored from in_progress
func Reset(au models.User, name string, action *models.Action) (int, []byte) {
	var e models.Env
	var b models.Build
	var builds []models.Build

	if st, res := h.IsAuthorizedToResource(&au, h.ResetBuild, e.GetType(), name); st != 200 {
		return st, res
	}

	if err := b.FindByEnvironmentName(name, &builds); err != nil {
		h.L.Warning(err.Error())
		return 500, []byte("Internal Error")
	}

	if len(builds) == 0 {
		return 404, []byte("No builds found for the specified environment")
	}

	if builds[0].Status != "in_progress" {
		return 200, []byte("Reset only applies to an 'in progress' environment, however environment '" + name + "' is on status '" + e.Status)
	}

	if err := builds[0].Reset(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	return 200, []byte("success")
}
