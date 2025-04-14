package event

import (
	"github.com/asaskevich/EventBus"
)

type eventManager struct {
	bus EventBus.Bus
}

func New() EventManager {
	return &eventManager{bus: EventBus.New()}
}

func (e *eventManager) RegisterEvent(topic string, callback any) error {
	return e.bus.SubscribeAsync(topic, callback, false)
}

func (e *eventManager) EmitEvent(topic string, args ...any) {
	e.bus.Publish(topic, args)
}
