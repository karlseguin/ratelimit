package ratelimit

import (
	"time"
)

type RateLimit struct {
	*Cache
	allowance int32
	ttl       time.Duration
}

func New(config *Configuration) *RateLimit {
	return &RateLimit{
		Cache:     NewCache(config.maxItems),
		allowance: int32(config.allowance),
	}
}

func (r *RateLimit) Track(key string) int32 {
	tracker := r.Fetch(key)
	return tracker.Track(r.allowance)
}
