/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builds

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Import : Imports an environment
func Import(au models.User, env string, action *models.Action) (int, []byte) {
	var e models.Env
	var m models.Mapping

	if !models.IsAlphaNumeric(env) {
		return 404, models.NewJSONError("Environment name contains invalid characters")
	}

	err := e.FindByName(env)
	if err != nil {
		h.L.Error(err.Error())
		return 404, models.NewJSONError("Environment not found")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.UpdateEnv, e.GetType(), e.Name); st != 200 {
		return st, res
	}

	err = m.Import(env, action.Options.Filters, au)
	if err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError(`"Couldn't map the import build"`)
	}

	b := models.Build{
		ID:            m["id"].(string),
		EnvironmentID: e.ID,
		UserID:        au.ID,
		Username:      au.Username,
		Type:          "import",
		Mapping:       m,
	}

	err = b.Save()
	if err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError(`"Couldn't create the build"`)
	}

	if err := b.RequestImport(&m); err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError(`"Couldn't call build.import"`)
	}

	action.ResourceID = b.ID
	action.ResourceType = "build"
	action.Status = "in_progress"

	data, err := json.Marshal(action)
	if err != nil {
		h.L.Error(err.Error())
		return 500, models.NewJSONError(`"Couldn't marshal response"`)
	}

	return http.StatusOK, data
}
