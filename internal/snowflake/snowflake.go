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

// Package snowflake ...
package snowflake

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	// epoch is the starting epoch. (first second of 2021)
	epoch = 1609459200000 * time.Millisecond

	// Nil is a nil/null snowflake.
	Nil = ^Snowflake(0)
)

var (
	NodeID = 1

	sequence *uint64
)

func init() {
	seq := uint64(0)
	sequence = &seq
}

// Snowflake represents a snowflake.
//
// Timestamp - 63 to 22 (42 bits)
// Node ID   - 21 to 12 (10 bits)
// Increment - 11 to  0 (12 bits)
type Snowflake uint64

var _ fmt.Stringer = (*Snowflake)(nil)
var _ json.Marshaler = (*Snowflake)(nil)
var _ json.Unmarshaler = (*Snowflake)(nil)

// New returns a new snowflake.
func New() Snowflake {
	id := newAtTime(time.Now())
	id |= uint64(NodeID << 53)

	atomic.AddUint64(sequence, 1)
	id |= atomic.LoadUint64(sequence)

	return Snowflake(id)
}

// NewAtTime returns a new snowflake using the specified time.
func NewAtTime(t time.Time) Snowflake {
	return Snowflake(newAtTime(t))
}

func newAtTime(t time.Time) uint64 {
	return uint64(((time.Duration(t.UnixNano()) - epoch) / time.Millisecond) << 22)
}

// Parse parses a string into a snowflake.
func Parse(snowflake string) Snowflake {
	if snowflake == "null" {
		return Nil
	}

	i, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return Nil
	}
	return Snowflake(i)
}

// Valid returns true if the snowflake is valid.
func (s Snowflake) Valid() bool {
	return int64(s) > 0
}

// Time returns the time at which the snowflake was generated at.
func (s Snowflake) Time() time.Time {
	if !s.Valid() {
		return time.Time{}
	}

	unixNano := (time.Duration(s)>>22)*time.Millisecond + epoch
	return time.Unix(0, int64(unixNano))
}

// String satisfies fmt.Stringer.
func (s Snowflake) String() string {
	if !s.Valid() {
		return ""
	}

	return strconv.FormatInt(int64(s), 10)
}

// MarshalJSON satisfies json.Marshaler.
func (s Snowflake) MarshalJSON() ([]byte, error) {
	if !s.Valid() {
		return []byte("null"), nil
	}

	return []byte(`"` + strconv.FormatInt(int64(s), 10) + `"`), nil
}

// UnmarshalJSON satisfies json.Unmarshaler.
func (s *Snowflake) UnmarshalJSON(v []byte) error {
	*s = Parse(strings.Trim(string(v), `"`))
	return nil
}
