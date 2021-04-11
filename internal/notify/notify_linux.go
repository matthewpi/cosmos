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

package notify

import (
	"io"
	"net"
	"os"
	"strings"
)

// notify .
func notify(path string, r io.Reader) error {
	s := &net.UnixAddr{
		Name: path,
		Net:  "unixgram",
	}
	c, err := net.DialUnix(s.Net, nil, s)
	if err != nil {
		return err
	}
	defer c.Close()

	if _, err := io.Copy(c, r); err != nil {
		return err
	}
	return nil
}

func socketNotify(payload string) error {
	v, ok := os.LookupEnv("NOTIFY_SOCKET")
	if !ok || v == "" {
		return nil
	}
	if err := notify(v, strings.NewReader(payload)); err != nil {
		return err
	}
	return nil
}

// readiness .
func readiness() error {
	return socketNotify("READY=1")
}

// reloading .
func reloading() error {
	return socketNotify("RELOADING=1")
}

// stopping .
func stopping() error {
	return socketNotify("STOPPING=1")
}
