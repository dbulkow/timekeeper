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
	"sync"
	"time"
)

type TimeRunner interface {
	Run(context.Context, ...interface{})
}

type Timer struct {
	Action  TimeRunner
	Args    []interface{}
	Expires time.Time
}

type TimerSet struct {
	Timers []*Timer
	sync.Mutex
}

func NewTimerSet() *TimerSet {
	return &TimerSet{Timers: make([]*Timer, 0)}
}

func (ts *TimerSet) Tick(ctx context.Context) bool {
	ts.Lock()
	defer ts.Unlock()

	for i, t := range ts.Timers {
		now := time.Now()
		if now.After(t.Expires) {
			ts.Timers = ts.Timers[:i+copy(ts.Timers[i:], ts.Timers[i+1:])]
			if t.Action != nil {
				t.Action.Run(ctx, t.Args)
				return true
			}
		}
	}

	return false
}

func (ts *TimerSet) Remove(r *Timer) {
	ts.Lock()
	defer ts.Unlock()

	for i, t := range ts.Timers {
		if t == r {
			ts.Timers = ts.Timers[:i+copy(ts.Timers[i:], ts.Timers[i+1:])]
			return
		}
	}
}

func (ts *TimerSet) After(dur time.Duration, action TimeRunner, args ...interface{}) *Timer {
	return ts.Add(time.Now().Add(dur), action, args...)
}

func (ts *TimerSet) Add(expires time.Time, action TimeRunner, args ...interface{}) (t *Timer) {
	t = &Timer{
		Action:  action,
		Args:    args,
		Expires: expires,
	}

	ts.Lock()
	defer ts.Unlock()

	for i, x := range ts.Timers {
		if t.Expires.Before(x.Expires) {
			ts.Timers = append(ts.Timers[:i], append([]*Timer{t}, ts.Timers[i:]...)...)
			return
		}
	}

	ts.Timers = append(ts.Timers, t)

	return
}

func (ts *TimerSet) Find(t *Timer) bool {
	ts.Lock()
	defer ts.Unlock()

	for _, v := range ts.Timers {
		if v == t {
			return true
		}
	}

	return false
}
