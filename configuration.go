package ratelimit

import (
	"time"
)

type Configuration struct {
	maxItems  int64
	allowance int
	ttl       time.Duration
}

func Configure() *Configuration {
	return &Configuration{
		maxItems:  5000,
		allowance: 5,
		ttl:       time.Minute * 10,
	}
}

// The max number of items to track
// [5000]
func (c *Configuration) MaxItems(max int) *Configuration {
	c.maxItems = int64(max)
	return c
}

// The number of events allowed per second
// [5]
func (c *Configuration) Allowance(allowance int) *Configuration {
	c.allowance = allowance
	return c
}

// The length of time to wait before purging an idle key
// [10 minutes]
func (c *Configuration) TTL(ttl time.Duration) *Configuration {
	c.ttl = ttl
	return c
}
