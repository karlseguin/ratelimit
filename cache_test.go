package ratelimit

import (
	. "github.com/karlseguin/expect"
	"strconv"
	"testing"
	"time"
)

type CacheTests struct{}

func Test_Cache(t *testing.T) {
	Expectify(new(CacheTests), t)
}

func (_ CacheTests) FetchesANewItem() {
	cache := NewCache(10)
	tracker := cache.Fetch("leto")
	Expect(tracker.Track(4)).To.Equal(int32(3))
}

func (_ CacheTests) FetchesAnExitingItem() {
	cache := NewCache(10)
	tracker := cache.Fetch("paul")
	Expect(tracker.Track(4)).To.Equal(int32(3))

	tracker = cache.Fetch("paul")
	Expect(tracker.Track(4)).To.Equal(int32(2))
}

func (_ CacheTests) PurgesWhenOverLimit() {
	cache := NewCache(10)
	for i := 0; i < 12; i++ {
		cache.Fetch(strconv.Itoa(i)).Track(4)
	}
	time.Sleep(time.Millisecond * 10) //let the gc run
	Expect(cache.Fetch("0").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("1").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("2").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("3").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("4").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("5").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("6").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("7").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("8").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("9").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("10").Track(4)).To.Equal(int32(2))
	Expect(cache.Fetch("11").Track(4)).To.Equal(int32(2))
}

func (_ CacheTests) StopsPurgingWhenNothingLeftToPurge() {
	cache := NewCache(10)
	for i := 0; i < 4; i++ {
		cache.Fetch(strconv.Itoa(i)).Track(5)
	}
	time.Sleep(time.Millisecond * 10) //let the gc run
	Expect(cache.Fetch("0").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("1").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("2").Track(4)).To.Equal(int32(3))
	Expect(cache.Fetch("3").Track(4)).To.Equal(int32(3))
}
