package config

import "fmt"

type MQTTConfig struct {
	Host     string
	Port     string
	ClientID string
	Username string
	Password string
	Topics   []string
}

func (m *MQTTConfig) BrokerAddress() string {
	return fmt.Sprintf("tcp://%s:%s", m.Host, m.Port)
}
