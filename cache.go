package ratelimit

import (
	"hash/fnv"
	"sync"
)

const (
	BUCKETS     = 16
	BUCKET_MASK = BUCKETS - 1
)

type CachedTracker struct {
	key     string
	tracker *Tracker
	prev    *CachedTracker
	next    *CachedTracker
}

type Cache struct {
	list        *List
	maxItems    int64
	items       int64
	purgeSize   int64
	buckets     []*bucket
	promotables chan *CachedTracker
}

func NewCache(maxItems int64) *Cache {
	c := &Cache{
		maxItems:    maxItems,
		list:        NewList(),
		purgeSize:   maxItems / 20,
		buckets:     make([]*bucket, BUCKETS),
		promotables: make(chan *CachedTracker, 1024),
	}
	for i := 0; i < BUCKETS; i++ {
		c.buckets[i] = &bucket{lookup: make(map[string]*CachedTracker)}
	}
	if c.purgeSize < 10 {
		c.purgeSize = 10
	}
	go c.worker()
	return c
}

func (c *Cache) Fetch(key string) *Tracker {
	bucket := c.bucket(key)
	item := bucket.fetch(key)
	c.promotables <- item
	return item.tracker
}

func (c *Cache) bucket(key string) *bucket {
	h := fnv.New32a()
	h.Write([]byte(key))
	return c.buckets[h.Sum32()&BUCKET_MASK]
}

func (c *Cache) worker() {
	for {
		item := <-c.promotables
		if item.prev == nil { //new item
			c.items += 1
			if c.items > c.maxItems {
				c.gc()
			}
		}
		c.list.PushToFront(item)
	}
}

func (c *Cache) gc() {
	item := c.list.tail.prev
	for i := int64(0); i < c.purgeSize; i++ {
		if item == c.list.head {
			break
		}
		c.bucket(item.key).remove(item.key)
		c.list.Remove(item)
		c.items -= 1
		item = item.prev
	}
}

type bucket struct {
	sync.RWMutex
	lookup map[string]*CachedTracker
}

func (b *bucket) fetch(key string) *CachedTracker {
	b.RLock()
	item := b.lookup[key]
	b.RUnlock()
	if item != nil {
		return item
	}

	// we might end up not needing this
	// but better to create it outside the lock
	// (since most of the time we will need it)
	newItem := &CachedTracker{
		key:     key,
		tracker: NewTracker(),
	}

	b.Lock()
	defer b.Unlock()
	if item := b.lookup[key]; item != nil {
		return item
	}
	b.lookup[key] = newItem
	return newItem
}

func (b *bucket) remove(key string) {
	b.Lock()
	defer b.Unlock()
	delete(b.lookup, key)
}
