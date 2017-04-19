/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package models

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"

	"github.com/Sirupsen/logrus"
	h "github.com/ernestio/api-gateway/helpers"
	"golang.org/x/crypto/scrypt"
)

const (
	// SaltSize is the lenght of the salt string
	SaltSize = 32
	// HashSize is the lenght of the hash string
	HashSize = 64
)

// User holds the user response from user-store
type User struct {
	ID          int    `json:"id"`
	GroupID     int    `json:"group_id"`
	GroupName   string `json:"group_name"`
	Username    string `json:"username"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"oldpassword,omitempty"`
	Salt        string `json:"salt,omitempty"`
	Admin       bool   `json:"admin"`
}

// Validate vaildate all of the user's input
func (u *User) Validate() error {
	r := regexp.MustCompile("^[a-zA-Z0-9@._-]*$")

	if u.Username == "" {
		return errors.New("Username cannot be empty")
	}

	if !r.MatchString(u.Username) {
		return errors.New("Username can only contain the following characters: a-z 0-9 @._-")
	}

	if u.Password == "" {
		return errors.New("Password cannot be empty")
	}

	if !r.MatchString(u.Password) {
		return errors.New("Password can only contain the following characters: a-z 0-9 @._-")
	}

	if len(u.Password) < 8 {
		return errors.New("Minimum password length is 8 characters")
	}

	return nil
}

// Map a user from a request's body and validates the input
func (u *User) Map(data []byte) error {
	if err := json.Unmarshal(data, &u); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error("Couldn't unmarshal given input")
		return NewError(InvalidInputCode, "Invalid input")
	}

	if err := u.Validate(); err != nil {
		h.L.WithFields(logrus.Fields{
			"input": string(data),
		}).Error(err.Error())
		return NewError(InvalidInputCode, err.Error())
	}

	return nil
}

// FindByUserName : find a user for the given username, and maps it on
// the fiven User struct
func (u *User) FindByUserName(name string, user *User) (err error) {
	query := make(map[string]interface{})
	query["username"] = name
	if err := NewBaseModel("user").GetBy(query, user); err != nil {
		return err
	}
	return nil
}

// FindAll : Searches for all users on the store current user
// has access to
func (u *User) FindAll(users *[]User) (err error) {
	query := make(map[string]interface{})
	if !u.Admin {
		query["group_id"] = u.GroupID
	}
	if err := NewBaseModel("user").FindBy(query, users); err != nil {
		return err
	}
	return nil
}

// FindByID : Searches a user by ID on the store current user
// has access to
func (u *User) FindByID(id string, user *User) (err error) {
	query := make(map[string]interface{})
	if query["id"], err = strconv.Atoi(id); err != nil {
		return err
	}
	if !u.Admin {
		query["group_id"] = u.GroupID
	}
	if err := NewBaseModel("user").GetBy(query, user); err != nil {
		return err
	}
	return nil
}

// Save : calls user.set with the marshalled current user
func (u *User) Save() (err error) {
	if err := NewBaseModel("user").Save(u); err != nil {
		return err
	}
	return nil
}

// Delete : will delete a user by its id
func (u *User) Delete(id string) (err error) {
	query := make(map[string]interface{})
	if query["id"], err = strconv.Atoi(id); err != nil {
		return err
	}
	if err := NewBaseModel("user").Delete(query); err != nil {
		return err
	}
	return nil
}

// Redact : removes all sensitive fields from the return
// data before outputting to the user
func (u *User) Redact() {
	u.Password = ""
	u.Salt = ""
}

// Improve : adds extra data as group name
func (u *User) Improve() {
	g := u.Group()
	u.GroupName = g.Name
}

// ValidPassword : checks if a submitted password matches
// the users password hash
func (u *User) ValidPassword(pw string) bool {
	userpass, err := base64.StdEncoding.DecodeString(u.Password)
	if err != nil {
		return false
	}

	usersalt, err := base64.StdEncoding.DecodeString(u.Salt)
	if err != nil {
		return false
	}

	hash, err := scrypt.Key([]byte(pw), usersalt, 16384, 8, 1, HashSize)
	if err != nil {
		return false
	}

	// Compare in constant time to mitigate timing attacks
	if subtle.ConstantTimeCompare(userpass, hash) == 1 {
		return true
	}

	return false
}

// Group : Gets the related user group if any
func (u *User) Group() (group Group) {
	if err := group.FindByID(u.GroupID); err != nil {
		h.L.Warning(err.Error())
	}

	return group
}

// Datacenters : Gets the related user datacenters if any
func (u *User) Datacenters() (ds []Datacenter, err error) {
	var d Datacenter

	err = d.FindByGroupID(u.GroupID, &ds)

	return ds, err
}

// FindAllKeyValue : Finds all users on a id:name hash
func (u *User) FindAllKeyValue() (list map[int]string) {
	var users []User
	list = make(map[int]string)
	if err := u.FindAll(&users); err != nil {
		h.L.Warning(err.Error())
	}
	for _, v := range users {
		list[v.ID] = v.Username
	}
	return list
}
