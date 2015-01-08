package ratelimit

type Configuration struct {
	maxItems  int64
	allowance int
}

func Configure() *Configuration {
	return &Configuration{
		maxItems:  5000,
		allowance: 5,
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
