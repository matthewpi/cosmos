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

// ConstraintType .
type ConstraintType string

const (
	CheckConstraint      ConstraintType = "CHECK"
	PrimaryKeyConstraint ConstraintType = "PRIMARY KEY"
	UniqueConstraint     ConstraintType = "UNIQUE"
)

// Constraint .
type Constraint interface{}

// constraint .
type constraint struct {
	typ  ConstraintType
	name string
}

var _ Constraint = (*constraint)(nil)

func (c *constraint) build(b *strings.Builder) {
	b.WriteString("CONSTRAINT ")
	b.WriteString(c.name)
	b.WriteByte(' ')
	switch c.typ {
	case CheckConstraint:
	case PrimaryKeyConstraint:
		b.WriteString("PRIMARY KEY")
	case UniqueConstraint:
		b.WriteString("UNIQUE")
	}
}
