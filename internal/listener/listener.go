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

// Package listener ...
package listener

import (
	"net"
	"time"
)

// TCPKeepAliveListener .
type TCPKeepAliveListener struct {
	*net.TCPListener

	KeepAlive time.Duration
}

// AcceptTCP .
func (l TCPKeepAliveListener) AcceptTCP() (*net.TCPConn, error) {
	c, err := l.TCPListener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	if err := c.SetKeepAlive(true); err != nil {
		return nil, err
	}
	if err := c.SetKeepAlivePeriod(l.KeepAlive); err != nil {
		return nil, err
	}
	return c, nil
}

// Accept .
func (l TCPKeepAliveListener) Accept() (net.Conn, error) {
	return l.AcceptTCP()
}

// Wrap .
func Wrap(listener net.Listener, keepAlive time.Duration) net.Listener {
	if tcpListener, ok := listener.(*net.TCPListener); ok {
		return TCPKeepAliveListener{
			TCPListener: tcpListener,
			KeepAlive:   keepAlive,
		}
	}
	return listener
}
