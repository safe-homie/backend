package infrastructure

import (
	"github.com/safe-homie/backend/internal/store"
	"github.com/safe-homie/backend/pkg/mqtt"
)

type AppInfra interface {
	Store() store.Store
	MQTT() mqtt.MQTTClient
}

type appInfra struct {
	store store.Store
	mqtt  mqtt.MQTTClient
}

func NewAppInfra(st store.Store, mc mqtt.MQTTClient) AppInfra {
	return &appInfra{
		store: st,
		mqtt:  mc,
	}
}

func (ai *appInfra) Store() store.Store {
	return ai.store
}

func (ai *appInfra) MQTT() mqtt.MQTTClient {
	return ai.mqtt
}
