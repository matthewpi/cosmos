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

package lexer

import (
	"testing"
)

type lexerTestCase struct {
	input    []byte
	expected []Token
}

func TestLexer(t *testing.T) {
	testCases := []lexerTestCase{
		{
			input: []byte(`host:123`),
			expected: []Token{
				{Line: 1, Text: "host:123"},
			},
		},
		{
			input: []byte(`host:123

					directive`),
			expected: []Token{
				{Line: 1, Text: "host:123"},
				{Line: 3, Text: "directive"},
			},
		},
		{
			input: []byte(`host:123 {
						directive
					}`),
			expected: []Token{
				{Line: 1, Text: "host:123"},
				{Line: 1, Text: "{"},
				{Line: 2, Text: "directive"},
				{Line: 3, Text: "}"},
			},
		},
		{
			input: []byte(`host:123 { directive }`),
			expected: []Token{
				{Line: 1, Text: "host:123"},
				{Line: 1, Text: "{"},
				{Line: 1, Text: "directive"},
				{Line: 1, Text: "}"},
			},
		},
		{
			input: []byte(`host:123 {
						#comment
						directive
						# comment
						foobar # another comment
					}`),
			expected: []Token{
				{Line: 1, Text: "host:123"},
				{Line: 1, Text: "{"},
				{Line: 3, Text: "directive"},
				{Line: 5, Text: "foobar"},
				{Line: 6, Text: "}"},
			},
		},
		{
			input: []byte(`host:123 {
						# hash inside string is not a comment
						redir / /some/#/path
					}`),
			expected: []Token{
				{Line: 1, Text: "host:123"},
				{Line: 1, Text: "{"},
				{Line: 3, Text: "redir"},
				{Line: 3, Text: "/"},
				{Line: 3, Text: "/some/#/path"},
				{Line: 4, Text: "}"},
			},
		},
		{
			input: []byte("# comment at beginning of file\n# comment at beginning of line\nhost:123"),
			expected: []Token{
				{Line: 3, Text: "host:123"},
			},
		},
		{
			input: []byte(`a "quoted value" b
					foobar`),
			expected: []Token{
				{Line: 1, Text: "a"},
				{Line: 1, Text: "quoted value"},
				{Line: 1, Text: "b"},
				{Line: 2, Text: "foobar"},
			},
		},
		{
			input: []byte(`A "quoted \"value\" inside" B`),
			expected: []Token{
				{Line: 1, Text: "A"},
				{Line: 1, Text: `quoted "value" inside`},
				{Line: 1, Text: "B"},
			},
		},
		{
			input: []byte("An escaped \"newline\\\ninside\" quotes"),
			expected: []Token{
				{Line: 1, Text: "An"},
				{Line: 1, Text: "escaped"},
				{Line: 1, Text: "newline\\\ninside"},
				{Line: 2, Text: "quotes"},
			},
		},
		{
			input: []byte("An escaped newline\\\noutside quotes"),
			expected: []Token{
				{Line: 1, Text: "An"},
				{Line: 1, Text: "escaped"},
				{Line: 1, Text: "newline"},
				{Line: 1, Text: "outside"},
				{Line: 1, Text: "quotes"},
			},
		},
		{
			input: []byte("line1\\\nescaped\nline2\nline3"),
			expected: []Token{
				{Line: 1, Text: "line1"},
				{Line: 1, Text: "escaped"},
				{Line: 3, Text: "line2"},
				{Line: 4, Text: "line3"},
			},
		},
		{
			input: []byte("line1\\\nescaped1\\\nescaped2\nline4\nline5"),
			expected: []Token{
				{Line: 1, Text: "line1"},
				{Line: 1, Text: "escaped1"},
				{Line: 1, Text: "escaped2"},
				{Line: 4, Text: "line4"},
				{Line: 5, Text: "line5"},
			},
		},
		{
			input: []byte(`"unescapable\ in quotes"`),
			expected: []Token{
				{Line: 1, Text: `unescapable\ in quotes`},
			},
		},
		{
			input: []byte(`"don't\escape"`),
			expected: []Token{
				{Line: 1, Text: `don't\escape`},
			},
		},
		{
			input: []byte(`"don't\\escape"`),
			expected: []Token{
				{Line: 1, Text: `don't\\escape`},
			},
		},
		{
			input: []byte(`un\escapable`),
			expected: []Token{
				{Line: 1, Text: `un\escapable`},
			},
		},
		{
			input: []byte(`A "quoted value with line
					break inside" {
						foobar
					}`),
			expected: []Token{
				{Line: 1, Text: "A"},
				{Line: 1, Text: "quoted value with line\n\t\t\t\t\tbreak inside"},
				{Line: 2, Text: "{"},
				{Line: 3, Text: "foobar"},
				{Line: 4, Text: "}"},
			},
		},
		{
			input: []byte(`"C:\php\php-cgi.exe"`),
			expected: []Token{
				{Line: 1, Text: `C:\php\php-cgi.exe`},
			},
		},
		{
			input: []byte(`empty "" string`),
			expected: []Token{
				{Line: 1, Text: `empty`},
				{Line: 1, Text: ``},
				{Line: 1, Text: `string`},
			},
		},
		{
			input: []byte("skip those\r\nCR characters"),
			expected: []Token{
				{Line: 1, Text: "skip"},
				{Line: 1, Text: "those"},
				{Line: 2, Text: "CR"},
				{Line: 2, Text: "characters"},
			},
		},
		{
			input: []byte("\xEF\xBB\xBF:8080"), // test with leading byte order mark
			expected: []Token{
				{Line: 1, Text: ":8080"},
			},
		},
		{
			input: []byte("simple `backtick quoted` string"),
			expected: []Token{
				{Line: 1, Text: `simple`},
				{Line: 1, Text: `backtick quoted`},
				{Line: 1, Text: `string`},
			},
		},
		{
			input: []byte("multiline `backtick\nquoted\n` string"),
			expected: []Token{
				{Line: 1, Text: `multiline`},
				{Line: 1, Text: "backtick\nquoted\n"},
				{Line: 3, Text: `string`},
			},
		},
		{
			input: []byte("nested `\"quotes inside\" backticks` string"),
			expected: []Token{
				{Line: 1, Text: `nested`},
				{Line: 1, Text: `"quotes inside" backticks`},
				{Line: 1, Text: `string`},
			},
		},
		{
			input: []byte("reverse-nested \"`backticks` inside\" quotes"),
			expected: []Token{
				{Line: 1, Text: `reverse-nested`},
				{Line: 1, Text: "`backticks` inside"},
				{Line: 1, Text: `quotes`},
			},
		},
	}

	for i, testCase := range testCases {
		actual, err := Tokenize("", testCase.input)
		if err != nil {
			t.Errorf("%v", err)
		}
		lexerCompare(t, i, testCase.expected, actual)
	}
}

func lexerCompare(t *testing.T, n int, expected, actual []Token) {
	if len(expected) != len(actual) {
		t.Errorf("Test case %d: expected %d token(s) but got %d", n, len(expected), len(actual))
	}

	for i := 0; i < len(actual) && i < len(expected); i++ {
		if actual[i].Line != expected[i].Line {
			t.Errorf("Test case %d token %d ('%s'): expected line %d but was line %d",
				n, i, expected[i].Text, expected[i].Line, actual[i].Line)
			break
		}
		if actual[i].Text != expected[i].Text {
			t.Errorf("Test case %d token %d: expected text '%s' but was '%s'",
				n, i, expected[i].Text, actual[i].Text)
			break
		}
	}
}
