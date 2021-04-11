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
	"sort"
	"strings"
)

// Table .
type Table interface {
	BigInt(name string) Column
	BigSerial(name string) Column
	Bit(name string) Column
	VarBit(name string, size uint) Column
	Bool(name string) Column
	Char(name string) Column
	VarChar(name string, size uint) Column
	Date(name string) Column
	Float8(name string) Column
	Inet(name string) Column
	Int(name string) Column
	JSON(name string) Column
	JSONB(name string) Column
	Real(name string) Column
	SmallInt(name string) Column
	SmallSerial(name string) Column
	Serial(name string) Column
	Text(name string) Column
	Time(name string) Column
	TimeTZ(name string) Column
	Timestamp(name string) Column
	TimestampTZ(name string) Column
	UUID(name string) Column

	// DropColumns .
	DropColumns(columns ...string)
	// RenameColumn .
	RenameColumn(old, new string)
}

type table struct {
	name    string
	columns map[string]*column
}

var _ Table = (*table)(nil)

func (t *table) BigInt(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  BigInt,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) BigSerial(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  BigSerial,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) Bit(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  Bit,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) VarBit(name string, size uint) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:      VarBit,
		typeSize: size,
		name:     name,
	}
	t.columns[name] = c
	return c
}

func (t *table) Bool(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  Bool,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) Char(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  Char,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) VarChar(name string, size uint) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:      VarChar,
		typeSize: size,
		name:     name,
	}
	t.columns[name] = c
	return c
}

func (t *table) Date(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  Date,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) Float8(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  Float8,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) Inet(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  Inet,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) Int(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  Int,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) JSON(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  JSON,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) JSONB(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  JSONB,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) Real(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  Real,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) SmallInt(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  SmallInt,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) SmallSerial(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  SmallSerial,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) Serial(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  Serial,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) Text(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  Text,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) Time(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  Time,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) TimeTZ(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  TimeTZ,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) Timestamp(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  Timestamp,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) TimestampTZ(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  TimestampTZ,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) UUID(name string) Column {
	c := &column{
		id:    len(t.columns),
		table: t.name,

		typ:  UUID,
		name: name,
	}
	t.columns[name] = c
	return c
}

func (t *table) DropColumns(columns ...string) {

}

func (t *table) RenameColumn(old, new string) {

}

func (t *table) build(ifNotExists bool) string {
	var b strings.Builder
	b.WriteString("CREATE TABLE ")
	if ifNotExists {
		b.WriteString("IF NOT EXISTS ")
	}
	b.WriteString(t.name)
	b.WriteString(" (\n")

	var columns []*column
	for _, c := range t.columns {
		columns = append(columns, c)
	}
	sort.Slice(columns, func(i, j int) bool {
		return columns[i].id < columns[j].id
	})

	columnCount := len(columns) - 1
	for i, c := range columns {
		b.WriteString("\t")
		c.build(&b)

		if i < columnCount {
			b.WriteString(",\n")
		} else {
			b.WriteString("\n")
		}
	}
	b.WriteString(");")
	return b.String()
}
