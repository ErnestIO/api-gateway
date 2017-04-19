/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package helpers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/nats-io/nats"
)

// ResponseError is
type ResponseError struct {
	Error     string          `json:"_error"`
	Code      string          `json:"_code"`
	HTTPError *echo.HTTPError `json:"-"`
}

// ResponseErr : ..
func ResponseErr(msg *nats.Msg) *ResponseError {
	var e ResponseError

	err := json.Unmarshal(msg.Data, &e)
	if err != nil || e.Error == "" {
		return nil
	}

	e.HTTPError = echo.NewHTTPError(http.StatusInternalServerError, e.Error)

	if strings.Contains(e.Error, "Not found") {
		e.HTTPError = ErrNotFound
	}

	return &e
}
