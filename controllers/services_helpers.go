package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
	"github.com/ghodss/yaml"
	"github.com/labstack/echo"
	"github.com/nu7hatch/gouuid"
)

// ServiceInput : service received by the endpoint
type ServiceInput struct {
	Datacenter string `json:"datacenter"`
	Name       string `json:"name"`
}

// ServicePayload : payload to be sent to workflow manager
type ServicePayload struct {
	ID         string           `json:"id"`
	PrevID     string           `json:"previous_id"`
	Datacenter *json.RawMessage `json:"datacenter"`
	Group      *json.RawMessage `json:"client"`
	Service    *json.RawMessage `json:"service"`
}

// Given an echo context, it will extract the json or yml
// request body and will processes it in order to extract
// a valid defintion
func mapInputService(c echo.Context) (s ServiceInput, definition []byte, jsonbody []byte, err error) {
	req := c.Request()
	definition, err = ioutil.ReadAll(req.Body)

	// Normalize input body to json
	ctype := req.Header.Get("Content-Type")

	if ctype != "application/json" && ctype != "application/yaml" {
		return s, definition, jsonbody, errors.New(`"Invalid input format"`)
	}

	if ctype == "application/yaml" {
		jsonbody, err = yaml.YAMLToJSON(definition)
		if err != nil {
			return s, definition, jsonbody, errors.New(`"Invalid yaml input"`)
		}
	} else {
		jsonbody = definition
	}

	if err = json.Unmarshal(jsonbody, &s); err != nil {
		return s, definition, jsonbody, errors.New(`"Invalid input"`)
	}

	return s, definition, jsonbody, nil
}

// Generates a service id composed by a random uuid, and
// a valid generated stream id
func generateServiceID(salt string) string {
	sufix := generateStreamID(salt)
	prefix, _ := uuid.NewV4()

	return prefix.String() + "-" + string(sufix[:])
}

func generateStreamID(salt string) string {
	compose := []byte(salt)
	hasher := md5.New()
	if _, err := hasher.Write(compose); err != nil {
		log.Println(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

func getDatacenter(name string, group int) (datacenter []byte, err error) {
	var d models.Datacenter
	var datacenters []models.Datacenter

	if err := d.FindByNameAndGroupID(name, group, &datacenters); err != nil {
		return datacenter, err
	}

	if len(datacenters) == 0 {
		return datacenter, errors.New(`"Specified datacenter does not exist"`)
	}

	datacenter, err = json.Marshal(datacenters[0])
	if err != nil {
		return datacenter, errors.New("Internal error trying to get the datacenter")
	}

	return datacenter, nil
}

func getGroup(id int) (group []byte, err error) {
	var g models.Group

	if err = g.FindByID(id); err != nil {
		return group, errors.New(`"Specified group does not exist"`)
	}

	if group, err = json.Marshal(g); err != nil {
		return group, errors.New(`"Internal error"`)
	}
	println(group)

	return group, nil
}

func getService(name string, group int) (service *models.Service, err error) {
	var s models.Service
	var services []models.Service

	if err = s.FindByNameAndGroupID(name, group, &services); err != nil {
		return service, h.ErrGatewayTimeout
	}

	if len(services) == 0 {
		return nil, nil
	}

	return &services[0], nil
}

func getServiceRaw(name string, group int) (service []byte, err error) {
	var s models.Service
	var services []models.Service

	if err = s.FindByNameAndGroupID(name, group, &services); err != nil {
		return nil, errors.New(`"Internal error"`)
	}

	if len(services) == 0 {
		return nil, errors.New(`"Service not found"`)
	}

	body, err := json.Marshal(services[0])
	if err != nil {
		return nil, errors.New("Internal error")
	}
	return body, nil
}

func getServicesOutput(filter map[string]interface{}) (list []views.ServiceRender, err error) {
	var s models.Service
	var services []models.Service
	var o views.ServiceRender

	if err := s.Find(filter, &services); err != nil {
		return list, err
	}

	return o.RenderCollection(services)
}
