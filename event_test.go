package septum

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"testing"
	"time"
)

func TestJsonEncodingOfEventData(t *testing.T) {

	buf := new(bytes.Buffer)
	encEvent := &Event{
		Id:          1,
		Timeline:    2,
		Timestamp:   time.Now(),
		ReceiveTime: time.Now(),
		Kind:        "TestEvent",
		Data:        []byte("Some event data"),
	}

	if err := json.NewEncoder(buf).Encode(encEvent); err != nil {
		t.Error("Encoding error:", err)
	}

	decEvent := &Event{}

	if err := json.NewDecoder(buf).Decode(decEvent); err != nil {
		t.Error("Decoding error:", err)
	}

	eventEquals(t, encEvent, decEvent)
}

func TestGobEncodingOfEventData(t *testing.T) {
	buf := new(bytes.Buffer)
	encEvent := newTestEvent()

	if err := gob.NewEncoder(buf).Encode(encEvent); err != nil {
		t.Error("Encoding error:", err)
	}

	decEvent := &Event{}
	if err := gob.NewDecoder(buf).Decode(decEvent); err != nil {
		t.Error("Decoding error:", err)
	}

	eventEquals(t, encEvent, decEvent)
}

func BenchmarkEventMarshalBinary(b *testing.B) {
	encEvent := newTestEvent()
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := encEvent.MarshalBinary(); err != nil {
			b.Error("Error during marshalling:", err)
		}
	}
}

func BenchmarkEventUnmarshalBinary(b *testing.B) {
	encEvent := newTestEvent()
	decEvent := &Event{}
	data, err := encEvent.MarshalBinary()
	b.Logf("Using data:%d", len(data))
	if err != nil {
		b.Error("Error during marshalling:", err)
	}
	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if err := decEvent.UnmarshalBinary(data); err != nil {
			b.Error("Error during unmarshalling:", err)
		}
	}
}

func newTestEvent() *Event {
	return &Event{
		Id:          1,
		Timeline:    2,
		Timestamp:   time.Now().UTC(),
		ReceiveTime: time.Now().UTC(),
		Kind:        "TestEvent",
		Data:        []byte("Some event data"),
	}
}

func eventEquals(t *testing.T, encEvent, decEvent *Event) {
	if encEvent.Id != decEvent.Id {
		t.Errorf("Error roundtripping Event.Id: %v != %v", encEvent.Id, decEvent.Id)
	}

	if encEvent.Timeline != decEvent.Timeline {
		t.Errorf("Error roundtripping Event.Timeline: %v != %v", encEvent.Timeline, decEvent.Timeline)
	}

	if encEvent.Timestamp != decEvent.Timestamp {
		t.Errorf("Error roundtripping Event.Timestamp: %v != %v", encEvent.ReceiveTime, decEvent.ReceiveTime)
	}

	if encEvent.ReceiveTime != decEvent.ReceiveTime {
		t.Errorf("Error roundtripping Event.ReceiveTime: %v != %v", encEvent.ReceiveTime, decEvent.ReceiveTime)
	}

	if encEvent.Kind != decEvent.Kind {
		t.Errorf("Error roundtripping Event.Kind: %v != %v", encEvent.Kind, decEvent.Kind)
	}

	if !sliceEquals(encEvent.Data, decEvent.Data) {
		t.Errorf("Error roundtripping Event.Data: %v != %v", encEvent.Data, decEvent.Data)
	}
}

func sliceEquals(x, y []byte) bool {
	if len(x) != len(y) {
		return false
	}
	for i := 0; i < len(x); i++ {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
