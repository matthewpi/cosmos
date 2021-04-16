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

package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/matthewpi/cosmos/internal/config"
	"go.uber.org/zap"

	"github.com/matthewpi/cosmos"
	"github.com/matthewpi/cosmos/internal/listener"
	"github.com/matthewpi/cosmos/internal/log"
	"github.com/matthewpi/cosmos/internal/server"
)

func main() {
	cfg, err := config.Load(".env/cosmos.conf")
	if err != nil {
		panic(err)
		return
	}

	var opts []log.Opt
	for _, s := range cfg.Key("log").Segments {
		d := s.Directive()
		switch d {
		case "output":
		case "level":
			if len(s) != 2 {
				return
			}
			l, ok := log.Levels[strings.ToLower(s[1].Text)]
			if !ok {
				return
			}
			opts = append(opts, log.WithLevel(l))
		}
	}

	l, err := log.New(opts...)
	if err != nil {
		panic(err)
		return
	}

	productionLogger, err := l.Production()
	if err != nil {
		panic(err)
		return
	}
	log.SetGlobal(productionLogger)
	defer cosmos.Log().Sync()

	cosmos.Log().Info("Hello, world!")

	s, err := server.New(
		server.WithListener(listener.Listener{
			Network: listener.NetworkTCP,
			Address: ":9000",
		}),
	)
	if err != nil {
		cosmos.Log().Fatal("failed to create new server", zap.Error(err))
		return
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.Close(ctx); err != nil {
			cosmos.Log().Error("failed to close server", zap.Error(err))
			return
		}
	}()

	cosmos.Log().Info("attempting to start listening...")
	if errs := s.Listen(context.Background()); errs != nil {
		var fields []zap.Field
		for _, err := range errs {
			fields = append(fields, zap.Error(err))
		}
		cosmos.Log().Error("failed to start listening", fields...)
		return
	}

	go func() {
		cosmos.Log().Info("attempting to start http servers...")
		if err := s.Serve(context.Background()); err != nil {
			cosmos.Log().Error("failed to start serving", zap.Error(err))
			return
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Close(ctx); err != nil {
		cosmos.Log().Error("failed to close server", zap.Error(err))
		return
	}
}
