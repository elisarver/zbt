package zbt

import (
	"context"
	"time"
)

// Ticker ticks as soon as its created, then waits until duration to
// tick again.
type Ticker struct {
	// C returns the time when a tick fires.
	C <-chan time.Time

	it         *time.Ticker
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewTicker creates a Ticker that comes pre-loaded with the first tick.
// duration must be non-zero or the underlying library will panic.
func NewTicker(duration time.Duration) *Ticker {
	// making c here so we can map it to a receive-only channel
	// in the Ticker
	c := make(chan time.Time, 1)
	zbt := Ticker{
		C: c,
	}

	zbt.ctx, zbt.cancelFunc = context.WithCancel(context.Background())

	// pre-seed c with an immediate value
	c <- time.Now()
	// create this right before the goroutine in case duration is
	// really small:
	zbt.it = time.NewTicker(duration)

	go func() {
		// We have to select against both done and the channel
		// in order to exit the goroutine cleanly on Close()
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

// Stop passes the stop command onto the internal ticker and
// closes the internal context.Context so the goroutine doesn't
// run forever.
func (zbt *Ticker) Stop() {
	zbt.it.Stop()
	zbt.cancelFunc()
	// drain the ticker to prevent accidental ticks.
	select {
	case _, _ = <-zbt.C:
	}
}
