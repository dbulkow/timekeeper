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
)

type testimm struct{}

func (t *testimm) Run(context.Context, ...interface{}) {}

func TestImmediateEvent(t *testing.T) {
	ts := NewTimerSet()
	d := NewImmediateEvent(ts)
	d.Trigger(&testimm{})
	var count int
	for range ts.Timers {
		count++
	}
	if count != 1 {
		t.Fatalf("expected 1 timer scheduled, found %d", count)
	}
}

func TestImmediateRetrigger(t *testing.T) {
	ts := NewTimerSet()
	d := NewImmediateEvent(ts)
	d.Trigger(&testimm{})
	d.Trigger(&testimm{})
	var count int
	for range ts.Timers {
		count++
	}
	if count != 2 {
		t.Fatalf("expected 1 timer scheduled, found %d", count)
	}
}
