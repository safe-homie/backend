package handler

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (m *mQTTHandler) HandleSensorMessage(client mqtt.Client, msg mqtt.Message) {
	// TODO: Need to parse payload, call service to insert data
	fmt.Printf("received temperature data: %s", msg.Payload())
}
