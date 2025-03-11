package handler

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/safe-homie/backend/internal/domain"
)

func (m *mQTTHandler) HandleSensorMessage(client mqtt.Client, msg mqtt.Message) {
	var sensorData domain.SensorDataMessage
	if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
		log.Printf("error parsing sensor message: %v", err)
		return
	}
	if _, err := m.service.InsertSensorData(&sensorData); err != nil {
		return
	}
}
