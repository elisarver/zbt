package zbt

import (
	"testing"
	"time"
)

func TestNewTicker_InitialTick(t *testing.T) {
	ticker := NewTicker(1 * time.Second)
	time.Sleep(500 * time.Millisecond)
	select {
	case <-ticker.C:
		return
	default:
		t.Fatal("unable to get first tick!")
	}
}

func TestNewTicker_ReTick(t *testing.T) {
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

func TestNewTicker_StopBeforeRetick(t *testing.T) {
	ticker := NewTicker(250 * time.Millisecond)
	time.Sleep(50 * time.Millisecond)
	count := 0
	cutoff := time.After(500 * time.Millisecond)
	for alive := true; alive; {
		select {
		case <-ticker.C:
			ticker.Stop()
			count++
			if count > 1 {
				t.Fatalf("reached count of 2, stop unsuccessful")
				return
			}
		case <-cutoff:
			alive = false
			t.Log("reached cutoff without multiple ticks, PASS")
		}
	}
	ticker.Stop()
}
