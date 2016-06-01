package septum

import (
	"fmt"
)

// ==== timeline ====

type (
	timeline struct {
		id     uint64
		length int
		head   *node
		tail   *node
	}
)

func newTimeline(id uint64) *timeline {
	return &timeline{id: id}
}

func (s *timeline) Id() uint64 {
	return s.id
}

func (s *timeline) Len() int {
	return s.length
}

func (s *timeline) Range() Range {
	return newRanger(s.head)
}

func (s *timeline) add(n *node) {
	if n == nil {
		panic(fmt.Errorf("Can't add nil node to timeline %d", s.id))
	}

	s.length = s.length + 1
	if s.head == nil {
		s.head = n
		s.tail = n
	} else {
		for h := s.head; h != nil; h = h.next {
			if h.e.Timestamp.After(n.e.Timestamp) {
				if h == s.head {
					s.head = n.insertBefore(n)
				} else {
					n.insertBefore(n)
				}
				return
			}
		}
		s.tail = s.tail.insertAfter(n)
	}
}

func (t *timeline) removeById(id uint64) bool {
	for n := t.head; n != nil; n = n.next {
		if n.e.Id == id {
			t.del(n)
			return true
		}
	}
	return false
}

func (s *timeline) del(n *node) {
	if n == nil {
		panic(fmt.Errorf("Can't delete nil node from timeline %d", s.id))
	}

	if s.head == n {
		s.head = n.next
	}
	if s.tail == n {
		s.tail = n.prev
	}
	n.del()
}

// ==== node ====

type (
	node struct {
		prev *node
		next *node
		e    *Event
	}
)

func newNode(e *Event) *node {
	if e == nil {
		panic("e is nil")
	}

	return &node{
		e: e,
	}
}

func (n *node) Range() Range {
	return newRanger(n)
}

func (n *node) insertBefore(curr *node) *node {
	curr.prev = n.prev
	curr.next = n

	if curr.prev != nil {
		curr.prev.next = curr
	}

	n.prev = curr
	return curr
}

func (n *node) insertAfter(curr *node) *node {
	curr.prev = n
	curr.next = n.next

	if curr.next != nil {
		curr.next.prev = curr
	}

	n.next = curr
	return curr
}

func (n *node) del() {
	prev := n.prev
	next := n.next

	if prev != nil {
		prev.next = next
		n.prev = nil
	}

	if next != nil {
		next.prev = prev
		n.next = nil
	}

	n.e = nil
}
