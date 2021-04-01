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

// Package forwarded implements a "Forwarded HTTP Extension" (RFC 7239)
// compatible parser for the "Forwarded" HTTP header used to denote
// what proxies a HTTP Request was proxied through.
//
// The Forwarded HTTP Extension is a better and standardized version
// of the "X-Forwarded-For" header which lacks detailed information
// about what proxies touched the request.
//
// Resources
// - https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Forwarded
// - https://tools.ietf.org/html/rfc7239
package forwarded

import (
	"net/http"
	"strings"
)

// Proto represents the protocol used with a specific Proxy.
type Proto string

const (
	// HTTP represents the HTTP protocol.
	HTTP Proto = "http"
	// HTTPS represents the HTTPS protocol.
	HTTPS Proto = "https"
)

// Proxy represents an individual proxy server that has forwarded a HTTP request.
//
// See https://tools.ietf.org/html/rfc7239#section-5 for more information.
type Proxy struct {
	// By identifies the user-agent facing interface of the proxy.
	By string
	// For identifies the node making the request to the proxy.
	For string
	// Host is the host request header field as received by the proxy.
	Host string
	// Proto indicates what protocol was used to make the request.
	Proto Proto
}

// String formats the proxy in the Forwarded header format.
//
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Forwarded
// for more information.
func (p Proxy) String() string {
	var b strings.Builder
	for k, v := range map[string]string{
		"by":    p.By,
		"for":   p.For,
		"host":  p.Host,
		"proto": string(p.Proto),
	} {
		if v == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(';')
		}
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteByte('"')
		b.WriteString(v)
		b.WriteByte('"')
	}
	return b.String()
}

// Get returns the proxies from the "Forwarded" headers of a request.
func Get(r *http.Request) ([]Proxy, error) {
	forwarded := r.Header.Values("Forwarded")
	if len(forwarded) < 1 {
		return nil, nil
	}

	// If there are multiple Forwarded headers in the request;
	// "merge" them into a single header, separated by commas.
	//
	// https://tools.ietf.org/html/rfc7239#section-7.1
	return Parse(strings.TrimSpace(strings.Join(forwarded, ",")))
}

// Parse parses a "Forwarded"-style header, usually from a HTTP request.
//
// See https://tools.ietf.org/html/rfc7239#section-5 and
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Forwarded
// for more information.
func Parse(header string) ([]Proxy, error) {
	var proxies []Proxy

	var key []rune
	var text []rune
	var content bool
	var inQuotes bool
	keys := map[string]string{}

	for _, c := range header {
		switch {
		case c == ' ':
			if !content {
				continue
			}
			text = append(text, c)
		case c == '"':
			inQuotes = !inQuotes
		case c == '=':
			if inQuotes {
				text = append(text, c)
				continue
			}
			content = true
		case c == ',':
			if inQuotes {
				text = append(text, c)
				continue
			}
			proxies = append(proxies, Proxy{
				By:    keys["by"],
				For:   keys["for"],
				Host:  keys["host"],
				Proto: Proto(keys["proto"]),
			})
			keys = map[string]string{}
			fallthrough
		case c == ';':
			if inQuotes {
				text = append(text, c)
				continue
			}
			keys[strings.ToLower(string(key))] = string(text)
			key = key[:0]
			text = text[:0]
			content = false
		default:
			if content {
				text = append(text, c)
			} else {
				key = append(key, c)
			}
		}
	}
	keys[strings.ToLower(string(key))] = string(text)
	proxies = append(proxies, Proxy{
		By:    keys["by"],
		For:   keys["for"],
		Host:  keys["host"],
		Proto: Proto(keys["proto"]),
	})
	return proxies, nil
}
