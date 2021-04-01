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

package forwarded_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/matthewpi/cosmos/internal/forwarded"
)

func TestGet(t *testing.T) {
	// TODO
}

func TestParse(t *testing.T) {
	for i, tc := range []struct {
		header    string
		expect    []forwarded.Proxy
		expectErr error
	}{
		{
			header: "for=\"_mdn\"",
			expect: []forwarded.Proxy{
				{
					By:    "",
					For:   "_mdn",
					Host:  "",
					Proto: "",
				},
			},
			expectErr: nil,
		},
		{
			header: "For=\"[2001:db8:cafe::17]:4711\"",
			expect: []forwarded.Proxy{
				{
					By:    "",
					For:   "[2001:db8:cafe::17]:4711",
					Host:  "",
					Proto: "",
				},
			},
			expectErr: nil,
		},
		{
			header: "by=\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36\"",
			expect: []forwarded.Proxy{
				{
					By:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
					For:   "",
					Host:  "",
					Proto: "",
				},
			},
			expectErr: nil,
		},
		{
			header: "by=\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36\";for=\"_mdn\"",
			expect: []forwarded.Proxy{
				{
					By:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
					For:   "_mdn",
					Host:  "",
					Proto: "",
				},
			},
			expectErr: nil,
		},
		{
			header: "by=\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36\";for=_mdn;host=github.com;proto=http",
			expect: []forwarded.Proxy{
				{
					By:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
					For:   "_mdn",
					Host:  "github.com",
					Proto: forwarded.HTTP,
				},
			},
			expectErr: nil,
		},
		{
			header: "by=\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36\";for=\"_mdn\";host=github.com:443;proto=https",
			expect: []forwarded.Proxy{
				{
					By:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
					For:   "_mdn",
					Host:  "github.com:443",
					Proto: forwarded.HTTPS,
				},
			},
			expectErr: nil,
		},
		{
			header: "By=\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36\";For=\"_mdn\";Host=github.com:443;Proto=https",
			expect: []forwarded.Proxy{
				{
					By:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36",
					For:   "_mdn",
					Host:  "github.com:443",
					Proto: forwarded.HTTPS,
				},
			},
			expectErr: nil,
		},
	} {
		result, err := forwarded.Parse(tc.header)
		if tc.expectErr != nil && err == nil {
			t.Errorf("Test #%d: Expected error return value, but got \"%v\"", i, err)
			continue
		}
		if tc.expectErr == nil && err != nil {
			t.Errorf("Test #%d: Should not have error return value, but received \"%v\"", i, err)
			continue
		}

		if tc.expect == nil {
			if result != nil {
				received, _ := json.Marshal(result)
				t.Errorf("Test #%d: Expected nil, but got \"%s\"", i, string(received))
			}
			continue
		}

		if !reflect.DeepEqual(tc.expect, result) {
			expected, _ := json.Marshal(tc.expect)
			received, _ := json.Marshal(result)
			t.Errorf("Test #%d: Expected \"%s\", but got \"%s\"", i, string(expected), string(received))
		}
	}
}
