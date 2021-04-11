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
	"fmt"
	"strconv"
	"strings"
)

// ColumnType .
type ColumnType string

func (c ColumnType) String() string {
	return string(c)
}

const (
	BigInt      ColumnType = "BIGINT"
	BigSerial   ColumnType = "BIGSERIAL"
	Bit         ColumnType = "BIT"
	VarBit      ColumnType = "VARBIT"
	Bool        ColumnType = "BOOL"
	Char        ColumnType = "CHAR"
	VarChar     ColumnType = "VARCHAR"
	Date        ColumnType = "DATE"
	Float8      ColumnType = "FLOAT8"
	Inet        ColumnType = "INET"
	Int         ColumnType = "INT"
	JSON        ColumnType = "JSON"
	JSONB       ColumnType = "JSONB"
	Real        ColumnType = "REAL"
	SmallInt    ColumnType = "SMALLINT"
	SmallSerial ColumnType = "SMALLSERIAL"
	Serial      ColumnType = "SERIAL"
	Text        ColumnType = "TEXT"
	Time        ColumnType = "TIME"
	TimeTZ      ColumnType = "TIME WITH TIME ZONE"
	Timestamp   ColumnType = "TIMESTAMP"
	TimestampTZ ColumnType = "TIMESTAMP WITH TIME ZONE"
	UUID        ColumnType = "UUID"
)

// Column .
type Column interface {
	// Change .
	Change()

	// Collation .
	Collation(string) Column

	// Default .
	Default(interface{}) Column

	// Nullable .
	Nullable() Column

	// Index .
	Index() Column
	// Primary .
	Primary() Column
	// Unique .
	Unique() Column

	// References .
	References(table, column string) Column
	// OnDelete .
	OnDelete(action ReferentialAction) Column
	// OnUpdate .
	OnUpdate(action ReferentialAction) Column
}

// column .
type column struct {
	id    int
	table string

	typ      ColumnType
	typeSize uint
	name     string
	def      string
	nullable bool

	primary   bool
	unique    bool
	reference *reference
}

var _ Column = (*column)(nil)

func (c *column) Change() {
}

func (c *column) Collation(string) Column {
	return c
}

func (c *column) Default(i interface{}) Column {
	switch i.(type) {
	case string:
		c.def = i.(string)
	case int, int8, int16, int32, int64:
		c.def = strconv.FormatInt(i.(int64), 10)
	case uint, uint8, uint16, uint32, uint64:
		c.def = strconv.FormatUint(i.(uint64), 10)
	default:
		c.def = i.(fmt.Stringer).String()
	}
	return c
}

func (c *column) Nullable() Column {
	c.nullable = true
	return c
}

func (c *column) Index() Column {
	return c
}

func (c *column) Primary() Column {
	c.primary = true
	return c
}

func (c *column) Unique() Column {
	c.unique = true
	return c
}

func (c *column) References(table, column string) Column {
	c.reference = &reference{
		table:        table,
		targetColumn: column,

		onDelete: NoAction,
		onUpdate: NoAction,
	}
	return c
}

func (c *column) OnDelete(action ReferentialAction) Column {
	if c.reference == nil {
		panic("*column#reference is nil")
	}
	c.reference.onDelete = action
	return c
}

func (c *column) OnUpdate(action ReferentialAction) Column {
	if c.reference == nil {
		panic("*column#reference is nil")
	}
	c.reference.onUpdate = action
	return c
}

func (c *column) build(b *strings.Builder) {
	b.WriteString(c.name)
	b.WriteByte(' ')
	b.WriteString(c.typ.String())
	if c.typeSize > 0 {
		b.WriteByte('(')
		b.WriteString(strconv.FormatUint(uint64(c.typeSize), 10))
		b.WriteByte(')')
	}
	if c.def != "" {
		b.WriteByte(' ')
		b.WriteString("DEFAULT ")
		b.WriteString(c.def)
	}
	b.WriteByte(' ')
	if c.nullable {
		b.WriteString("NULL")
	} else {
		b.WriteString("NOT NULL")
	}
	if c.primary {
		b.WriteByte(' ')
		(&constraint{
			typ:  PrimaryKeyConstraint,
			name: c.table + "_pk",
		}).build(b)
	}
	if c.unique {
		b.WriteByte(' ')
		(&constraint{
			typ:  UniqueConstraint,
			name: c.table + "_" + c.name + "_uindex",
		}).build(b)
	}
	if c.reference != nil {
		b.WriteString(" CONSTRAINT ")
		b.WriteString(c.table + "_" + c.reference.table + "_" + c.reference.targetColumn + "_fk")
		c.reference.build(b)
	}
}
