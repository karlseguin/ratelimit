package ratelimit

import (
	. "github.com/karlseguin/expect"
	"testing"
	"time"
)

type RateLimitTests struct{}
func Test_RateLimit(t *testing.T) {
	Expectify(new(RateLimitTests), t)
}

func (e *RateLimitTests) OkForNewItem() {
	limiter := New(Configure().Allowance(3))
	Expect(limiter.Track("test")).GreaterOrEqual.To(int32(0))
}

func (e *RateLimitTests) OkForExistingItemOverAllowance() {
	limiter := New(Configure().Allowance(3))
	Expect(limiter.Track("test")).GreaterOrEqual.To(int32(0))
	Expect(limiter.Track("test")).GreaterOrEqual.To(int32(0))
	Expect(limiter.Track("test")).GreaterOrEqual.To(int32(0))
}

func (e *RateLimitTests) FloodWhenPastThreshold() {
	limiter := New(Configure().
		Allowance(2).
		TTL(time.Minute))
	Expect(limiter.Track("test")).GreaterOrEqual.To(int32(0))
	Expect(limiter.Track("test")).GreaterOrEqual.To(int32(0))
	Expect(limiter.Track("test")).Less.Than(int32(0))
	Expect(limiter.Track("other")).GreaterOrEqual.To(int32(0))
	Expect(limiter.Track("other")).GreaterOrEqual.To(int32(0))
}
