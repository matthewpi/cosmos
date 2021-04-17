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
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/matthewpi/cosmos/internal/server/listener"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/matthewpi/cosmos"
	"github.com/matthewpi/cosmos/internal/metrics"
)

var (
	// ErrNoListeners is returned when Server.Serve() is called before
	// Server.Listen()
	ErrNoListeners = errors.New("server: no listeners defined")

	// ErrMissingListener is returned by Server.Serve() when the number of
	// listeners created by Server.Listen() doesn't match the number of
	// listeners set in Config.Listeners.
	ErrMissingListener = errors.New("server: missing listener")

	// ErrAlreadyServing is returned by Server.Serve() when there are already
	// servers listening.  Before calling Server.Serve() again, Server.Close()
	// should be called to ensure all servers and listeners are closed.
	ErrAlreadyServing = errors.New("server: already serving")
)

var defaultTLSConfig = &tls.Config{
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

// Server .
type Server struct {
	config *Config

	listeners []net.Listener
	servers   []*http.Server

	newRouter func(l listener.Listener) *chi.Mux
}

// New .
func New(ops ...Opt) (*Server, error) {
	s := &Server{
		config: &Config{},

		newRouter: func(l listener.Listener) *chi.Mux {
			r := chi.NewRouter()
			r.Use(loggingAndRecovery)
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			if l.Metrics != "" {
				r.Get(l.Metrics, func(w http.ResponseWriter, r *http.Request) {
					metrics.WritePrometheus(w, true)
				})
			}
			return r
		},
	}
	for _, op := range ops {
		if err := op(s); err != nil {
			return nil, err
		}
	}
	return s, nil
}

// Listen .
func (s *Server) Listen(ctx context.Context) []error {
	var errs []error
	for _, lc := range s.config.Listeners {
		var lc2 net.ListenConfig
		lc2.KeepAlive = lc.KeepAlive
		l, err := lc2.Listen(ctx, lc.Network.String(), lc.Address)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		cosmos.Log().Info("listening on " + l.Addr().String())
		s.listeners = append(s.listeners, l)
	}
	return errs
}

// Serve .
func (s *Server) Serve(ctx context.Context) error {
	if s.listeners == nil || len(s.listeners) < 1 {
		return ErrNoListeners
	}
	if len(s.listeners) != len(s.config.Listeners) {
		return ErrMissingListener
	}
	if s.servers != nil && len(s.servers) > 0 {
		return ErrAlreadyServing
	}
	defer func() {
		s.servers = nil
		s.listeners = nil
	}()

	stdLog := zap.NewStdLog(cosmos.Log())
	g, ctx := errgroup.WithContext(ctx)
	for i, lc := range s.config.Listeners {
		l := s.listeners[i]
		addr := l.Addr().String()
		hs := &http.Server{
			Addr:    addr,
			Handler: s.newRouter(lc),

			TLSConfig: nil,

			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 3 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       30 * time.Second,

			MaxHeaderBytes: http.DefaultMaxHeaderBytes,

			ErrorLog: stdLog,
		}
		s.servers = append(s.servers, hs)
		cosmos.Log().Info("serving on " + addr)

		if certPath, keyPath := lc.CertPath, lc.KeyPath; certPath == "" || keyPath == "" {
			g.Go(func() error {
				if err := hs.Serve(l); err != nil &&
					err != http.ErrServerClosed &&
					!strings.HasSuffix(err.Error(), " use of closed network connection") {
					return err
				}
				return nil
			})
		} else {
			hs.TLSConfig = defaultTLSConfig
			g.Go(func() error {
				if err := hs.ServeTLS(l, certPath, keyPath); err != nil &&
					err != http.ErrServerClosed &&
					!strings.HasSuffix(err.Error(), " use of closed network connection") {
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

func loggingAndRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		defer func() {
			if err := recover(); err != nil {
				// I swear this log call will panic while we try
				// to recover from a panic.
				cosmos.Log().Error(
					"recovered from panic in http#Handler",
					zap.String("error", err.(string)),
				)
				// TODO: Write http#InternalServerError
			}

			remoteAddr, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				panic(err)
				return
			}

			var route string
			if ctx := chi.RouteContext(r.Context()); ctx != nil {
				if route = ctx.RoutePattern(); route != "/" {
					route = strings.TrimSuffix(route, "/")
				}
			} else {
				route = r.URL.Path
			}
			if route == "" {
				return
			}

			method := r.Method
			code := 200
			duration := time.Since(start)

			if code != 404 {
				metrics.RequestsTotal(method, route, code).Inc()
				metrics.RequestDuration(route).Update(duration.Seconds())
			}

			cosmos.Log().Info(
				"handled request",
				zap.String("remote", remoteAddr),
				zap.String("method", method),
				zap.String("route", route),
				zap.Int("code", code),
				zap.Duration("duration", duration.Round(time.Microsecond)),
			)
		}()

		w.Header().Set("Server", "Cosmos")
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, r)
	})
}
