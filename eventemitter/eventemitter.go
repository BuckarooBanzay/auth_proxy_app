package eventemitter

import (
	"sync"
)

type EventEmitter struct {
	handlers []EventHandler
	mutex    *sync.RWMutex
}

type EventHandler interface {
	OnEvent(name string, payload interface{})
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		handlers: []EventHandler{},
		mutex:    &sync.RWMutex{},
	}
}

func (ee *EventEmitter) AddListener(name string, handler EventHandler) {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()

	ee.handlers = append(ee.handlers, handler)
}

func (ee *EventEmitter) Emit(name string, payload interface{}) {
	ee.mutex.RLock()
	defer ee.mutex.RUnlock()

	for _, handler := range ee.handlers {
		handler.OnEvent(name, payload)
	}
}

func (ee *EventEmitter) RemoveListener(name string, remove_handler EventHandler) {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()

	newlist := []EventHandler{}

	for _, handler := range ee.handlers {
		if remove_handler != handler {
			newlist = append(newlist, handler)
		}
	}

	ee.handlers = newlist
}
