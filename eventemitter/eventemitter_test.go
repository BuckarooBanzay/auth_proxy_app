package eventemitter

import (
	"testing"
)

type MyHandler struct {
	count int
}

func (mh *MyHandler) OnEvent(name string, payload interface{}) {
	str, ok := payload.(string)
	if !ok {
		panic("not a string")
	}
	if str != "test" {
		panic("string does not match")
	}
	mh.count++
}

func TestEventEmitter(t *testing.T) {
	ee := NewEventEmitter()
	handler := &MyHandler{}

	if handler.count != 0 {
		t.Fatal("Count not 0")
	}
	ee.AddListener("my-event", handler)
	ee.Emit("my-event", "test")
	ee.RemoveListener("my-event", handler)

	if handler.count != 1 {
		t.Fatal("Count not 0")
	}
}
