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

package uuid_test

import (
	"testing"

	"github.com/matthewpi/cosmos/internal/uuid"
)

func TestNew(t *testing.T) {
	id, err := uuid.New()
	if err != nil {
		t.Errorf("Should not have error return value, but received \"%v\"", err)
		return
	}

	if len(id) != 16 {
		t.Errorf("Expected length of 16, but got %d", len(id))
		return
	}
}

func TestUUID_Dashed(t *testing.T) {
	id, err := uuid.New()
	if err != nil {
		t.Errorf("Should not have error return value, but received \"%v\"", err)
		return
	}

	if len(id) != 16 {
		t.Errorf("Expected length of 16, but got %d", len(id))
		return
	}

	dashed := id.Dashed()
	if len(dashed) != 36 {
		t.Errorf("Expected length of 36, but got %d", len(dashed))
		return
	}
}

func TestUUID_String(t *testing.T) {
	id, err := uuid.New()
	if err != nil {
		t.Errorf("Should not have error return value, but received \"%v\"", err)
		return
	}

	if len(id) != 16 {
		t.Errorf("Expected length of 16, but got %d", len(id))
		return
	}

	idString := id.String()
	if len(idString) != 32 {
		t.Errorf("Expected length of 32, but got %d", len(idString))
		return
	}
}
