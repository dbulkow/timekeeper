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
	"fmt"
	"log"
	"testing"
	"time"
)

func TestTimerAdd(t *testing.T) {
	dur := []time.Duration{
		time.Second,
		15 * time.Second,
		7 * time.Second,
		30 * time.Second,
		31 * time.Second,
		0,
	}

	ts := NewTimerSet()

	for _, d := range dur {
		ts.After(d, nil)
	}

	for _, t := range ts.Timers {
		fmt.Println(t.Expires)
	}
}

type ticktest struct {
	count *int
}

func (t *ticktest) Run(context.Context, ...interface{}) {
	*t.count -= 1
	fmt.Println("run", time.Now(), *t.count)
}

func TestTimerTick(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	dur := []time.Duration{
		time.Second,
		5 * time.Second,
		7 * time.Second,
		8 * time.Second,
		3 * time.Second,
		5 * time.Second,
	}

	ts := NewTimerSet()

	var count int

	for _, d := range dur {
		ts.After(d, &ticktest{&count})
	}

	for _, t := range ts.Timers {
		fmt.Println(t.Expires)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for count = len(dur); count > 0; {
		time.Sleep(500 * time.Millisecond)
		// fmt.Println("here", time.Now(), count)
		ts.Tick(ctx)
	}
}

type ticktesttwo struct {
	count *int
}

func (t *ticktesttwo) Run(context.Context, ...interface{}) {
	*t.count -= 1
	log.Println("run", *t.count)
	time.Sleep(4 * time.Second)
}

func TestTimerSlowRun(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	dur := []time.Duration{
		time.Second,
		5 * time.Second,
		7 * time.Second,
		8 * time.Second,
		3 * time.Second,
		5 * time.Second,
	}

	ts := NewTimerSet()

	var count int

	for _, d := range dur {
		ts.After(d, &ticktesttwo{&count})
	}

	for _, t := range ts.Timers {
		fmt.Println(t.Expires)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var action bool
	start := time.Now()
	for count = len(dur); count > 0; {
		if !action {
			time.Sleep(500 * time.Millisecond)
		}
		log.Println("here", time.Since(start).Truncate(time.Second), count)
		action = ts.Tick(ctx)
	}
}

func TestTimerRemove(t *testing.T) {
	dur := []time.Duration{
		time.Second,
		15 * time.Second,
		7 * time.Second,
		30 * time.Second,
		31 * time.Second,
		0,
	}

	ts := NewTimerSet()

	timers := make([]*Timer, 0)

	for _, d := range dur {
		timers = append(timers, ts.After(d, nil))
	}

	ts.Remove(timers[2])

	for _, t := range ts.Timers {
		fmt.Println(t.Expires)
	}
}
