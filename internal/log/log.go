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

// Package log ...
package log

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var bufferPool = buffer.NewPool()

// Logger .
type Logger struct {
	Config *Config

	logger *zap.Logger
}

// New .
func New(ops ...Opt) (*Logger, error) {
	l := &Logger{
		Config: &Config{
			Level: InfoLevel,
		},
	}
	for _, op := range ops {
		if err := op(l); err != nil {
			return nil, err
		}
	}
	return l, nil
}

// Production .
func (l *Logger) Production() (*zap.Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("15:04:05"))
	}
	config.EncodeDuration = zapcore.StringDurationEncoder
	config.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		if !caller.Defined {
			enc.AppendString("undefined")
			return
		}

		idx := strings.LastIndexByte(caller.File, '/')
		if idx == -1 {
			enc.AppendString(caller.FullPath())
			return
		}

		buf := bufferPool.Get()

		buf.AppendString(caller.File[idx+1:])
		buf.AppendByte(':')
		buf.AppendInt(int64(caller.Line))
		c := buf.String()

		buf.Free()

		enc.AppendString(c)
	}

	writeSyncer, _, err := zap.Open("stderr")
	if err != nil {
		return nil, errors.Wrap(err, "failed to open stderr")
	}

	return zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(config),
			writeSyncer,
			zapcore.Level(l.Config.Level),
		),
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	), nil
}

// SetGlobal .
func SetGlobal(l *zap.Logger) {
	zap.ReplaceGlobals(l)
}
