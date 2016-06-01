package septum

// ==== ranger ====

type (
	ranger struct {
		start *node
		curr  *node
		event *Event
	}
)

func ToSlice(r Ranger) []*Event {
	events := make([]*Event, 0, 4)

	for x := r.Range(); x.Next(); {
		events = append(events, x.Event())
	}

	return events
}

func newRanger(n *node) *ranger {
	return &ranger{start: n, curr: n}
}

func (r *ranger) Range() Range {
	return newRanger(r.start)
}

func (r *ranger) Next() bool {
	if r.curr == nil {
		r.event = nil
		return false
	}

	r.event = r.curr.e
	r.curr = r.curr.next
	return true
}

func (r *ranger) Event() *Event {
	return r.event
}

// ==== whereRanger ===

type (
	whereRanger struct {
		inner  Ranger
		filter func(*Event) bool
	}

	whereRange struct {
		inner  Range
		filter func(*Event) bool
		event  *Event
	}
)

// Where creates a filtered range containing only the events for which filter returns true
func Where(source Ranger, filter func(*Event) bool) Ranger {
	return &whereRanger{
		inner:  source,
		filter: filter,
	}
}

func (r *whereRanger) Range() Range {
	return &whereRange{
		inner:  r.inner.Range(),
		filter: r.filter,
	}
}

func (r *whereRange) Next() bool {
	for r.inner.Next() {
		if r.filter(r.inner.Event()) {
			r.event = r.inner.Event()
			return true
		}
	}

	r.event = nil
	return false
}

func (r *whereRange) Event() *Event {
	return r.event
}

// ==== sliceRanger ====

func ToRange(events []*Event) Ranger {
	return newSliceRange(events)
}

func newSliceRange(events []*Event) *sliceRange {
	return &sliceRange{
		events: events,
		index:  -1,
	}
}

type sliceRange struct {
	index  int
	events []*Event
}

func (r *sliceRange) Range() Range {
	return newSliceRange(r.events)
}

func (r *sliceRange) Next() bool {
	if r.index < len(r.events) {
		r.index += 1
	}
	return r.index < len(r.events)
}

func (r *sliceRange) Event() *Event {
	if 0 > r.index || r.index >= len(r.events) {
		return nil
	}
	return r.events[r.index]
}
