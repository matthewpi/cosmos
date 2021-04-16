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

package log

import (
	"fmt"
	"strings"

	"github.com/matthewpi/cosmos/internal/config/lexer"
)

// Config .
type Config struct {
	// Level .
	Level Level `json:"level"`
}

func FromLexer(b lexer.Block) (*Config, error) {
	c := &Config{}
	for _, s := range b.Segments {
		d := s.Directive()
		switch d {
		case "output":
		case "level":
			if len(s) < 2 {
				return nil, fmt.Errorf("missing level after level directive")
			}
			if len(s) > 2 {
				return nil, fmt.Errorf("too many arguments after level directive")
			}
			k := strings.ToLower(s[1].Text)
			l, ok := Levels[k]
			if !ok {
				return nil, fmt.Errorf("unknown level: \"%s\"", k)
			}
			c.Level = l
		}
	}
	return c, nil
}
