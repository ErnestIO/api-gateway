/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"net/http"
	"strings"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
)

// Get : responds to GET /builds/:build with the
// details of an existing build
func Get(au models.User, id string) (int, []byte) {
	var o views.BuildRender
	var err error
	var body []byte
	var e models.Env
	var b models.Build
	var p models.Project
	var r models.Role
	var roles []models.Role

	if err = b.FindByID(id); err != nil {
		h.L.Error(err.Error())
		return 404, h.ErrMessage("Specified environment build does not exist")
	}

	if err := e.FindByID(int(b.EnvironmentID)); err != nil {
		return 404, h.ErrMessage("Environment not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.GetEnv, e.GetType(), e.Name); st != 200 {
		return st, res
	}

	if err := p.FindByID(int(e.ProjectID)); err != nil {
		return 404, h.ErrMessage("Project not found")
	}

	if err := r.FindAllByResource(e.GetID(), e.GetType(), &roles); err == nil {
		for _, v := range roles {
			o.Roles = append(o.Roles, v.UserID+" ("+v.Role+")")
		}
	}

	o.Name = strings.Split(e.Name, "/")[1]
	o.Project = p.Name
	o.Provider = e.Type

	if err := o.Render(b); err != nil {
		h.L.Warning(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}
	if body, err = o.ToJSON(); err != nil {
		return 500, []byte(err.Error())
	}

	return http.StatusOK, body
}
