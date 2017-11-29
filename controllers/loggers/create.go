/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package loggers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ernestio/api-gateway/models"
)

// Create : responds to POST /loggers/ by creating a logger
// on the data store
func Create(au models.User, body []byte) (int, []byte) {
	var l models.Logger
	var err error

	if l.Map(body) != nil {
		return 400, []byte("Invalid input")
	}

	if err = l.Save(); err != nil {
		var e struct {
			Msg []byte `json:"_error"`
		}
		parts := strings.Split(err.Error(), "message=")
		if len(parts) > 0 {
			return 500, []byte(parts[1])
		}
		return 500, []byte(e.Msg)
	}

	if body, err = json.Marshal(l); err != nil {
		return 500, []byte("Internal server error")
	}
	return http.StatusOK, body
}
