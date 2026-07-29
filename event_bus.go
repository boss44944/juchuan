package main

import (
	"encoding/json"
	"sync"
)

// EventBus is a lightweight application event broadcaster.
// WebSocket layer can subscribe to these events.
type EventBus struct {
	mu sync.RWMutex
	subscribers []func([]byte)
}

func NewEventBus() *EventBus {
	return &EventBus{}
}

func (b *EventBus) Subscribe(fn func([]byte)) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers = append(b.subscribers, fn)
}

func (b *EventBus) Publish(event interface{}) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, fn := range b.subscribers {
		fn(data)
	}
}
