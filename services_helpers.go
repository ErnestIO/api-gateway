package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/labstack/echo"
	"github.com/nats-io/nats"
	"github.com/nu7hatch/gouuid"
)

// ServiceInput : service received by the endpoint
type ServiceInput struct {
	Datacenter string `json:"datacenter"`
	Provider   string `json:"provider"`
	Name       string `json:"name"`
}

type ServicePayload struct {
	ID         string           `json:"id"`
	PrevID     string           `json:"previous_id"`
	Datacenter *json.RawMessage `json:"datacenter"`
	Group      *json.RawMessage `json:"client"`
	Service    *json.RawMessage `json:"service"`
}

// Maps input as a valid Serviceinput
func mapInputService(c echo.Context) (s ServiceInput, body []byte, err error) {
	req := c.Request()
	body, err = ioutil.ReadAll(req.Body())

	// Normalize input body to json
	ctype := req.Header().Get("Content-Type")

	if ctype != "application/json" && ctype != "application/yaml" {
		return s, body, errors.New(`"Invalid input format"`)
	} else if ctype == "application/yaml" {
		if body, err = yaml.JSONToYAML(body); err != nil {
			return s, body, errors.New(`"Invalid yaml input"`)
		}
	}

	if err = json.Unmarshal(body, &s); err != nil {
		return s, body, errors.New(`"Invalid input"`)
	}

	return s, body, nil
}

// Generates a service ID based on an input service
func generateServiceID(salt string) string {
	sufix := generateStreamID(salt)
	prefix, _ := uuid.NewV4()

	return prefix.String() + "-" + string(sufix[:])
}

func generateStreamID(salt string) string {
	compose := []byte(salt)
	hasher := md5.New()
	hasher.Write(compose)
	return hex.EncodeToString(hasher.Sum(nil))
}

func getDatacenter(id string, group int, provider string) (datacenter []byte, err error) {
	var msg *nats.Msg

	query := fmt.Sprintf(`{"id": %s, "group_id": %d}`, id, group)
	if msg, err = n.Request("datacenter.find", []byte(query), 1*time.Second); err != nil {
		return datacenter, ErrGatewayTimeout
	}
	if string(msg.Data) == `[]` {
		return datacenter, errors.New(`"Specified datacenter does not exist"`)
	}

	// FIXME This is just a temporal fix until we introduce typed datacenters
	if provider == "fake" {
		datacenter = []byte(`{"id":"fake","name":"fake","username":"fake","password":"fake_pwd","region":"fake","type":"fake","external_network":"fake","vse_url":"http://vse.url/","vcloud_url":"fake"}`)
	}

	return msg.Data, nil
}

func getGroup(id int) (group []byte, err error) {
	var msg *nats.Msg

	query := fmt.Sprintf(`{"id": %d}`, group)
	if msg, err = n.Request("group.get", []byte(query), 1*time.Second); err != nil {
		return group, ErrGatewayTimeout
	}
	if strings.Contains(string(msg.Data), `"error"`) {
		return group, errors.New(`"Specified group does not exist"`)
	}
	return msg.Data, nil
}

func getService(name string, group int) (service *Service, err error) {
	var msg *nats.Msg

	query := fmt.Sprintf(`{"name":"%s","group_id":%d}`, name, group)
	if msg, err = n.Request("service.find", []byte(query), 1*time.Second); err != nil {
		return service, ErrGatewayTimeout
	}

	p := []Service{}
	json.Unmarshal(msg.Data, &p)
	if len(p) == 0 {
		return nil, nil
	}

	return &p[0], nil
}

func mapCreateDefinition(payload ServicePayload) (body []byte, err error) {
	var msg *nats.Msg

	if body, err = json.Marshal(payload); err != nil {
		return body, errors.New("Provided yaml is not valid")
	}

	if msg, err = n.Request("definition.map.creation", body, 1*time.Second); err != nil {
		return body, errors.New("Provided yaml is not valid")
	}

	return msg.Data, nil
}

func getServiceRaw(name string, group int) (service []byte, err error) {
	var msg *nats.Msg

	query := fmt.Sprintf(`{"name":"%s","group_id":%d}`, name, group)
	if msg, err = n.Request("service.find", []byte(query), 1*time.Second); err != nil {
		return service, ErrGatewayTimeout
	}
	p := []*json.RawMessage{}

	if err = json.Unmarshal(msg.Data, &p); err != nil {
		return nil, errors.New(`"Internal error"`)
	}

	if len(p) == 0 {
		return nil, errors.New(`"Service not found"`)
	}

	if body, err := p[0].MarshalJSON(); err != nil {
		return nil, errors.New("Internal error")
	} else {
		return body, nil
	}
}

type OutputService struct {
	ID           string `json:"id"`
	DatacenterID int    `json:"datacenter_id"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Status       string `json:"status"`
	Options      string `json:"options"`
	Endpoint     string `json:"endpoint"`
	Definition   string `json:"definition"`
}

func getServicesOutput(filter map[string]interface{}) (list []OutputService, err error) {
	var msg *nats.Msg

	query, err := json.Marshal(filter)
	if err != nil {
		return list, err
	}

	if msg, err = n.Request("service.find", query, 1*time.Second); err != nil {
		return list, ErrGatewayTimeout
	}

	if err := json.Unmarshal(msg.Data, &list); err != nil {
		return list, errors.New("Internal error")
	}

	return list, nil
}

func resetService(au User, name string) (status int, err error) {
	var list []OutputService
	filter := make(map[string]interface{})
	filter["group_id"] = au.GroupID
	filter["name"] = name

	if list, err = getServicesOutput(filter); err != nil {
		return 500, errors.New("Internal error")
	}
	if len(list) == 0 {
		return 404, errors.New(`No services found with for '` + name + `'`)
	}
	if list[0].Status != "in_progress" {
		return 200, errors.New("Reset only applies to 'in progress' serices, however service '" + name + "' is on status '" + list[0].Status)
	}

	query := `{"id":"` + list[0].ID + `","status":"errored"}`
	if _, err := n.Request("service.update", []byte(query), 1*time.Second); err != nil {
		return 500, errors.New("Could not update the service")
	}

	return 200, nil
}
