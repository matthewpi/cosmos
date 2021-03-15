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

// Package server ...
package server

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/matthewpi/cosmos"
	"github.com/matthewpi/cosmos/internal/listener"
)

var (
	ErrNoListeners = errors.New("server: no listeners defined")
)

var (
	defaultTLSConfig = &tls.Config{
		NextProtos: []string{
			"h2",
			"http/1.1",
		},

		CipherSuites: []uint16{
			// TLS 1.0 - 1.2
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,

			// TLS 1.3
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},

		PreferServerCipherSuites: true,

		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,

		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
	}
)

// Server .
type Server struct {
	config *Config

	listeners []net.Listener
	servers   []*http.Server
	router    *chi.Mux
}

// New .
func New(ops ...Opt) (*Server, error) {
	s := &Server{
		config: &Config{},
		router: chi.NewRouter(),
	}
	for _, op := range ops {
		if err := op(s); err != nil {
			return nil, err
		}
	}
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	return s, nil
}

// Listen .
func (s *Server) Listen() []error {
	var errs []error
	for _, lc := range s.config.Listeners {
		l, err := net.Listen(lc.Network.String(), lc.Address)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if l2, ok := l.(*net.TCPListener); ok {
			l = listener.Wrap(l2, lc.KeepAlive)
		}
		s.listeners = append(s.listeners, l)
	}
	return errs
}

// Serve .
func (s *Server) Serve(ctx context.Context) error {
	if s.listeners == nil || len(s.listeners) < 1 {
		return ErrNoListeners
	}
	l := zap.NewStdLog(cosmos.Log())
	g, ctx := errgroup.WithContext(ctx)
	for i, lc := range s.config.Listeners {
		hs := &http.Server{
			Addr:    lc.Address,
			Handler: s.router,

			TLSConfig: nil,

			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 3 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       30 * time.Second,

			MaxHeaderBytes: http.DefaultMaxHeaderBytes,

			ErrorLog: l,
		}
		s.servers = append(s.servers, hs)

		if certPath, keyPath := lc.CertPath, lc.KeyPath; certPath == "" || keyPath == "" {
			g.Go(func() error {
				if err := hs.Serve(s.listeners[i]); err != http.ErrServerClosed {
					return err
				}
				return nil
			})
		} else {
			hs.TLSConfig = defaultTLSConfig
			g.Go(func() error {
				if err := hs.ServeTLS(s.listeners[i], certPath, keyPath); err != http.ErrServerClosed {
					return err
				}
				return nil
			})
		}
	}
	return g.Wait()
}

// Close .
func (s *Server) Close(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, hs := range s.servers {
		g.Go(func() error {
			hs.SetKeepAlivesEnabled(false)
			return hs.Shutdown(ctx)
		})
	}
	return g.Wait()
}
