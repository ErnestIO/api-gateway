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

// Policy holds the policy response from policy
type Policy struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

// Validate : validates the policy
func (l *Policy) Validate() error {
	if l.Definition == "" {
		return errors.New("Policy definition is empty")
	}

	return nil
}

// Map : maps a datacenter from a request's body and validates the input
func (l *Policy) Map(data []byte) error {
	if err := json.Unmarshal(data, &l); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	return nil
}

// FindAll : Searches for all policys on the system
func (l *Policy) FindAll(policys *[]Policy) (err error) {
	query := make(map[string]interface{})
	return NewBaseModel("policy").FindBy(query, policys)
}

// FindByID : Gets a policy by ID
func (l *Policy) FindByID(id string, policy *Policy) (err error) {
	query := make(map[string]interface{})
	if query["id"], err = strconv.Atoi(id); err != nil {
		return err
	}
	return NewBaseModel("policy").GetBy(query, policy)
}

// FindByName : Searches for all policys with a name equal to the specified
func (l *Policy) FindByName(name string, policy *Policy) (err error) {
	query := make(map[string]interface{})
	query["name"] = name
	return NewBaseModel("policy").GetBy(query, policy)
}

// Save : calls policy.set with the marshalled current policy
func (l *Policy) Save() (err error) {
	return NewBaseModel("policy").Save(l)
}

// Delete : will delete a policy by its type
func (l *Policy) Delete() (err error) {
	query := make(map[string]interface{})
	query["id"] = l.ID
	return NewBaseModel("policy").Delete(query)
}
