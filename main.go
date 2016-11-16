/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nats-io/nats"
)

var n *nats.Conn
var secret string

func main() {
	log.Println("starting gateway")
	setup()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/auth", authenticate)
	e.GET("/status", getStatusHandler)

	// Setup JWT auth & protected routes
	api := e.Group("/api")
	api.Use(middleware.JWT([]byte(secret)))
	setupRoutes(api)

	e.Start(":8080")
}
