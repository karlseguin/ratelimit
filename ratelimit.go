package ratelimit

import (
	"github.com/karlseguin/ccache"
	"time"
)

type RateLimit struct {
	*ccache.Cache
	allowance int32
	ttl       time.Duration
}

func New(config *Configuration) *RateLimit {
	return &RateLimit{
		Cache:     ccache.New(ccache.Configure().MaxItems(config.maxItems).GetsPerPromote(1)),
		allowance: int32(config.allowance),
		ttl:       config.ttl,
	}
}

func (r *RateLimit) Track(key string) int32 {
	tracker, _ := r.Fetch(key, r.ttl, func() (interface{}, error) { return NewTracker(), nil })
	return tracker.(*Tracker).Track(r.allowance)
}
