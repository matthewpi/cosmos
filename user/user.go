//
// Copyright (c) 2021 Matthew Penner
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

// Package user ...
package user

import (
	"time"

	"github.com/matthewpi/cosmos/internal/argon2"
	"github.com/matthewpi/cosmos/internal/snowflake"
)

// User represents a Cosmos User.
type User struct {
	// ID is the user's unique identifier. (unique, crypto-secure random)
	ID snowflake.Snowflake

	// Email is the user's email address. (unique, end-user data)
	Email string `json:"email,omitempty"`

	// password is an argon2 hash of the user's password. (technically end-user data)
	password string

	// Confirmed represents if the User has confirmed their email address.
	Confirmed bool `json:"confirmed"`

	// Locked represents if the User's account is locked.
	Locked bool `json:"-"`

	// Avatar is a hash of the User's avatar.
	Avatar string `json:"avatar"`

	// CreatedAt is a timestamp of when the account was created.
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// New .
func New(email string, password []byte) (*User, error) {
	now := time.Now()
	u := &User{
		ID:        snowflake.NewAtTime(now),
		Email:     email,
		Confirmed: false,
		Locked:    false,
		CreatedAt: now,
	}
	if password != nil {
		if err := u.SetPassword(password); err != nil {
			return nil, err
		}
	}
	return u, nil
}

// HasPassword .
func (u *User) HasPassword() bool {
	return u.password != ""
}

// SetPassword hashes a raw password and updates the user's password.
func (u *User) SetPassword(password []byte) error {
	h, err := argon2.Hash(password)
	if err != nil {
		return err
	}
	u.password = h
	return nil
}

// VerifyPassword takes a password and the user's hashed password and verifies the password against the hashed password.
func (u *User) VerifyPassword(password []byte) error {
	return argon2.Verify(password, u.password)
}
