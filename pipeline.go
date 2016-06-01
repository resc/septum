package septum

import (
	"errors"
)

type (
	// Filter determines if an event is selected for processing
	Filter func(e *Event) bool

	// Process processes an event
	Process func(e *Event) error

	//Pipeline combines event processing steps
	Pipeline struct {
		filter  Filter
		process Process
		prev    *Pipeline
		next    *Pipeline
	}
)

// NewPipeline creates a new Pipeline
func NewPipeline(process Process) *Pipeline {
	head := Pipeline{}
	head.Append(process)
	h, _ := head.Next()
	head.Delete()
	return h
}

// NewPipeline creates a new Pipeline with a filtered first step.
func NewPipelineFiltered(filter Filter, process Process) *Pipeline {
	head := Pipeline{}
	head.AppendFiltered(filter, process)
	h, _ := head.Next()
	head.Delete()
	return h
}

// Next returns the next segment in the pipeline and true, or nil and false.
func (p *Pipeline) Next() (*Pipeline, bool) {
	return p.next, p.next != nil
}

// Delete deletes the segment from the pipeline.
func (p *Pipeline) Delete() {
	prev := p.prev
	next := p.next

	if prev != nil {
		prev.next = next
		p.prev = nil
	}

	if next != nil {
		next.prev = prev
		p.next = nil
	}
}

// Append appends a  processing step at the end of the pipeline.
func (p *Pipeline) Append(process Process) {
	p.AppendFiltered(acceptAllEvents, process)
}

// AppendFilter appends a filtered processing step at the end of the pipeline.
func (p *Pipeline) AppendFiltered(filter Filter, process Process) {
	n := p
	for n.next != nil {
		n = n.next
	}
	n.InsertAfterFiltered(filter, process)
}

// InsertAfter inserts a processing step after the current step of the pipeline.
func (p *Pipeline) InsertAfter(process Process) {
	p.InsertAfterFiltered(acceptAllEvents, process)
}

// InsertAfter inserts a filtered processing step after the current step of the pipeline.
func (p *Pipeline) InsertAfterFiltered(filter Filter, process Process) {
	if filter == nil {
		panic(errors.New("filter parameter is nil"))
	}

	if process == nil {
		panic(errors.New("process parameter is nil"))
	}

	n := &Pipeline{
		filter:  filter,
		process: process,
		prev:    p,
		next:    p.next,
	}

	p.next = n
}

// Process processes the event using the pipeline.
// If any processing step fails the processing is aborted and the error is returned.
func (p *Pipeline) Process(e *Event) error {
	if e == nil {
		return errors.New("e is nil")
	}

	for n := p; n != nil; n = n.next {
		if n.filter(e) {
			err := n.process(e)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func acceptAllEvents(e *Event) bool {
	return true
}
