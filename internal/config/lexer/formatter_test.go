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

// Copyright 2015 Matthew Holt and The Caddy Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lexer

import (
	"strings"
	"testing"
)

func TestFormatter(t *testing.T) {
	for i, tc := range []struct {
		description string
		input       string
		expect      string
	}{
		{
			description: "very simple",
			input: `abc   def
	g hi jkl
mn`,
			expect: `abc def
g hi jkl
mn`,
		},
		{
			description: "basic indentation, line breaks, and nesting",
			input: `  a
b

	c {
		d
}

e { f
}



g {
h {
i
}
}

j { k {
l
}
}

m {
	n { o
	}
	p { q r
s }
}

	{
{ t
		u

	v

w
}
}`,
			expect: `a
b

c {
	d
}

e {
	f
}

g {
	h {
		i
	}
}

j {
	k {
		l
	}
}

m {
	n {
		o
	}
	p {
		q r
		s
	}
}

{
	{
		t
		u

		v

		w
	}
}`,
		},
		{
			description: "block spacing",
			input: `a{
	b
}

c{ d
}`,
			expect: `a {
	b
}

c {
	d
}`,
		},
		{
			description: "advanced spacing",
			input: `abc {
	def
}ghi{
	jkl mno
pqr}`,
			expect: `abc {
	def
}

ghi {
	jkl mno
	pqr
}`,
		},
		{
			description: "env var placeholders",
			input: `{$A}

b {
{$C}
}

d { {$E}
}

{ {$F}
}
`,
			expect: `{$A}

b {
	{$C}
}

d {
	{$E}
}

{
	{$F}
}`,
		},
		{
			description: "comments",
			input: `#a "\n"

 #b {
	c
}

d {
e#f
# g
}

h { # i
}`,
			expect: `#a "\n"

#b {
c
}

d {
	e#f
	# g
}

h {
	# i
}`,
		},
		{
			description: "quotes and escaping",
			input: `"a \"b\" "#c
	d

e {
"f"
}

g { "h"
}

i {
	"foo
bar"
}

j {
"\"k\" l m"
}`,
			expect: `"a \"b\" "#c
d

e {
	"f"
}

g {
	"h"
}

i {
	"foo
bar"
}

j {
	"\"k\" l m"
}`,
		},
		{
			description: "bad nesting (too many open)",
			input: `a
{
	{
}`,
			expect: `a {
	{
	}
`,
		},
		{
			description: "bad nesting (too many close)",
			input: `a
{
	{
}}}`,
			expect: `a {
	{
	}
}
}
`,
		},
		{
			description: "json",
			input: `foo
bar      "{\"key\":34}"
`,
			expect: `foo
bar "{\"key\":34}"`,
		},
		{
			description: "escaping after spaces",
			input:       `foo \"literal\"`,
			expect:      `foo \"literal\"`,
		},
		{
			description: "simple placeholders as standalone tokens",
			input:       `foo {bar}`,
			expect:      `foo {bar}`,
		},
		{
			description: "simple placeholders within tokens",
			input:       `foo{bar} foo{bar}baz`,
			expect:      `foo{bar} foo{bar}baz`,
		},
		{
			description: "placeholders and malformed braces",
			input:       `foo{bar} foo{ bar}baz`,
			expect: `foo{bar} foo {
	bar
}

baz`,
		},
		{
			description: "hash within string is not a comment",
			input:       `redir / /some/#/path`,
			expect:      `redir / /some/#/path`,
		},
		{
			description: "brace does not fold into comment above",
			input: `# comment
{
	foo
}`,
			expect: `# comment
{
	foo
}`,
		},
	} {
		// the formatter should output a trailing newline,
		// even if the tests aren't written to expect that
		if !strings.HasSuffix(tc.expect, "\n") {
			tc.expect += "\n"
		}

		actual := Format([]byte(tc.input))

		if string(actual) != tc.expect {
			t.Errorf("\n[TEST %d: %s]\n====== EXPECTED ======\n%s\n====== ACTUAL ======\n%s^^^^^^^^^^^^^^^^^^^^^",
				i, tc.description, string(tc.expect), string(actual))
		}
	}
}