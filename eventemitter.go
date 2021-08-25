package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

type EventEmitter struct {
	emitters map[string][]EventHandler
	mutex    *sync.RWMutex
}

type EventHandler func(json.RawMessage)

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		emitters: make(map[string][]EventHandler),
		mutex:    &sync.RWMutex{},
	}
}

func (ee EventEmitter) getList(name string) []EventHandler {
	list := ee.emitters[name]
	if list == nil {
		list = make([]EventHandler, 0)
	}
	return list
}

func (ee EventEmitter) setList(name string, list []EventHandler) {
	ee.emitters[name] = list
}

func (ee EventEmitter) AddListener(name string, handler EventHandler) {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()

	list := ee.getList(name)
	list = append(list, handler)
	ee.setList(name, list)
}

func (ee EventEmitter) Emit(name string, msg json.RawMessage) {
	ee.mutex.RLock()
	defer ee.mutex.RUnlock()

	list := ee.getList(name)
	for _, handler := range list {
		handler(msg)
	}
}

func (ee EventEmitter) RemoveListener(name string, remove_handler EventHandler) {
	ee.mutex.Lock()
	defer ee.mutex.Unlock()

	list := ee.getList(name)
	newlist := make([]EventHandler, 0)

	remove_handler_p := fmt.Sprintf("%v", remove_handler)

	for _, handler := range list {
		handler_p := fmt.Sprintf("%v", handler)
		if remove_handler_p != handler_p {
			newlist = append(newlist, handler)
		}
	}

	ee.setList(name, newlist)
}
