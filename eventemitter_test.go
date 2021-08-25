package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestEventEmitter(t *testing.T) {
	ee := NewEventEmitter()
	count := 0
	handler := func(msg json.RawMessage) {
		fmt.Println(msg)
		count++
	}

	if count != 0 {
		t.Fatal("Count not 0")
	}
	ee.AddListener("my-event", handler)
	ee.Emit("my-event", []byte{1, 2, 3})
	ee.RemoveListener("my-event", handler)

	if count != 1 {
		t.Fatal("Count not 0")
	}
}
