/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package loggers

import (
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// Delete : responds to DELETE /loggers/:id: by deleting an
// existing logger
func Delete(au models.User, logger string) (int, []byte) {
	var l models.Logger

	if logger == "basic" {
		return 400, []byte("Basic logger can't be deleted")
	}

	if err := l.Delete(logger); err != nil {
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, []byte("Logger successfully deleted")
}
