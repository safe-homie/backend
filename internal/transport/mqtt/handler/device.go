package handler

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/safe-homie/backend/internal/domain"
)

func (h *mQTTHandler) HandleDeviceMessage(client mqtt.Client, msg mqtt.Message) {
	var deviceMsg domain.DeviceMessage
	topic := msg.Topic()

	if err := json.Unmarshal(msg.Payload(), &deviceMsg); err != nil {
		log.Printf("error parsing device message: %v", err)
		return
	}

	parts := strings.Split(topic, "/")
	if len(parts) < 2 {
		log.Printf("Invalid device topic format")
		return
	}
	idStr := parts[1]
	idInt, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		log.Printf("invalid deviceID: %v", err)
		return
	}
	deviceMsg.DeviceID = int32(idInt)

	switch {
	case strings.HasSuffix(topic, "/status"):
		if err := h.Device.InsertStatusData(&deviceMsg); err != nil {
			log.Printf("Failed to update device status: %v", err)
		}
	case strings.HasSuffix(topic, "/control"):
		if err := h.Device.HandleControlMessage(&deviceMsg); err != nil {
			log.Printf("Failed to handle control message: %v", err)
		}
	default:
		log.Printf("Unknown device message type: %s", topic)
	}
}
