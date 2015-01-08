package ratelimit

import (
	. "github.com/karlseguin/expect"
	"testing"
)

type RateLimitTests struct{}

func Test_RateLimit(t *testing.T) {
	Expectify(new(RateLimitTests), t)
}

func (_ RateLimitTests) OkForNewItem() {
	limiter := New(Configure().Allowance(3))
	Expect(limiter.Track("test")).GreaterOrEqual.To(int32(0))
}

func (_ RateLimitTests) OkForExistingItemOverAllowance() {
	limiter := New(Configure().Allowance(3))
	Expect(limiter.Track("test")).GreaterOrEqual.To(int32(0))
	Expect(limiter.Track("test")).GreaterOrEqual.To(int32(0))
	Expect(limiter.Track("test")).GreaterOrEqual.To(int32(0))
}

func (_ RateLimitTests) FloodWhenPastThreshold() {
	limiter := New(Configure().Allowance(2))
	Expect(limiter.Track("test")).GreaterOrEqual.To(int32(0))
	Expect(limiter.Track("test")).GreaterOrEqual.To(int32(0))
	Expect(limiter.Track("test")).Less.Than(int32(0))
	Expect(limiter.Track("other")).GreaterOrEqual.To(int32(0))
	Expect(limiter.Track("other")).GreaterOrEqual.To(int32(0))
}
