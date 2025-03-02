package mqtt

import mqtt "github.com/eclipse/paho.mqtt.golang"

type MQTTClient interface {
	Connect() error
	Publish(topic string, payload interface{}) error
	Subscribe(topic string, handler mqtt.MessageHandler) error
	Close()
}
