package zbt

import (
	"testing"
	"time"
)

func Test_InitialTick(t *testing.T) {
	ticker := NewTicker(1 * time.Second)
	time.Sleep(500 * time.Millisecond)
	select {
	case <-ticker.C:
		return
	default:
		t.Fatal("unable to get first tick!")
	}
}

func Test_ReTick(t *testing.T) {
	ticker := NewTicker(250 * time.Millisecond)
	time.Sleep(50 * time.Millisecond)
	count := 0
	cutoff := time.After(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			count++
			if count > 1 {
				t.Log("reached count of 2 before 1 second, PASS")
				return
			}
		case <-cutoff:
			t.Fatal("cutoff reached before two ticks!")
		}
	}
}

func Test_Stop(t *testing.T) {
	ticker := NewTicker(time.Millisecond)
	ticker.Stop()
	safety := time.After(50 * time.Millisecond)
	select {
	case <-ticker.ctx.Done():
		t.Log("PASS received done; internal function stopped")
	case <-safety:
		t.Fatal("didn't get Done() closed after 50ms")
	}
}
