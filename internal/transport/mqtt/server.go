package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/safe-homie/backend/infrastructure"
	_mqtt "github.com/safe-homie/backend/pkg/mqtt"
)

type Server struct {
	infra  infrastructure.AppContext
	client _mqtt.MQTTClient
}

func NewServer(infra infrastructure.AppContext, client _mqtt.MQTTClient) *Server {
	return &Server{
		infra:  infra,
		client: client,
	}
}

func (s *Server) Start() error {
	s.infra.Logger().Info("mqtt server started")
	if err := s.client.Connect(); err != nil {
		s.infra.Logger().Error(fmt.Sprintf("failed to connect to mqtt broker: %s", err.Error()))
		return err
	}
	// TODO: Refactor to handle multple topics
	s.infra.Logger().Info("subscribe to sensors/temp")
	if err := s.client.Subscribe("sensors/temp", s.handleTemperature); err != nil {
		return err
	}
	return nil
}

func (s *Server) handleTemperature(client mqtt.Client, msg mqtt.Message) {
	s.infra.Logger().Info(fmt.Sprintf("received temperature data: %s\n", msg.Payload()))
}
