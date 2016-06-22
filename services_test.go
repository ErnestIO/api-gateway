/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServices(t *testing.T) {
	os.Setenv("JWT_SECRET", "test")
	setup()

	Convey("Scenario: getting a list of services", t, func() {
		Convey("Given services exist on the store", func() {
			findServiceSubcriber()
			Convey("When I call GET /services/", func() {
				resp, err := doRequest("GET", "/services/", nil, nil, getServicesHandler, nil)

				Convey("It should return the correct set of data", func() {
					var s []Service
					So(err, ShouldBeNil)
					err = json.Unmarshal(resp, &s)
					So(err, ShouldBeNil)
					So(len(s), ShouldEqual, 2)
					So(s[0].ID, ShouldEqual, "1")
					So(s[0].Name, ShouldEqual, "test")
					So(s[0].GroupID, ShouldEqual, 1)
				})
			})
		})
	})

	Convey("Scenario: getting a single services", t, func() {
		Convey("Given the service exists on the store", func() {
			getServiceSubcriber()
			Convey("And I call /service/:service on the api", func() {
				params := make(map[string]string)
				params["service"] = "1"
				resp, err := doRequest("GET", "/services/:service", params, nil, getServiceHandler, nil)

				Convey("When I'm authenticated as an admin user", func() {
					Convey("Then I should get the existing service", func() {
						var d Service

						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &d)

						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, "1")
						So(d.Name, ShouldEqual, "test")
					})
				})

				Convey("When the service group matches the authenticated users group", func() {
					ft := generateTestToken(1, "test", false)

					params := make(map[string]string)
					params["service"] = "1"
					resp, err := doRequest("GET", "/services/:service", params, nil, getServiceHandler, ft)

					Convey("Then I should get the existing service", func() {
						var d Service
						So(err, ShouldBeNil)
						err = json.Unmarshal(resp, &d)
						So(err, ShouldBeNil)
						So(d.ID, ShouldEqual, "1")
						So(d.Name, ShouldEqual, "test")
					})
				})

				Convey("When the service group does not match the authenticated users group", func() {
					ft := generateTestToken(2, "test2", false)

					params := make(map[string]string)
					params["service"] = "1"
					_, err := doRequest("GET", "/services/:service", params, nil, getServiceHandler, ft)

					Convey("Then I should get a 404 error as it doesn't exist", func() {
						So(err, ShouldNotBeNil)
						So(err.(*echo.HTTPError).Code, ShouldEqual, 404)
					})
				})
			})
		})
	})

	Convey("Scenario: creating a service", t, func() {
		params := make(map[string]string)
		Convey("Given I do a post call to /services ", func() {
			createServiceSubcriber()

			Convey("And the content type is non json and non yaml", func() {
				data := []byte("bla")
				resp, err := doRequest("POST", "/services/", params, data, createServiceHandler, nil)
				Convey("Then I should get a 400 response", func() {
					So(err, ShouldEqual, nil)
					So(string(resp), ShouldEqual, `"Invalid input format"`)
				})
			})

			Convey("And the content type is a non valid yaml", func() {
				data := []byte("{asd}")
				headers := map[string]string{}
				headers["Content-Type"] = "application/yaml"
				resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
				Convey("Then I should get a 400 response", func() {
					So(err, ShouldEqual, nil)
					So(string(resp), ShouldEqual, `"Invalid input"`)
				})
			})

			Convey("And the content type is a non valid json", func() {
				data := []byte(`{"name"}`)
				headers := map[string]string{}
				headers["Content-Type"] = "application/json"
				resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
				Convey("Then I should get a 400 response", func() {
					So(err, ShouldEqual, nil)
					So(string(resp), ShouldEqual, `"Invalid input"`)
				})
			})

			Convey("And the specified datacenter does not exist", func() {
				notFoundSubscriber("datacenter.find", 1)
				data := []byte(`{"name":"test"}`)
				headers := map[string]string{}
				headers["Content-Type"] = "application/json"
				resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
				Convey("Then I should get a 404 response", func() {
					So(err, ShouldEqual, nil)
					So(string(resp), ShouldEqual, `"Specified datacenter does not exist"`)
				})
			})

			Convey("And the specified group does not exist", func() {
				notFoundSubscriber("group.get", 1)
				foundSubscriber("datacenter.find", `{"id":"1"}`, 1)
				data := []byte(`{"name":"test"}`)
				headers := map[string]string{}
				headers["Content-Type"] = "application/json"
				resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
				Convey("Then I should get a 404 response", func() {
					So(err, ShouldEqual, nil)
					So(string(resp), ShouldEqual, `"Specified group does not exist"`)
				})
			})

			Convey("And I provide a valid input", func() {
				foundSubscriber("group.get", `{"id":"1"}`, 1)
				foundSubscriber("datacenter.find", `{"id":"1"}`, 1)
				data := []byte(`{"name":"test"}`)
				headers := map[string]string{}
				headers["Content-Type"] = "application/json"

				Convey("And the service does not exist", func() {
					notFoundSubscriber("service.get", 1)
					foundSubscriber("service.create", `{"id":"1"}`, 1)
					foundSubscriber("definition.map_create", `{"id":"1"}`, 1)
					resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
					Convey("Then I should get a response with a valid id", func() {
						So(err, ShouldEqual, nil)
						So(strings.Contains(string(resp), `{"id":"`), ShouldEqual, true)
						So(strings.Contains(string(resp), `-d29d2764b65cae3f4114164bb6cf80cb`), ShouldEqual, true)
					})
				})

				Convey("And the service already exists", func() {
					foundSubscriber("service.create", `{"id":"1"}`, 1)
					Convey("And the existing service is done", func() {
						foundSubscriber("definition.map_create", `{"id":"1"}`, 1)
						foundSubscriber("service.get", `{"id":"foo-bar","status":"done"}`, 1)
						resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
						Convey("Then I should get a response with the existing id", func() {
							So(err, ShouldEqual, nil)
							So(strings.Contains(string(resp), `{"id":"`), ShouldEqual, true)
							So(strings.Contains(string(resp), `-d29d2764b65cae3f4114164bb6cf80cb`), ShouldEqual, true)
						})
					})

					Convey("And the existing service is in progress", func() {
						foundSubscriber("service.get", `{"id":"foo-bar","status":"in_progress"}`, 1)
						resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
						Convey("Then I should get an error as an in_progress service can't be modified", func() {
							So(err, ShouldEqual, nil)
							So(string(resp), ShouldEqual, `"Your service process is 'in progress' if your're sure you want to fix it please reset it first"`)
						})
					})

					Convey("And the existing service is errored", func() {
						foundSubscriber("service.get", `{"id":"foo-bar","status":"errored"}`, 1)
						foundSubscriber("definition.map_create", `{"id":"1"}`, 1)
						foundSubscriber("service.patch", `{"id":"1"}`, 1)
						resp, err := doRequestHeaders("POST", "/services/", params, data, createServiceHandler, nil, headers)
						Convey("Then I should get a response with the existing id", func() {
							So(err, ShouldEqual, nil)
							So(strings.Contains(string(resp), `{"id":"`), ShouldEqual, true)
							So(strings.Contains(string(resp), `-d29d2764b65cae3f4114164bb6cf80cb`), ShouldEqual, true)
						})
					})
				})
			})
			/*
				Convey("When I do a post to /services/ with invalid innput type", func() {
					mockDC := Service{
						GroupID: 1,
						Name:    "new-test",
						Type:    "vcloud",
					}

					data, _ := json.Marshal(mockDC)
					params := make(map[string]string)
					params["service"] = "test"

					Convey("And I am logged in as an admin", func() {
						resp, err := doRequest("POST", "/services/", params, data, createServiceHandler, nil)

						Convey("Then a service should be created", func() {
							var d Service
							So(err, ShouldBeNil)
							err = json.Unmarshal(resp, &d)
							So(err, ShouldBeNil)
							So(d.ID, ShouldEqual, 3)
							So(d.Name, ShouldEqual, "new-test")
						})
					})

					Convey("And the service group matches the authenticated users group", func() {
						ft := generateTestToken(1, "test", false)
						resp, err := doRequest("POST", "/services/", params, data, createServiceHandler, ft)

						Convey("It should create the service and return the correct set of data", func() {
							var d Service
							So(err, ShouldBeNil)
							err = json.Unmarshal(resp, &d)
							So(err, ShouldBeNil)
							So(d.ID, ShouldEqual, 3)
							So(d.Name, ShouldEqual, "new-test")
						})
					})

					Convey("And the service group does not match the authenticated users group", func() {
						ft := generateTestToken(2, "test2", false)
						_, err := doRequest("POST", "/services/", params, data, createServiceHandler, ft)

						Convey("It should return an 403 unauthorized error", func() {
							So(err, ShouldNotBeNil)
							So(err.(*echo.HTTPError).Code, ShouldEqual, 403)
						})
					})
				})
			*/
		})
	})
	SkipConvey("Scenario: deleting a service", t, func() {
		Convey("Given a service exists on the store", func() {
			deleteServiceSubcriber()

			Convey("When I call DELETE /services/:service", func() {
				ft := generateTestToken(1, "test", false)

				params := make(map[string]string)
				params["service"] = "test"
				_, err := doRequest("DELETE", "/services/:service", params, nil, deleteServiceHandler, ft)

				Convey("It should delete the service and return ok", func() {
					So(err, ShouldBeNil)
				})
			})

		})

	})
}
