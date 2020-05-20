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
	"sync"
	"time"
)

type When struct {
	Hour   int
	Minute int
}

type Absolute struct {
	TimerSet   *TimerSet
	When       *When          // wall clock time to execute event
	Timer      *Timer         // timer we scheduled for execution
	Excludes   []time.Weekday // range of excluded days
	sync.Mutex                // single thread to protect timer lookup
}

func NewAbsoluteEvent(ts *TimerSet, hour, min int, excludes []time.Weekday) Event {
	return &Absolute{
		TimerSet: ts,
		When:     &When{Hour: hour, Minute: min},
		Excludes: excludes,
	}
}

func (a *Absolute) adjust(when time.Time) time.Time {
	if time.Until(when) < 0 {
		when = when.AddDate(0, 0, 1)
	}

	for i := 0; i < len(a.Excludes); i++ {
		if when.Weekday() == a.Excludes[i] {
			when = when.AddDate(0, 0, 1)
			i = 0 // start over to deal with out of order Excludes
		}
	}

	return when
}

func (a *Absolute) Trigger(action TimeRunner, args ...interface{}) *Timer {
	a.Lock()
	defer a.Unlock()

	if a.TimerSet.Find(a.Timer) {
		return a.Timer
	}

	now := time.Now().Local()
	when := time.Date(now.Year(), now.Month(), now.Day(), a.When.Hour, a.When.Minute, 0, 0, time.Local)

	a.Timer = a.TimerSet.Add(a.adjust(when), action, args)

	return a.Timer
}
