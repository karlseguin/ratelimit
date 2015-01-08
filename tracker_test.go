package ratelimit

import (
	. "github.com/karlseguin/expect"
	"testing"
	"time"
)

type TrackerTests struct{}

func Test_Tracker(t *testing.T) {
	Expectify(new(TrackerTests), t)
}

func (_ TrackerTests) OkWhenOverAllowance() {
	tracker := &Tracker{}
	Expect(tracker.Track(5)).GreaterOrEqual.To(int32(0))
}

func (_ TrackerTests) OkWhenOverThreshold() {
	tracker := &Tracker{Allowance: 2, LastRead: time.Now().Unix()}
	Expect(tracker.Track(5)).GreaterOrEqual.To(int32(0))
}

func (_ TrackerTests) FloodWhenBelowTreshold() {
	tracker := &Tracker{Allowance: 2, LastRead: time.Now().Unix()}
	Expect(tracker.Track(5)).GreaterOrEqual.To(int32(0))
	Expect(tracker.Track(5)).GreaterOrEqual.To(int32(0))
	Expect(tracker.Track(5)).Less.Than(int32(0))
	Expect(tracker.Track(5)).Less.Than(int32(0))
}

func (_ TrackerTests) OkAfterARestPeriod() {
	tracker := &Tracker{Allowance: 1, LastRead: time.Now().Unix()}
	Expect(tracker.Track(5)).GreaterOrEqual.To(int32(0))
	Expect(tracker.Track(5)).Less.Than(int32(0))
	Expect(tracker.Track(5)).Less.Than(int32(0))
	time.Sleep(time.Second * 3)
	Expect(tracker.Track(5)).GreaterOrEqual.To(int32(0))
	Expect(tracker.Track(5)).Less.Than(int32(0))
}
