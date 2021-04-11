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

// Package db ...
package db

import (
	"fmt"
)

// DB .
type DB interface {
	// Create .
	Create(name string, f func(Table)) error

	// CreateIfExists .
	CreateIfExists(name string, f func(Table)) error

	// Drop .
	Drop(name string) error

	// DropIfExists .
	DropIfExists(name string) error

	// Table .
	Table(name string, f func(Table)) error
}

// database .
type database struct{}

var _ DB = (*database)(nil)

// New .
func New() DB {
	return &database{}
}

func (db *database) Create(name string, f func(Table)) error {
	return db.create(name, f, false)
}

func (db *database) CreateIfExists(name string, f func(Table)) error {
	return db.create(name, f, true)
}

func (db *database) create(name string, f func(Table), ifNotExists bool) error {
	t := &table{
		name:    name,
		columns: make(map[string]*column),
	}
	f(t)

	fmt.Println(t.build(ifNotExists) + "\n")
	return nil
}

func (db *database) Drop(name string) error {
	// "DROP TABLE " + name + ";"
	return nil
}

func (db *database) DropIfExists(name string) error {
	// "DROP TABLE IF EXISTS " + name + ";"
	return nil
}

func (db *database) Table(name string, f func(Table)) error {
	return nil
}
