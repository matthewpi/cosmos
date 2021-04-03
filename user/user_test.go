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

package user_test

import (
	"testing"

	"github.com/matthewpi/cosmos/user"
)

func Test_New(t *testing.T) {
	for i, tc := range []struct {
		email          string
		password       []byte
		expectPassword bool
		expectErr      error
	}{
		{
			email:          "matthew@example.com",
			password:       []byte("what a really secure password"),
			expectPassword: true,
			expectErr:      nil,
		},
		{
			email:          "matthew@example.com",
			password:       nil,
			expectPassword: false,
			expectErr:      nil,
		},
	} {
		u, err := user.New(tc.email, tc.password)

		if tc.expectErr != nil && err == nil {
			t.Errorf("Test #%d: Expected error return value, but got \"%v\"", i, err)
			continue
		}
		if tc.expectErr == nil && err != nil {
			t.Errorf("Test #%d: Should not have error return value, but received \"%v\"", i, err)
			continue
		}

		// If we expected an error and received an error, skip the rest of the test cases.
		if tc.expectErr != nil && err != nil {
			continue
		}

		// This case just ensures that `u` cannot be nil, even though it should never be.
		if u == nil {
			t.Errorf("Test #%d: Expected user, received nil", i)
			continue
		}
		if tc.expectPassword && !u.HasPassword() {
			t.Errorf("Test #%d: Expected user to have password, none set", i)
			continue
		}
		if !tc.expectPassword && u.HasPassword() {
			t.Errorf("Test #%d: Expected user to not have password, one was set", i)
			continue
		}
	}
}

func TestUser_HasPassword(t *testing.T) {
	// TODO
}

func TestUser_SetPassword(t *testing.T) {
	// TODO
}

func TestUser_VerifyPassword(t *testing.T) {
	// TODO
}
