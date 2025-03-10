package handler

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (m *mQTTHandler) HandleSensorMessage(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("received temperature data: %s", msg.Payload())
	fmt.Printf("received temperature data: %s", msg.Topic())
}
