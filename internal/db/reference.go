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

package db

import (
	"strings"
)

// ReferentialAction .
type ReferentialAction string

const (
	// Cascade .
	Cascade ReferentialAction = "CASCADE"
	// NoAction .
	NoAction ReferentialAction = "NO ACTION"
	// Restrict .
	Restrict ReferentialAction = "RESTRICT"
	// SetDefault .
	SetDefault ReferentialAction = "SET DEFAULT"
	// SetNull .
	SetNull ReferentialAction = "SET NULL"
)

// Reference .
type Reference interface{}

// reference .
type reference struct {
	table        string
	targetColumn string

	onDelete ReferentialAction
	onUpdate ReferentialAction
}

var _ Reference = (*reference)(nil)

func (r *reference) build(b *strings.Builder) {
	b.WriteString(" REFERENCES ")
	b.WriteString(r.table)
	b.WriteString(" ON DELETE ")
	b.WriteString(string(r.onDelete))
	b.WriteString(" ON UPDATE ")
	b.WriteString(string(r.onUpdate))
}
