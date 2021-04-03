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

// Package backoff ...
package backoff

import (
	"context"
	"math"
	"time"
)

// maxInt64 is used to avoid overflowing a time.Duration (int64) value.
const maxInt64 = float64(math.MaxInt64 - 512)

// Backoff represents an exponential backoff.
type Backoff struct {
	// n is the current attempt and defaults to 0.  The first attempt will not
	// have any delay before it is ran.
	n uint

	// MaxAttempts is the max number of attempts that can occur.
	MaxAttempts uint
	// Factor is the factor at which Min will increase after each failed attempt.
	Factor float64
	// Min is the initial backoff time to wait after the first failed attempt.
	Min time.Duration
	// Max is the maximum time to wait before retrying.
	Max time.Duration

	// NewTimer is used for unit tests.  For actual use, this should be set to
	// time.NewTimer.
	NewTimer func(time.Duration) *time.Timer
}

// New returns a new Backoff instance.
func New(maxAttempts uint, factor float64, min, max time.Duration) *Backoff {
	return &Backoff{
		n: 0,

		MaxAttempts: maxAttempts,
		Factor:      factor,
		Min:         min,
		Max:         max,

		NewTimer: time.NewTimer,
	}
}

// Attempt returns the current attempt.
func (b *Backoff) Attempt() uint {
	return b.n
}

// Duration returns a time.Duration to wait for a specified attempt.
func (b *Backoff) Duration(attempt float64) time.Duration {
	if attempt == 0 {
		return 0
	}

	durF := float64(b.Min) * math.Pow(b.Factor, attempt)
	if durF > maxInt64 {
		return b.Max
	}

	dur := time.Duration(durF)
	if dur < b.Min {
		return b.Min
	}
	if dur > b.Max {
		return b.Max
	}
	return dur
}

// Next increments the attempt, then waits for the Duration of the attempt to
// pass, returning true.  Next will return false if the attempt will exceed the
// MaxAttempts limit or if the context has been cancelled.
func (b *Backoff) Next(ctx context.Context) bool {
	if b.n >= b.MaxAttempts {
		return false
	}
	d := b.Duration(float64(b.n))
	b.n++

	t := b.NewTimer(d)
	select {
	case <-ctx.Done():
		t.Stop()
		return false
	case <-t.C:
		return true
	}
}
