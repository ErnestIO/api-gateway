package roles

import (
	"encoding/json"
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
)

// Create : responds to POST /roles/ by creating a
// role on the data store
func Create(au models.User, body []byte) (int, []byte) {
	var err error
	var d models.Role

	if d.Map(body) != nil {
		return 400, []byte("Input is not valid")
	}

	err = d.Validate()
	if err != nil {
		h.L.Error(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}

	if !d.ResourceExists() {
		return 404, []byte("Specified resource not found")
	}

	if !d.UserExists() {
		return 404, []byte("Specified user not found")
	}

	if !au.IsAdmin() {
		if ok := au.IsOwner(d.ResourceType, d.ResourceID); !ok {
			return 403, []byte("You're not authorized to perform this action")
		}
	}

	existing, err := d.Get(d.UserID, d.ResourceID, d.ResourceType)
	if err == nil && existing != nil {
		d.ID = existing.ID
	}

	if err = d.Save(); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	if body, err = json.Marshal(d); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal server error")
	}

	return http.StatusOK, body
}
