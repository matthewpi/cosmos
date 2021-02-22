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

package snowflake_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/matthewpi/cosmos/internal/snowflake"
)

func TestNew(t *testing.T) {
	s := snowflake.New()
	if s == 0 {
		t.Errorf("Expected non-zero value, but got one")
		return
	}
}

func TestNewAtTime(t *testing.T) {
	expected := time.Date(2021, 5, 24, 12, 38, 38, 333*int(time.Millisecond), time.UTC)
	s := snowflake.NewAtTime(expected)
	if s == 0 {
		t.Errorf("Expected non-zero value, but got one")
		return
	}

	if expected != s.Time().UTC() {
		t.Errorf("Expected \"%v\", but got \"%v\"", expected, s.Time().UTC())
		return
	}
}

func TestParse(t *testing.T) {
	for i, tc := range []struct {
		snowflake string
		expect    snowflake.Snowflake
	}{
		{
			snowflake: "52466462028201985",
			expect:    snowflake.Snowflake(52466462028201985),
		},
		{
			snowflake: "null",
			expect:    snowflake.Nil,
		},
		{
			snowflake: "invalid snowflake",
			expect:    snowflake.Nil,
		},
	} {
		result := snowflake.Parse(tc.snowflake)

		if tc.expect != result {
			t.Errorf("Test #%d: Expected \"%s\", but got \"%s\"", i, tc.expect, result)
			continue
		}
	}
}

func TestSnowflake_Valid(t *testing.T) {
	for i, tc := range []struct {
		snowflake snowflake.Snowflake
		expect    bool
	}{
		{
			snowflake: snowflake.Snowflake(52466462028201985),
			expect:    true,
		},
		{
			snowflake: snowflake.Nil,
			expect:    false,
		},
		{
			snowflake: snowflake.Zero,
			expect:    false,
		},
	} {
		result := tc.snowflake.Valid()

		if tc.expect != result {
			t.Errorf("Test #%d: Expected \"%t\", but got \"%t\"", i, tc.expect, result)
			continue
		}
	}
}

func TestSnowflake_Time(t *testing.T) {
	for i, tc := range []struct {
		snowflake snowflake.Snowflake
		expect    time.Time
	}{
		{
			snowflake: snowflake.NewAtTime(time.Date(2021, 5, 24, 12, 38, 38, 333*int(time.Millisecond), time.UTC)),
			expect:    time.Date(2021, 5, 24, 12, 38, 38, 333*int(time.Millisecond), time.UTC).UTC(),
		},
		{
			snowflake: snowflake.Nil,
			expect:    snowflake.Nil.Time(),
		},
		{
			snowflake: snowflake.Zero,
			expect:    snowflake.Zero.Time(),
		},
	} {
		result := tc.snowflake.Time().UTC()

		if tc.expect != result {
			t.Errorf("Test #%d: Expected \"%s\", but got \"%s\"", i, tc.expect, result)
			continue
		}
	}
}

func TestSnowflake_String(t *testing.T) {
	for i, tc := range []struct {
		snowflake snowflake.Snowflake
		expect    string
	}{
		{
			snowflake: snowflake.Snowflake(52466462028201985),
			expect:    "52466462028201985",
		},
		{
			snowflake: snowflake.Nil,
			expect:    "",
		},
		{
			snowflake: snowflake.Zero,
			expect:    "",
		},
	} {
		result := tc.snowflake.String()

		if tc.expect != result {
			t.Errorf("Test #%d: Expected \"%s\", but got \"%s\"", i, tc.expect, result)
			continue
		}
	}
}

func TestSnowflake_MarshalJSON(t *testing.T) {
	for i, tc := range []struct {
		snowflake snowflake.Snowflake
		expect    []byte
		expectErr error
	}{
		{
			snowflake: snowflake.Snowflake(52466462028201985),
			expect:    []byte("\"52466462028201985\""),
			expectErr: nil,
		},
		{
			snowflake: snowflake.Nil,
			expect:    []byte("null"),
			expectErr: nil,
		},
		{
			snowflake: snowflake.Zero,
			expect:    []byte("null"),
			expectErr: nil,
		},
	} {
		result, err := tc.snowflake.MarshalJSON()

		if tc.expectErr != nil && err == nil {
			t.Errorf("Test #%d: Expected error return value, but got \"%v\"", i, err)
			continue
		}
		if tc.expectErr == nil && err != nil {
			t.Errorf("Test #%d: Should not have error return value, but received \"%v\"", i, err)
			continue
		}
		if !bytes.Equal(tc.expect, result) {
			t.Errorf("Test #%d: Expected \"%s\", but got \"%s\"", i, tc.expect, result)
			continue
		}
	}
}

func TestSnowflake_UnmarshalJSON(t *testing.T) {
	for i, tc := range []struct {
		v         []byte
		expect    snowflake.Snowflake
		expectErr error
	}{
		{
			v:         []byte("\"52466462028201985\""),
			expect:    snowflake.Snowflake(52466462028201985),
			expectErr: nil,
		},
		{
			v:         []byte("null"),
			expect:    snowflake.Nil,
			expectErr: nil,
		},
	} {
		result := snowflake.Snowflake(0)
		err := result.UnmarshalJSON(tc.v)

		if tc.expectErr != nil && err == nil {
			t.Errorf("Test #%d: Expected error return value, but got \"%v\"", i, err)
			continue
		}
		if tc.expectErr == nil && err != nil {
			t.Errorf("Test #%d: Should not have error return value, but received \"%v\"", i, err)
			continue
		}
		if tc.expect != result {
			t.Errorf("Test #%d: Expected \"%s\", but got \"%s\"", i, tc.expect, result)
			continue
		}
	}
}
