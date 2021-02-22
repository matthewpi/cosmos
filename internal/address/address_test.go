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

package address_test

import (
	"strings"
	"testing"

	"github.com/matthewpi/cosmos/internal/address"
)

func TestParse(t *testing.T) {
	for i, tc := range []struct {
		address   string
		mailbox   string
		domain    string
		expectErr error
	}{
		{
			address:   "mailbox@domain.com",
			mailbox:   "mailbox",
			domain:    "domain.com",
			expectErr: nil,
		},
		{
			address:   "somebody@gmail.com",
			mailbox:   "somebody",
			domain:    "gmail.com",
			expectErr: nil,
		},
		{
			address:   "user@protonmail.com",
			mailbox:   "user",
			domain:    "protonmail.com",
			expectErr: nil,
		},
		{
			address:   "ayyyyyyyyyyyyyyyyyyyyyyyyyyyyyyycomeondude@yahoo.com",
			mailbox:   "ayyyyyyyyyyyyyyyyyyyyyyyyyyyyyyycomeondude",
			domain:    "yahoo.com",
			expectErr: nil,
		},
		{
			address:   strings.Repeat("ab", 160) + "a",
			mailbox:   "",
			domain:    "",
			expectErr: address.ErrExceedsMaximumLength,
		},
		{
			address:   "does_not_contain_an_at_sign_haha",
			mailbox:   "",
			domain:    "",
			expectErr: address.ErrMalformedAddress,
		},
		{
			address:   "domain.com",
			mailbox:   "",
			domain:    "",
			expectErr: address.ErrMalformedAddress,
		},
		{
			address:   "@domain.com",
			mailbox:   "",
			domain:    "",
			expectErr: address.ErrInvalidMailbox,
		},
		{
			address:   strings.Repeat("ab", 32) + "a@domain.com",
			mailbox:   "",
			domain:    "",
			expectErr: address.ErrInvalidMailbox,
		},
		{
			address:   "test@whatanamazingdomainthatisreallylongthatnobodyshouldeverbuythanksforreadingthismessofcode.com",
			mailbox:   "",
			domain:    "",
			expectErr: address.ErrInvalidDomain,
		},
		{
			address:   "does_contain_an_at_sign@",
			mailbox:   "",
			domain:    "",
			expectErr: address.ErrInvalidDomain,
		},
		{
			address:   "a@d",
			mailbox:   "",
			domain:    "",
			expectErr: address.ErrInvalidDomain,
		},
	} {
		mailbox, domain, err := address.Parse(tc.address)

		if tc.expectErr != nil && err == nil {
			t.Errorf("Test #%d: Expected error return value, but got \"%v\"", i, err)
			continue
		}
		if tc.expectErr == nil && err != nil {
			t.Errorf("Test #%d: Should not have error return value, but received \"%v\"", i, err)
			continue
		}
		if tc.mailbox != mailbox {
			t.Errorf("Test #%d: Expected \"%s\", but got \"%s\"", i, tc.mailbox, mailbox)
			continue
		}
		if tc.domain != domain {
			t.Errorf("Test #%d: Expected \"%s\", but got \"%s\"", i, tc.mailbox, mailbox)
			continue
		}
	}
}

func TestIsValid(t *testing.T) {
	for i, tc := range []struct {
		address string
		expect  bool
	}{
		{
			address: "mailbox@domain.com",
			expect:  true,
		},
		{
			address: "somebody@gmail.com",
			expect:  true,
		},
		{
			address: "user@protonmail.com",
			expect:  true,
		},
		{
			address: "ayyyyyyyyyyyyyyyyyyyyyyyyyyyyyyycomeondude@yahoo.com",
			expect:  true,
		},
		{
			address: strings.Repeat("ab", 160) + "a",
			expect:  false,
		},
		{
			address: "does_not_contain_an_at_sign_haha",
			expect:  false,
		},
		{
			address: "domain.com",
			expect:  false,
		},
		{
			address: "@domain.com",
			expect:  false,
		},
		{
			address: strings.Repeat("ab", 32) + "a@domain.com",
			expect:  false,
		},
		{
			address: "test@whatanamazingdomainthatisreallylongthatnobodyshouldeverbuythanksforreadingthismessofcode.com",
			expect:  false,
		},
		{
			address: "does_contain_an_at_sign@",
			expect:  false,
		},
		{
			address: "a@d",
			expect:  false,
		},
	} {
		result := address.IsValid(tc.address)

		if tc.expect != result {
			t.Errorf("Test #%d: Expected \"%t\", but got \"%t\"", i, tc.expect, result)
			continue
		}
	}
}

func TestIsMailboxValid(t *testing.T) {
	for i, tc := range []struct {
		mailbox string
		expect  bool
	}{
		{
			mailbox: "a",
			expect:  true,
		},
		{
			mailbox: "abcdefghijklmnopqrstuvwxyz",
			expect:  true,
		},
		{
			mailbox: "matthew",
			expect:  true,
		},
		{
			mailbox: "",
			expect:  false,
		},
		{
			mailbox: strings.Repeat("ab", 32) + "a",
			expect:  false,
		},
	} {
		result := address.IsMailboxValid(tc.mailbox)

		if tc.expect != result {
			t.Errorf("Test #%d: Expected \"%t\", but got \"%t\"", i, tc.expect, result)
			continue
		}
	}
}

func TestIsDomainValid(t *testing.T) {
	for i, tc := range []struct {
		domain string
		expect bool
	}{
		{
			domain: "gmail.com",
			expect: true,
		},
		{
			domain: "hotmail.com",
			expect: true,
		},
		{
			domain: "outlook.com",
			expect: true,
		},
		{
			domain: "pm.me",
			expect: true,
		},
		{
			domain: "protonmail.com",
			expect: true,
		},
		{
			domain: "protonmail.ch",
			expect: true,
		},
		{
			domain: "yahoo.com",
			expect: true,
		},
		{
			domain: "sub.domain.com",
			expect: true,
		},
		{
			domain: "sub.domain.co.uk",
			expect: true,
		},
		{
			domain: "",
			expect: false,
		},
		{
			domain: "domain..com",
			expect: false,
		},
		{
			domain: "domain-whatever--something.com",
			expect: false,
		},
		{
			domain: strings.Repeat("ab", 32) + "a.domain.com",
			expect: false,
		},
		{
			domain: "whatanamazingdomainthatisreallylongthatnobodyshouldeverbuythanksforreadingthismessofcode.com",
			expect: false,
		},
	} {
		result := address.IsDomainValid(tc.domain)

		if tc.expect != result {
			t.Errorf("Test #%d: Expected \"%t\", but got \"%t\"", i, tc.expect, result)
			continue
		}
	}
}
