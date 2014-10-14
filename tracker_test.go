package ratelimit

import (
	"github.com/karlseguin/gspec"
	"testing"
	"time"
)

func Test_Tracker_OkWhenOverAllowance(t *testing.T) {
	spec := gspec.New(t)
	tracker := &Tracker{}
	spec.Expect(tracker.Track(5)).ToEqual(OK)
}

func Test_Tracker_OkWhenOverThreshold(t *testing.T) {
	spec := gspec.New(t)
	tracker := &Tracker{Allowance: 2, LastRead: time.Now().Unix()}
	spec.Expect(tracker.Track(5)).ToEqual(OK)
}

func Test_Tracker_FloodWhenBelowTreshold(t *testing.T) {
	spec := gspec.New(t)
	tracker := &Tracker{Allowance: 2, LastRead: time.Now().Unix()}
	spec.Expect(tracker.Track(5)).ToEqual(OK)
	spec.Expect(tracker.Track(5)).ToEqual(OK)
	spec.Expect(tracker.Track(5)).ToEqual(FLOOD)
	spec.Expect(tracker.Track(5)).ToEqual(FLOOD)
}

func Test_Tracker_OkAfterARestPeriod(t *testing.T) {
	spec := gspec.New(t)
	tracker := &Tracker{Allowance: 1, LastRead: time.Now().Unix()}
	spec.Expect(tracker.Track(5)).ToEqual(OK)
	spec.Expect(tracker.Track(5)).ToEqual(FLOOD)
	spec.Expect(tracker.Track(5)).ToEqual(FLOOD)
	time.Sleep(time.Second * 3)
	spec.Expect(tracker.Track(5)).ToEqual(OK)
	spec.Expect(tracker.Track(5)).ToEqual(FLOOD)
}
