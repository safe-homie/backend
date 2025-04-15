package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/safe-homie/backend/internal/config"
)

type mqttClient struct {
	client mqtt.Client
	config config.MQTTConfig
}

func New(cfg *config.Config) MQTTClient {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.MQTT.BrokerAddress())
	opts.SetClientID(cfg.MQTT.ClientID)
	// opts.SetUsername(cfg.MQTT.Username)
	// opts.SetPassword(cfg.MQTT.Password)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectTimeout(10 * time.Second)
	client := mqtt.NewClient(opts)
	return &mqttClient{
		client: client,
		config: cfg.MQTT,
	}
}

func (c *mqttClient) Connect() error {
	token := c.client.Connect()
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *mqttClient) Publish(topic string, payload any) error {
	var finalPayload interface{} = payload
	if _, ok := payload.(map[string]interface{}); ok {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal map to json: %w", err)
		}
		finalPayload = string(jsonBytes)
	}
	token := c.client.Publish(topic, 0, false, finalPayload)
	token.Wait()
	return token.Error()
}

func (c *mqttClient) Subscribe(topic string, handler mqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, 0, handler)
	token.Wait()
	if token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *mqttClient) Disconnect() {
	c.client.Disconnect(250)
}
