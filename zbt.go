package zbt

import (
	"time"
)

// Ticker ticks as soon as its created, then waits until duration to
// tick again. Can be canceled by ctx.
type Ticker struct {
	C     chan time.Time
	it    *time.Ticker
	alive bool
}

// NewTicker creates a Ticker that immediately fires into its channel.
// duration must be non-zero or the underlying library will panic.
func NewTicker(duration time.Duration) *Ticker {
	zbt := &Ticker{
		C:  make(chan time.Time, 1),
		it: time.NewTicker(duration),
	}
	go func() {
		zbt.C <- time.Now()
		for zbt.alive = true; zbt.alive; {
			select {
			case now := <-zbt.it.C:
				zbt.C <- now
			}
		}
	}()
	return zbt
}

// Stop passes the stop command onto the internal ticker.
func (zbt *Ticker) Stop() {
	zbt.alive = false
	zbt.it.Stop()
}
