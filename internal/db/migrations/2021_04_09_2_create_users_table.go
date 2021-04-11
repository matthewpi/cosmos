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
	addMigration(&M202104092CreateUsersTable{})
}

type M202104092CreateUsersTable struct{}

var _ db.Migration = (*M202104092CreateUsersTable)(nil)

func (m *M202104092CreateUsersTable) Up(d db.DB) error {
	return d.Create("users", func(t db.Table) {
		t.BigSerial("id").
			Primary().
			Unique()
		t.VarChar("email", 255).
			Unique()
		t.VarChar("password", 255)
		t.BigInt("role_id").
			Nullable().
			References("roles", "id").
			OnDelete(db.SetDefault).
			OnUpdate(db.Cascade)
		t.TimestampTZ("created_at").
			Default("now()")
		t.TimestampTZ("updated_at").
			Default("now()")
	})
}

func (m *M202104092CreateUsersTable) Down(d db.DB) error {
	return d.DropIfExists("users")
}
