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
	"crypto/rand"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"math"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

const (
	// epoch is the starting epoch. (first second of 2021)
	epoch = 1609459200000 * time.Millisecond

	// Nil is a nil/null snowflake.
	Nil Snowflake = -1

	// Zero is a zero snowflake. (unset/omitted)
	Zero Snowflake = 0
)

var (
	nodeID = getNodeID()

	sequence *int64
)

func init() {
	seq := int64(0)
	sequence = &seq
}

// Snowflake represents a snowflake.
type Snowflake int64

var _ fmt.Stringer = (*Snowflake)(nil)
var _ json.Marshaler = (*Snowflake)(nil)
var _ json.Unmarshaler = (*Snowflake)(nil)

// New returns a new snowflake.
func New() Snowflake {
	return NewAtTime(time.Now())
}

// NewAtTime returns a new snowflake using the specified time.
func NewAtTime(t time.Time) Snowflake {
	id := int64(((time.Duration(t.UnixNano()) - epoch) / time.Millisecond) << 22)
	id |= int64(nodeID << 20)

	atomic.AddInt64(sequence, 1)
	id |= atomic.LoadInt64(sequence)

	return Snowflake(id)
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

	unixNano := ((time.Duration(s) >> 22) * time.Millisecond) + epoch
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

// ---------------------------------- \\

func getMac() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	var mac string
	for _, i := range interfaces {
		addr := i.HardwareAddr.String()
		if addr != "" {
			mac += addr
		}
	}

	return mac
}

func getNodeID() int {
	var id int
	mac := getMac()
	if mac != "" {
		h := fnv.New32a()
		id = int(h.Sum32())
	}

	if id == 0 {
		n, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			panic(err)
			return 0
		}

		return int(n.Int64())
	}

	return id & (int)(math.Pow(2, float64(10))-1)
}
