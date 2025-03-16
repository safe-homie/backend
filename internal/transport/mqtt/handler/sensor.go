package handler

import (
	"encoding/json"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/safe-homie/backend/internal/domain"
	"github.com/safe-homie/backend/internal/util"
)

func (m *mQTTHandler) HandleSensorMessage(client mqtt.Client, msg mqtt.Message) {
	var sensorData domain.SensorDataMessage
	if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
		log.Printf("error parsing sensor message: %v", err)
		return
	}
	topic := msg.Topic()
	sensorType := extractSensorTypeFromTopic(topic)
	sensorData.Type = sensorType
	if _, err := m.service.InsertSensorData(&sensorData); err != nil {
		log.Printf("error inserting sensor message: %v", err)
		return
	}
}

func extractSensorTypeFromTopic(topic string) string {
	parts := strings.Split(topic, "/")
	return util.ExtractPart(parts, 1, "unknown")
}
