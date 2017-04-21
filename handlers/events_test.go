package handlers

import (
	"reflect"
	"testing"
)

type TestHandler struct{}

func TestHandle(t *testing.T) {
	mux := NewEventsMux()
	handler := &TestHandler{}
	mux.Handle("message", handler)

	actual := mux.m["message"]
	expected := muxEntry{explicit: true, h: handler, event: "message"}

	if actual != expected {
		t.Errorf("Failed")
	}
}

func TestEventMatch(t *testing.T) {
	actual := eventMatch("message", "message")
	if actual != true {
		t.Errorf("Failed")
	}
}

func TestMatch(t *testing.T) {
	mux := NewEventsMux()
	handler := &TestHandler{}
	mux.m = make(map[string]muxEntry)
	mux.m["message"] = muxEntry{explicit: true, h: handler, event: "message"}
	actual, _ := mux.match("message")
	match := reflect.DeepEqual(actual, handler)
	if !match {
		t.Errorf("Failed")
	}
}
