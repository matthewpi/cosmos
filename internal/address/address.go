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

// Package address ...
package address

import (
	"errors"
	"strings"
)

const (
	// Max Address Length (RFC 3696)
	maxLength = 320

	// Max Mailbox Length (MAILBOX@domain.com)
	maxMailboxLength = 64

	// Max Domain Length (mailbox@DOMAIN.COM)
	maxDomainLength = 255

	// Max Domain Section Length (SOMETHING.domain.com)
	maxDomainSectionLength = 64
)

var (
	// ErrExceedsMaximumLength is an error for an address that exceeds the maximum length.
	ErrExceedsMaximumLength = errors.New("address: exceeds the maximum length")

	// ErrMalformedAddress is an error for a malformed address.
	ErrMalformedAddress = errors.New("address: malformed address")

	// ErrInvalidMailbox is an error for an invalid mailbox.
	ErrInvalidMailbox = errors.New("address: invalid mailbox on left-side of address")

	// ErrInvalidDomain is an error for an invalid domain.
	ErrInvalidDomain = errors.New("address: invalid domain on right-side of address")
)

// Parse parses an email address into it's sections.
func Parse(address string) (string, string, error) {
	// Check if the address is longer than the max length.
	if len(address) > maxLength {
		return "", "", ErrExceedsMaximumLength
	}

	// Split the address at the @ sign.
	sections := strings.Split(address, "@")
	if len(sections) != 2 {
		return "", "", ErrMalformedAddress
	}

	// Check if the mailbox is valid.
	if !IsMailboxValid(sections[0]) {
		return "", "", ErrInvalidMailbox
	}

	// Check if the domain is valid.
	if !IsDomainValid(sections[1]) {
		return "", "", ErrInvalidDomain
	}

	return sections[0], sections[1], nil
}

// IsValid checks if an email address is valid.
func IsValid(address string) bool {
	// Check if the address is longer than the max address length.
	if len(address) > maxLength {
		return false
	}

	// It would be a good idea to check if the address contains an @ sign before splitting it,
	// however the split method will return a slice with the length of 1 that only contains
	// the "address" string if it does not contain the @ sign.

	// Split the address at the @ sign.
	sections := strings.Split(address, "@")
	if len(sections) != 2 {
		return false
	}

	return IsMailboxValid(sections[0]) && IsDomainValid(sections[1])
}

// IsMailboxValid checks if a mailbox name is valid. (left side of the @)
func IsMailboxValid(mailbox string) bool {
	// According to the RFC spec, email addresses can be almost anything a user wants...
	//
	// Because of this, there is only a max length limit enforced, most special characters
	// are allowed. This method is designed for both incoming and outgoing validation,
	// meaning, there should not be a strict ruleset for validation, preventing false-positives
	// and improving execution times compared to using a complex and unnecessary regex.
	//
	// Honestly, the only true test to if an email address is valid is to try sending an email..

	if len(mailbox) < 1 || len(mailbox) > maxMailboxLength {
		return false
	}

	return true
}

// IsDomainValid checks if a domain is valid. (right side of the @)
func IsDomainValid(domain string) bool {
	// Check if the domain is longer than the max domain length.
	if len(domain) > maxDomainLength {
		return false
	}

	// Check if the domain starts or ends with a period.
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return false
	}

	// Check if the domain contains two periods beside each other.
	if strings.Contains(domain, "..") {
		return false
	}

	// Check if the domain contains two dashes beside each other.
	if strings.Contains(domain, "--") {
		return false
	}

	// Make sure all sections of the domains are not longer than the max section length
	// (sections refer the domain split where any periods are)
	labels := strings.Split(domain, ".")
	for _, label := range labels {
		if len(label) < 1 || len(label) > maxDomainSectionLength {
			return false
		}
	}

	// Only return true if the labels slice has more than 1 item (if domain contains at least 1 dot)
	return len(labels) > 1
}
