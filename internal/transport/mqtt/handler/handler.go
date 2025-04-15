package handler

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/safe-homie/backend/internal/service/device"
	"github.com/safe-homie/backend/internal/service/sensor"
)

type MQTTHandler interface {
	HandleSensorMessage(client mqtt.Client, msg mqtt.Message)
	HandleDeviceMessage(client mqtt.Client, msg mqtt.Message)
}

type mQTTHandler struct {
	Sensor sensor.SensorService
	Device device.DeviceService
}

func NewMQTTHandler(sensorSvc sensor.SensorService, deviceSvc device.DeviceService) MQTTHandler {
	return &mQTTHandler{
		Sensor: sensorSvc,
		Device: deviceSvc,
	}
}
