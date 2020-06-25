// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package acmeclient

import (
	"crypto"

	"github.com/go-acme/lego/v3/registration"
)

// User implements registration.User (go-acme/lego).
type User struct {
	Email        string
	Registration *registration.Resource
	Key          crypto.PrivateKey
}

// NewUser creates a new User.
func NewUser(email string, key crypto.PrivateKey) *User {
	return &User{
		Email: email,
		Key:   key,
	}
}

// GetEmail returns the email address associated with the User.
func (u *User) GetEmail() string {
	return u.Email
}

// GetRegistration returns the registration resource associated with the User.
func (u *User) GetRegistration() *registration.Resource {
	return u.Registration
}

// SetRegistration sets the registration resource associated with the User.
func (u *User) SetRegistration(r *registration.Resource) {
	u.Registration = r
}

// GetPrivateKey returns the private key associated with the User.
func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}
