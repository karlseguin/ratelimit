package ratelimit

import (
	"github.com/karlseguin/gspec"
	"testing"
	"time"
)

func Test_RateLimit_OkForNewItem(t *testing.T) {
	spec := gspec.New(t)
	limiter := New(Configure().Allowance(3))
	spec.Expect(limiter.Track("test")).ToEqual(OK)
}

func Test_RateLimit_OkForExistingItemOverAllowance(t *testing.T) {
	spec := gspec.New(t)
	limiter := New(Configure().Allowance(3))
	spec.Expect(limiter.Track("test")).ToEqual(OK)
	spec.Expect(limiter.Track("test")).ToEqual(OK)
	spec.Expect(limiter.Track("test")).ToEqual(OK)
}

func Test_RateLimit_FloodWhenPastThreshold(t *testing.T) {
	spec := gspec.New(t)
	limiter := New(Configure().
		Allowance(2).
		TTL(time.Minute))
	spec.Expect(limiter.Track("test")).ToEqual(OK)
	spec.Expect(limiter.Track("test")).ToEqual(OK)
	spec.Expect(limiter.Track("test")).ToEqual(FLOOD)
	spec.Expect(limiter.Track("other")).ToEqual(OK)
	spec.Expect(limiter.Track("other")).ToEqual(OK)
}
