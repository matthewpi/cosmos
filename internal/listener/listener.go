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

// Network .
type Network string

func (n Network) String() string {
	return string(n)
}

const (
	NetworkTCP        Network = "tcp"
	NetworkTCP4       Network = "tcp4"
	NetworkTCP6       Network = "tcp6"
	NetworkUNIX       Network = "unix"
	NetworkUNIXPacket Network = "unixpacket"
)

// Listener .
type Listener struct {
	// Network .
	Network Network
	// Address .
	Address string

	// KeepAlive .
	KeepAlive time.Duration

	// CertPath .
	CertPath string
	// KeyPath .
	KeyPath string
}

// TCPKeepAliveListener .
type TCPKeepAliveListener struct {
	*net.TCPListener

	KeepAlive time.Duration
}

var _ net.Listener = (*TCPKeepAliveListener)(nil)

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

// Wrap .
func Wrap(l *net.TCPListener, keepAlive time.Duration) net.Listener {
	return TCPKeepAliveListener{
		TCPListener: l,
		KeepAlive:   keepAlive,
	}
}
