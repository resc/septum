package septum

import (
	"time"
)

type (
	// Environment abstracts the environment functions like the current time, and id generators.
	Environment interface {

		// Now returns the current time (in UTC)
		Now() time.Time

		// NextEventId generates the next event id
		NextEventId() uint64
	}

	// Nexus brings all timelines together
	Nexus interface {

		// PeekNextTimer returns the next timer to expire, or nil, false if there are no timers
		// it doesn't get removed from the timeline
		PeekNextTimer() (*Event, bool)

		// RemoveExpiredTimer removes the oldest expired event from the timer queue and returns the event and true or nil, false if there are no timers.
		RemoveExpiredTimer(now time.Time) (*Event, bool)

		// GetTimeline retrieves the timeline for the given id, there is a timeline for every id you can think of.
		GetTimeline(timelineId uint64) Timeline

		// Add adds a new event to the indicated timeline
		AddEvent(timelineId uint64, timestamp time.Time, kind string, data []byte) *Event

		// AddTimer adds a new event to the indicated timeline and to the timer queue
		AddTimer(timelineId uint64, timestamp time.Time, kind string, data []byte) *Event

		// Remove removes the event from its timeline and the timer queue, returns true if the event was removed, false if it doesn't exist in the timelines
		Remove(e *Event) bool

		// RemoveRange removes the events from their timeline and the timer queue, returns the number of events removed
		RemoveRange(r Ranger) int
	}

	// Timeline is a event list ordered by detect time
	Timeline interface {
		// Id returns the timeline id
		Id() uint64
		// Range returns all the events in the timeline in chronological order
		Range() Range
		// The number of events in the timeline
		Len() int
	}

	// Ranger can return a Range of events
	Ranger interface {
		Range() Range
	}

	// Range is an event iterator
	Range interface {
		// Next moves the cursor to the next event
		Next() bool
		// Event retrieves the current event, only valid after a call to Next that returns true.
		Event() *Event
	}
)
