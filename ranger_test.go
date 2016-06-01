// septum_test.go
package septum

import (
	"testing"
)

func TestEmptySliceRange(t *testing.T) {
	events := make([]*Event, 0, 1)
	r := ToRange(events)

	count := 0
	for i := r.Range(); i.Next(); {
		count = count + 1
	}

	if count != 0 {
		t.Errorf("Expected count to be zero, but got %d", count)
	}

}

func TestSingleEntrySliceRange(t *testing.T) {
	events := make([]*Event, 1, 1)
	r := ToRange(events)
	count := 0
	for i := r.Range(); i.Next(); {
		count = count + 1
	}

	if count != 1 {
		t.Errorf("Expected count to be 1, but got %d", count)
	}
}

func TestWhereSliceRangeWithNils(t *testing.T) {
	events := make([]*Event, 1, 1)
	r := Where(ToRange(events), func(e *Event) bool { return e != nil })
	count := 0
	for i := r.Range(); i.Next(); {
		count = count + 1
	}

	if count != 0 {
		t.Errorf("Expected count to be 1, but got %d", count)
	}
}

func TestWhereSliceRangeNonNils(t *testing.T) {
	events := make([]*Event, 3, 3)
	events[0] = &Event{}
	events[1] = &Event{}
	r := Where(ToRange(events), func(e *Event) bool { return e != nil })
	count := 0
	for i := r.Range(); i.Next(); {
		count = count + 1
	}

	if count != 2 {
		t.Errorf("Expected count to be 2, but got %d", count)
	}
}

func expectPanic(t *testing.T, f func()) {
	defer func() {
		var err = recover()
		if err == nil {
			t.Error("Expected a panic")
		}
	}()
	f()
}
