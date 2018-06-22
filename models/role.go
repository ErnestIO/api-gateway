/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"encoding/json"
	"errors"
	"strconv"

	h "github.com/ernestio/api-gateway/helpers"
	"github.com/sirupsen/logrus"
)

// Role holds the role response from role
type Role struct {
	ID           uint   `json:"id"`
	UserID       string `json:"user_id"`
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	Role         string `json:"role"`
}

// Validate : validates the role
func (l *Role) Validate() error {
	if l.UserID == "" {
		return errors.New("User is empty")
	}

	if !IsAlphaNumeric(l.UserID) {
		return errors.New("User ID contains invalid characters")
	}

	if l.ResourceID == "" {
		return errors.New("Resource is empty")
	}

	if !IsAlphaNumeric(l.ResourceID) {
		return errors.New("Resource ID contains invalid characters")
	}

	if l.ResourceType != "project" && l.ResourceType != "environment" && l.ResourceType != "policy" {
		return errors.New("Resource type accepted values are ['project', 'environment', 'policy']")
	}

	if l.Role == "" {
		return errors.New("Role is empty")
	}

	return nil
}

// Map : maps a role from a request's body and validates the input
func (l *Role) Map(data []byte) error {
	if err := json.Unmarshal(data, &l); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	return nil
}

// FindAll : Searches for all roles on the system
func (l *Role) FindAll(roles *[]Role) (err error) {
	query := make(map[string]interface{})
	return NewBaseModel("authorization").FindBy(query, roles)
}

// FindByID : Gets a role by ID
func (l *Role) FindByID(id string, role *Role) (err error) {
	query := make(map[string]interface{})
	if query["id"], err = strconv.Atoi(id); err != nil {
		return err
	}
	return NewBaseModel("authorization").GetBy(query, role)
}

// FindAllByUser : Searches for all roles on the system by user
func (l *Role) FindAllByUser(u string, roles *[]Role) (err error) {
	query := make(map[string]interface{})
	query["user_id"] = u

	return NewBaseModel("authorization").FindBy(query, roles)
}

// FindAllByUserAndResource : Searches for all roles on the system by user and resource type
func (l *Role) FindAllByUserAndResource(u, r string, roles *[]Role) (err error) {
	query := make(map[string]interface{})
	query["user_id"] = u
	query["resource_type"] = r

	return NewBaseModel("authorization").FindBy(query, roles)
}

// FindAllIDsByUserAndType : Searches for all resource_ids by user and resource type
func (l *Role) FindAllIDsByUserAndType(u, r string) (ids []string, err error) {
	var rs []Role

	if err = l.FindAllByUserAndResource(u, r, &rs); err != nil {
		return
	}

	for _, r := range rs {
		ids = append(ids, r.ResourceID)
	}

	return
}

// FindAllByResource : Searches for all roles on the system by user and resource type
func (l *Role) FindAllByResource(id, r string, roles *[]Role) (err error) {
	query := make(map[string]interface{})
	query["resource_id"] = id
	query["resource_type"] = r

	return NewBaseModel("authorization").FindBy(query, roles)
}

// Save : calls role.set with the marshalled current role
func (l *Role) Save() (err error) {
	return NewBaseModel("authorization").Save(l)
}

// Get : will delete a role by its type
func (l *Role) Get(userID, resourceID, resourceType string) (role *Role, err error) {
	var roles []Role
	query := make(map[string]interface{})
	query["resource_id"] = resourceID
	query["resource_type"] = resourceType
	query["user_id"] = userID
	if err = NewBaseModel("authorization").FindBy(query, &roles); err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, nil
	}
	return &roles[0], nil
}

// Delete : will delete a role by its type
func (l *Role) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = l.ID
	query["resource_id"] = l.ResourceID
	query["resource_type"] = l.ResourceType
	query["user_id"] = l.UserID
	query["role"] = l.Role

	return NewBaseModel("authorization").Delete(query)
}

// ResourceExists : check if related resource exists
func (l *Role) ResourceExists() bool {
	if l.ResourceType == "project" {
		var r Project
		err := r.FindByName(l.ResourceID)
		if err == nil && &r != nil {
			return true
		}
	} else if l.ResourceType == "environment" {
		var r Env
		err := r.FindByName(l.ResourceID)
		if err == nil && &r != nil {
			return true
		}
	} else if l.ResourceType == "policy" {
		var r Policy
		err := r.GetByName(l.ResourceID, &r)
		if err == nil && &r != nil {
			return true
		}
	}
	return false
}

// UserExists : check if related user exists
func (l *Role) UserExists() bool {
	var r User
	err := r.FindByUserName(l.UserID, &r)
	if err == nil && &r != nil {
		return true
	}
	return false
}
