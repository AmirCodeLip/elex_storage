package message_broker

import (
	"encoding/json"
	"net/http"
	"sync"
)

type EventBus struct {
	subscribers map[string][]chan string
	mutex       sync.RWMutex
}

type EventRequest struct {
	Topic string          `json:"topic"`
	Data  json.RawMessage `json:"data"`
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan string),
	}
}

func (bus *EventBus) Subscribe(topic string) <-chan string {
	ch := make(chan string, 1)

	bus.mutex.Lock()
	bus.subscribers[topic] = append(bus.subscribers[topic], ch)
	bus.mutex.Unlock()

	return ch
}

func (bus *EventBus) Publish(topic string, msg interface{}) {
	// bus.mutex.RLock()
	// subs := bus.subscribers[topic]
	// bus.mutex.RUnlock()

	// for _, ch := range subs {
	// 	go func(c chan interface{}) {
	// 		c <- msg
	// 	}(ch)
	// }
}
func (bus *EventBus) Listen(w http.ResponseWriter, r *http.Request) {
	var eventRequest EventRequest

	if err := json.NewDecoder(r.Body).Decode(&eventRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bus.mutex.RLock()
	subs := bus.subscribers[eventRequest.Topic]
	bus.mutex.RUnlock()

	for _, ch := range subs {
		go func(c chan string, data string) {
			c <- data
		}(ch, string(eventRequest.Data))
	}

	w.WriteHeader(http.StatusOK)
}
