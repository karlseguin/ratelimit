package ratelimit

import (
	"sync/atomic"
	"time"
)

type Tracker struct {
	Allowance int32
	LastRead  int64
}

func NewTracker() *Tracker {
	return new(Tracker)
}

func (t *Tracker) Track(allowedPerSecond int32) int32 {
	now := time.Now().Unix()
	earned := int32(now - atomic.SwapInt64(&t.LastRead, now))
	if earned > allowedPerSecond {
		earned = allowedPerSecond
	}
	allowance := atomic.AddInt32(&t.Allowance, earned-1)
	if allowance > allowedPerSecond {
		allowance = allowedPerSecond
		atomic.StoreInt32(&t.Allowance, allowedPerSecond)
	}
	return allowance
}
