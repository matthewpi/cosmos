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

// Package lexer ...
package lexer

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

// Token .
type Token struct {
	// File .
	File string
	// Line .
	Line int
	// Text .
	Text string

	inSnippet   bool
	snippetName string
}

// Tokenize .
func Tokenize(filename string, v []byte) ([]Token, error) {
	l := lexer{}
	if err := l.load(bytes.NewReader(v)); err != nil {
		return nil, err
	}
	var tokens []Token
	for l.next() {
		l.token.File = filename
		tokens = append(tokens, l.token)
	}
	return tokens, nil
}

// lexer .
type lexer struct {
	reader *bufio.Reader

	token        Token
	line         int
	skippedLines int
}

// load .
func (l *lexer) load(input io.Reader) error {
	l.reader = bufio.NewReader(input)
	l.line = 1

	// discard byte order mark, if present
	firstCh, _, err := l.reader.ReadRune()
	if err != nil {
		return err
	}
	if firstCh != 0xFEFF {
		err := l.reader.UnreadRune()
		if err != nil {
			return err
		}
	}

	return nil
}

// next .
func (l *lexer) next() bool {
	var (
		val []rune

		comment  bool
		quoted   bool
		btQuoted bool
		escaped  bool
	)

	makeToken := func() bool {
		l.token.Text = string(val)
		return true
	}

	for {
		ch, _, err := l.reader.ReadRune()
		if err != nil {
			if len(val) > 0 {
				return makeToken()
			}
			if err == io.EOF {
				return false
			}
			panic(err)
		}

		if !escaped && !btQuoted && ch == '\\' {
			escaped = true
			continue
		}

		if quoted || btQuoted {
			if quoted && escaped {
				// all is literal in quoted area,
				// so only escape quotes
				if ch != '"' {
					val = append(val, '\\')
				}
				escaped = false
			} else {
				if quoted && ch == '"' {
					return makeToken()
				}
				if btQuoted && ch == '`' {
					return makeToken()
				}
			}
			if ch == '\n' {
				l.line += 1 + l.skippedLines
				l.skippedLines = 0
			}
			val = append(val, ch)
			continue
		}

		if unicode.IsSpace(ch) {
			if ch == '\r' {
				continue
			}
			if ch == '\n' {
				if escaped {
					l.skippedLines++
					escaped = false
				} else {
					l.line += 1 + l.skippedLines
					l.skippedLines = 0
				}
				comment = false
			}
			if len(val) > 0 {
				return makeToken()
			}
			continue
		}

		if ch == '#' && len(val) == 0 {
			comment = true
		}
		if comment {
			continue
		}

		if len(val) == 0 {
			l.token = Token{Line: l.line}
			if ch == '"' {
				quoted = true
				continue
			}
			if ch == '`' {
				btQuoted = true
				continue
			}
		}

		if escaped {
			val = append(val, '\\')
			escaped = false
		}

		val = append(val, ch)
	}
}
