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
func Update(au models.User, name string, body []byte) (int, []byte) {
	var err error
	var n models.Notification
	var existing models.Notification
	var e models.Env
	var p models.Project
	var envs []models.Env
	var projects []models.Project

	if n.Map(body) != nil {
		return http.StatusBadRequest, models.NewJSONError("Invalid input")
	}

	err = n.Validate()
	if err != nil {
		return 400, models.NewJSONError(err.Error())
	}

	err = existing.FindByName(name, &existing)
	if err != nil {
		return 404, models.NewJSONError("Not found")
	}

	err = p.FindAll(au, &projects)
	if err != nil {
		return 400, models.NewJSONError(err.Error())
	}

	err = e.FindAll(au, &envs)
	if err != nil {
		return 400, models.NewJSONError(err.Error())
	}

SOURCELOOP:
	for _, source := range n.Sources {
		for _, project := range projects {
			if project.Name == source {
				continue SOURCELOOP
			}
		}

		for _, env := range envs {
			if env.Name == source {
				continue SOURCELOOP
			}
		}

		return 400, models.NewJSONError("Notification source '" + source + "' does not exist")
	}

	existing.Config = n.Config
	existing.Sources = n.Sources

	err = existing.Save()
	if err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError("Internal server error")
	}

	body, err = json.Marshal(n)
	if err != nil {
		return 500, models.NewJSONError("Internal server error")
	}

	return http.StatusOK, body
}
