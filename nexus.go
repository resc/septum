package septum

import (
	"fmt"
	"time"
)

// ==== Nexus ====

type (
	nexus struct {
		environment Environment
		timelines   map[uint64]*timeline
		timers      *timeline
	}
)

// NewNexus creates a new Nexus for the given environment.
func NewNexus(env Environment) Nexus {
	if env == nil {
		panic("env is nil")
	}

	return &nexus{
		environment: env,
		timelines:   make(map[uint64]*timeline),
		timers:      newTimeline(0),
	}
}

func (m *nexus) PeekNextTimer() (*Event, bool) {
	return m.peekNextTimer()
}

func (m *nexus) peekNextTimer() (*Event, bool) {
	h := m.timers.head
	if h != nil {
		return h.e, true
	}
	return nil, false
}

func (m *nexus) RemoveExpiredTimer(now time.Time) (*Event, bool) {
	return m.removeExpiredTimer(now)
}

func (m *nexus) removeExpiredTimer(now time.Time) (*Event, bool) {
	h := m.timers.head
	if h != nil && now.After(h.e.Timestamp) {
		e := h.e
		m.timers.del(h)
		return e, true
	}

	return nil, false
}

func (m *nexus) GetTimeline(timelineId uint64) Timeline {
	if timelineId < 1 {
		panic(fmt.Errorf("Invalid timelineId: %d", timelineId))
	}

	return m.getTimeline(timelineId)
}

func (m *nexus) getTimeline(timelineId uint64) *timeline {
	s, ok := m.timelines[timelineId]
	if !ok {
		s = newTimeline(timelineId)
		m.timelines[timelineId] = s
	}
	return s
}

func (m *nexus) AddEvent(timelineId uint64, timestamp time.Time, kind string, data []byte) *Event {
	if timelineId < 1 {
		panic(fmt.Errorf("Invalid timelineId: %d", timelineId))
	}

	return m.addEvent(timelineId, timestamp, kind, data)
}

func (m *nexus) addEvent(timelineId uint64, timestamp time.Time, kind string, data []byte) *Event {
	s := m.getTimeline(timelineId)
	e := NewEvent(m.environment.NextEventId(), timelineId, kind, timestamp, m.environment.Now().UTC(), data)
	n := newNode(e)
	s.add(n)
	return e
}

func (m *nexus) AddTimer(timelineId uint64, timestamp time.Time, kind string, data []byte) *Event {
	if timelineId < 1 {
		panic(fmt.Errorf("Invalid timelineId: %d", timelineId))
	}
	return m.addTimer(timelineId, timestamp, kind, data)
}

func (m *nexus) addTimer(timelineId uint64, timestamp time.Time, kind string, data []byte) *Event {
	e := m.addEvent(timelineId, timestamp, kind, data)
	n := newNode(e)
	m.timers.add(n)
	return e
}

func (m *nexus) Remove(e *Event) bool {
	t := m.getTimeline(e.Timeline)
	return t.removeById(e.Id)
}

func (m *nexus) RemoveRange(r Ranger) int {
	count := 0
	for i := r.Range(); i.Next(); {
		if m.Remove(i.Event()) {
			count = count + 1
		}
	}
	return count
}
