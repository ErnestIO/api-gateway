package envs

import (
	"net/http"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/ernestio/api-gateway/models"
	"github.com/ernestio/api-gateway/views"
)

// Get : responds to GET /services/:service with the
// details of an existing service
func Get(au models.User, name string) (int, []byte) {
	var o views.ServiceRender
	var body []byte
	var s models.Env
	var err error
	var r models.Role
	var d models.Project
	var roles []models.Role

	if s, err = s.FindLastByName(name); err != nil {
		h.L.Error(err.Error())
		return 500, []byte("Internal error")
	}

	if s.ID == "" {
		return 404, []byte("Specified environment name does not exist")
	}

	if st, res := h.IsAuthorizedToResource(&au, h.GetEnv, s.GetType(), name); st != 200 {
		return st, res
	}

	if err := r.FindAllByResource(s.GetID(), s.GetType(), &roles); err == nil {
		for _, v := range roles {
			s.Roles = append(s.Roles, v.UserID+" ("+v.Role+")")
		}
	}

	if err := d.FindByID(s.DatacenterID); err != nil {
		return 404, []byte("Project not found")
	}

	s.Project = d.Name
	s.Provider = d.Type

	if err := o.Render(s); err != nil {
		h.L.Warning(err.Error())
		return http.StatusBadRequest, []byte(err.Error())
	}
	if body, err = o.ToJSON(); err != nil {
		return 500, []byte(err.Error())
	}

	return http.StatusOK, body
}
