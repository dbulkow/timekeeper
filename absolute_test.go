/*
MIT License

Copyright (c) 2020 David Bulkow

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package timekeeper

import (
	"context"
	"testing"
	"time"
)

func TestAbsoluteAdjust(t *testing.T) {
	a := &Absolute{}
	now := time.Now().Local()

	for hour := now.Hour() - 5; hour < now.Hour()+5; hour++ {
		tm := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, time.Local)
		if hour/24 > 0 {
			tm.AddDate(0, 0, 1)
		}
		day := tm.Day()
		if time.Until(tm) < 0 {
			day++
		}
		tm = a.adjust(tm)
		if tm.Day() != day {
			t.Fatalf("expected day to be %d got %d", day, tm.Day())
		}
	}
}

func TestAbsoluteExcludes(t *testing.T) {
	tests := []struct {
		name string
		excl []time.Weekday
	}{
		{"weekend", []time.Weekday{time.Saturday, time.Sunday}},
		{"mon,tue", []time.Weekday{time.Monday, time.Tuesday}},
		{"tue,thu", []time.Weekday{time.Tuesday, time.Thursday}},
	}

	inexcl := func(tm time.Time, excl []time.Weekday) bool {
		weekday := tm.Weekday()
		for _, e := range excl {
			if weekday == e {
				return true
			}
		}
		return false
	}

	for _, tr := range tests {
		now := time.Now()
		t.Run(tr.name, func(t *testing.T) {
			a := &Absolute{Excludes: tr.excl}
			when := time.Date(now.Year(), now.Month(), now.Day(), 5, 30, 0, 0, time.Local)
			for i := 0; i < 7; i++ {
				tm := when.AddDate(0, 0, i)
				ta := a.adjust(tm)
				if inexcl(tm, tr.excl) && tm.Day() > ta.Day() {
					t.Fatalf("expected delay: before %s after %s (%s)", tm, ta, tm.Weekday())
				}
			}
		})
	}
}

type testabs struct{}

func (t *testabs) Run(context.Context, ...interface{}) {}

func TestAbsoluteTrigger(t *testing.T) {
	ts := NewTimerSet()
	a := NewAbsoluteEvent(ts, 14, 25, []time.Weekday{time.Sunday, time.Saturday})
	a.Trigger(&testabs{})
	var count int
	for range ts.Timers {
		count++
	}
	if count != 1 {
		t.Fatalf("expected 1 timer scheduled, found %d", count)
	}
}

func TestAbsoluteRetrigger(t *testing.T) {
	ts := NewTimerSet()
	a := NewAbsoluteEvent(ts, 14, 25, []time.Weekday{time.Sunday, time.Saturday})
	a.Trigger(&testabs{})
	a.Trigger(&testabs{})
	var count int
	for range ts.Timers {
		count++
	}
	if count != 1 {
		t.Fatalf("expected 1 timer scheduled, found %d", count)
	}
}
