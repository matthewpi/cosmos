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

package server

import (
	"fmt"

	"github.com/matthewpi/cosmos/internal/config/lexer"
	"github.com/matthewpi/cosmos/internal/server/listener"
)

// Config represents the configuration for a Server.
type Config struct {
	// Listeners is a slice of listeners to bind to.
	Listeners []listener.Listener
}

// FromLexer .
func FromLexer(b lexer.Block) (*Server, error) {
	var opts []Opt

	var tokens []lexer.Token
	for _, s := range b.Segments {
		tokens = append(tokens, s...)
	}
	d := lexer.NewDispenser(tokens)

	for d.Next() {
		dir := d.Val()
		switch dir {
		case "listen":
			l := listener.Listener{
				Network: listener.NetworkTCP,
			}
			for d.NextArg() {
				v := d.Val()
				if len(v) < 1 {
					return nil, fmt.Errorf("empty argument after listen directive")
				}
				l.Address = v
			}
			if l.Address == "" {
				return nil, fmt.Errorf("missing argument after \"listen\" directive")
			}
			for nesting := d.Nesting(); d.NextBlock(nesting); {
				subdir := d.Val()
				switch subdir {
				case "metrics":
					l.Metrics = "/metrics"
					for d.NextArg() {
						return nil, fmt.Errorf("unexpected argument after metrics directive")
					}
				case "{":
					return nil, fmt.Errorf("unexpected start of block")
				case "}":
					return nil, fmt.Errorf("unexpected closing of block")
				default:
					return nil, fmt.Errorf("unknown sub-directive: \"" + subdir + "\"")
				}
			}
			opts = append(opts, WithListener(l))
		default:
			return nil, fmt.Errorf("unknown directive: \"" + d.Val() + "\"")
		}
	}

	return New(opts...)
}
