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

// Package config ...
package config

import (
	"os"

	"github.com/matthewpi/cosmos/internal/config/lexer"
)

// Config .
type Config []lexer.Block

// Load .
func Load(file string) (Config, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	blocks, err := lexer.Parse(file, f)
	if err != nil {
		return nil, err
	}
	return blocks, nil
}

func (c Config) Key(key string) lexer.Block {
	for _, v := range c {
		for _, k := range v.Keys {
			if k != key {
				continue
			}
			return v
		}
	}
	return lexer.Block{}
}
