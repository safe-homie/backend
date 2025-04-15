package infrastructure

import (
	"github.com/safe-homie/backend/internal/store"
	"github.com/safe-homie/backend/pkg/event"
	"github.com/safe-homie/backend/pkg/mqtt"
)

type AppInfra interface {
	Store() store.Store
	MQTT() mqtt.MQTTClient
	Event() event.EventManager
}

type appInfra struct {
	store store.Store
	mqtt  mqtt.MQTTClient
	event event.EventManager
}

func NewAppInfra(st store.Store, mc mqtt.MQTTClient, ev event.EventManager) AppInfra {
	return &appInfra{
		store: st,
		mqtt:  mc,
		event: ev,
	}
}

func (ai *appInfra) Store() store.Store {
	return ai.store
}

func (ai *appInfra) MQTT() mqtt.MQTTClient {
	return ai.mqtt
}

func (ai *appInfra) Event() event.EventManager {
	return ai.event
}
