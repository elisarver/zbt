package zbt

import (
	"context"
	"time"
)

// Ticker ticks as soon as its created, then waits until duration to
// tick again.
type Ticker struct {
	C <-chan time.Time

	it         *time.Ticker
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewTicker creates a Ticker that immediately fires into its channel.
// duration must be non-zero or the underlying library will panic.
func NewTicker(duration time.Duration) *Ticker {
	ctx, done := context.WithCancel(context.Background())
	c := make(chan time.Time, 1)
	zbt := Ticker{
		C:          c,
		ctx:        ctx,
		cancelFunc: done,
	}
	c <- time.Now()
	zbt.it = time.NewTicker(duration)

	go func() {
		// We have to select against both done and
		for {
			select {
			case <-zbt.ctx.Done():
				return
			case now := <-zbt.it.C:
				c <- now
			}
		}
	}()

	return &zbt
}

// Stop passes the stop command onto the internal ticker.
func (zbt *Ticker) Stop() {
	zbt.it.Stop()
	zbt.cancelFunc()
	// drain
	select {
	case _, _ = <-zbt.C:
	}
}
