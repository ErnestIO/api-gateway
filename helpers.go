/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

var (
	ErrUnauthorized   = echo.NewHTTPError(http.StatusForbidden, "")
	ErrNotFound       = echo.NewHTTPError(http.StatusNotFound, "")
	ErrBadReqBody     = echo.NewHTTPError(http.StatusBadRequest, "")
	ErrGatewayTimeout = echo.NewHTTPError(http.StatusGatewayTimeout, "")
	ErrInternal       = echo.NewHTTPError(http.StatusInternalServerError, "")
	ErrNotImplemented = echo.NewHTTPError(http.StatusNotImplemented, "")
)

// Get the authenticated user from the JWT Token
func authenticatedUser(c echo.Context) User {
	var u User

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(ErnestClaims)

	u.GroupID = claims.GroupID
	u.Username = claims.Username
	u.Admin = claims.Admin

	return u
}
