package septum

import (
	"encoding/binary"
	"errors"
	"time"
)

// ==== event ====

type (
	Event struct {
		Id          uint64    `json:"id"`
		Timeline    uint64    `json:"timeline"`
		Timestamp   time.Time `json:"timestamp"`
		ReceiveTime time.Time `json:"receive_time"`
		Kind        string    `json:"kind"`
		Data        []byte    `json:"data,omitempty"`
	}

	EventData struct {
		Timeline  uint64    `json:"timeline"`
		Timestamp time.Time `json:"timestamp"`
		Kind      string    `json:"kind"`
		Data      []byte    `json:"data,omitempty"`
	}
)

func NewEvent(eventId uint64, timeline uint64, kind string, timestamp time.Time, receiveTime time.Time, data []byte) *Event {
	return &Event{
		Id:          eventId,
		Timeline:    timeline,
		Kind:        kind,
		Timestamp:   timestamp,
		ReceiveTime: receiveTime,
		Data:        data,
	}
}

func (e *Event) MarshalBinary() ([]byte, error) {

	if !e.Timestamp.IsZero() && e.Timestamp.Location() != time.UTC {
		return nil, errors.New("event timestamp is not in UTC time")
	}
	if !e.ReceiveTime.IsZero() && e.ReceiveTime.Location() != time.UTC {
		return nil, errors.New("event receive time is not in UTC time")
	}

	l := 8 + 8 + (8 + 4 + 2) + (8 + 4 + 2) + 8 + len(e.Kind) + 8 + len(e.Data)
	b := make([]byte, l)
	o := binary.PutUvarint(b, e.Id)
	o += binary.PutUvarint(b[o:], e.Timeline)
	o += binary.PutVarint(b[o:], e.Timestamp.Unix())
	o += binary.PutVarint(b[o:], int64(e.Timestamp.Nanosecond()))
	o += binary.PutVarint(b[o:], e.ReceiveTime.Unix())
	o += binary.PutVarint(b[o:], int64(e.ReceiveTime.Nanosecond()))

	o += binary.PutVarint(b[o:], int64(len(e.Kind)))
	if len(e.Kind) > 0 {
		o += copy(b[o:], e.Kind)
	}

	o += binary.PutVarint(b[o:], int64(len(e.Data)))
	if len(e.Data) > 0 {
		o += copy(b[o:], e.Data)
	}
	return b[:o], nil

}

func (v *Event) UnmarshalBinary(data []byte) error {
	_, err := v.unmarshalBinary(data)
	return err
}

func (v *Event) unmarshalBinary(data []byte) (int, error) {
	o := 0
	id, size := binary.Uvarint(data[o:])
	if size <= 0 {
		return o, errors.New("Error unmarshalling Event.Id")
	}

	o += size
	timeline, size := binary.Uvarint(data[o:])
	if size <= 0 {
		return o, errors.New("Error unmarshalling Event.Timeline")
	}

	o += size
	tUnix, size := binary.Varint(data[o:])
	if size <= 0 {
		return o, errors.New("Error unmarshalling Event.Timestamp")
	}

	o += size
	tNano, size := binary.Varint(data[o:])
	if size <= 0 {
		return o, errors.New("Error unmarshalling Event.Timestamp")
	}

	o += size
	rUnix, size := binary.Varint(data[o:])
	if size <= 0 {
		return o, errors.New("Error unmarshalling Event.ReceiveTime")
	}

	o += size
	rNano, size := binary.Varint(data[o:])
	if size <= 0 {
		return o, errors.New("Error unmarshalling Event.ReceiveTime")
	}

	o += size
	lKind, size := binary.Varint(data[o:])
	if size <= 0 || lKind < 0 || len(data[o:])-size < int(lKind) {
		return o, errors.New("Error unmarshalling Event.Kind")
	}

	o += size
	kind := string(data[o : o+int(lKind)])

	o += int(lKind)
	lData, size := binary.Varint(data[o:])
	if size <= 0 || lData < 0 || len(data[o:])-size < int(lData) {
		return o, errors.New("Error unmarshalling Event.Data")
	}

	o += size
	data = data[o : o+int(lData)]
	o += int(lData)

	v.Id = id
	v.Timeline = timeline
	v.Timestamp = time.Unix(tUnix, tNano).UTC()
	v.ReceiveTime = time.Unix(rUnix, rNano).UTC()
	v.Kind = kind
	v.Data = data
	return o, nil
}
