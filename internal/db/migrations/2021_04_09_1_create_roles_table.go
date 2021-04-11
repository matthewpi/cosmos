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

package migrations

import (
	"github.com/matthewpi/cosmos/internal/db"
)

func init() {
	addMigration(&M202104091CreateRolesTable{})
}

type M202104091CreateRolesTable struct{}

var _ db.Migration = (*M202104091CreateRolesTable)(nil)

func (m *M202104091CreateRolesTable) Up(d db.DB) error {
	return d.Create("roles", func(t db.Table) {
		t.BigSerial("id").Primary().Unique()
		t.VarChar("name", 32).Unique()
		t.Text("description")
		t.JSON("permissions")
		t.Int("sort_id")
	})
}

func (m *M202104091CreateRolesTable) Down(d db.DB) error {
	return d.DropIfExists("roles")
}
