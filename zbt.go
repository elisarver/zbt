package zbt

import (
	"context"
	"time"
)

// Ticker ticks as soon as its created, then waits until duration to
// tick again. Can be canceled by ctx.
type Ticker struct {
	C  chan time.Time
	it *time.Ticker
}

// NewTicker creates a Ticker that immediately fires into its channel.
// ctx, if canceled, will immediately release the ticker.
// duration must be non-zero or the underlying library will panic.
func NewTicker(ctx context.Context, duration time.Duration) *Ticker {
	zbt := &Ticker{
		C:  make(chan time.Time, 1),
		it: time.NewTicker(duration),
	}
	go func() {
		zbt.C <- time.Now()
		for alive := true; alive; {
			select {
			case <-ctx.Done():
				// prevent leaks, and exit.
				zbt.it.Stop()
				alive = false
			case now := <-zbt.it.C:
				// pass ticks along
				zbt.C <- now
			}
		}
	}()
	return zbt
}

// Stop passes the stop command onto the internal ticker.
func (zbt *Ticker) Stop() {
	zbt.it.Stop()
}
