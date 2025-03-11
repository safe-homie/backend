package handler

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/safe-homie/backend/internal/service/sensor"
)

type MQTTHandler interface {
	HandleSensorMessage(client mqtt.Client, msg mqtt.Message)
}

type mQTTHandler struct {
	service sensor.SensorService
}

func NewMQTTHandler(service sensor.SensorService) MQTTHandler {
	return &mQTTHandler{service: service}
}
