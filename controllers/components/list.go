package components

import (
	"encoding/json"
	"net/http"

	"github.com/ernestio/api-gateway/models"
)

// List : ...
func List(datacenter, service, component string) (int, []byte) {
	var body []byte
	var d models.Datacenter

	if err := d.FindByName(datacenter, &d); err != nil {
		return 404, []byte("Datacenter not found")
	}

	tags := make(map[string]string)
	if service != "" {
		tags["ernest.service"] = service
	}
	aws := models.AWSComponent{
		Datacenter: &d,
		Name:       component,
		Tags:       tags,
	}
	list, err := aws.FindBy()
	if err != nil {
		return 500, []byte("An internal error occured")
	}

	if body, err = json.Marshal(list); err != nil {
		return 500, []byte("Oops, somethign went wrong")
	}

	return http.StatusOK, body
}
