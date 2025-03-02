package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/safe-homie/backend/infrastructure"
	_mqtt "github.com/safe-homie/backend/pkg/mqtt"
)

type Server struct {
	context infrastructure.AppContext
	client  _mqtt.MQTTClient
}

func NewServer(context infrastructure.AppContext, infra infrastructure.AppInfra) *Server {
	// TODO: Server would contains service and passing infra as dependencies
	return &Server{
		context: context,
		client:  infra.MQTT(),
	}
}

func (s *Server) Start() error {
	s.context.Logger().Info("mqtt server started")
	if err := s.client.Connect(); err != nil {
		s.context.Logger().Error(fmt.Sprintf("failed to connect to mqtt broker: %s", err.Error()))
		return err
	}
	// TODO: Refactor to handle multple topics
	s.context.Logger().Info("subscribe to sensors/temp")
	if err := s.client.Subscribe("sensors/temp", s.handleTemperature); err != nil {
		return err
	}
	return nil
}

func (s *Server) handleTemperature(client mqtt.Client, msg mqtt.Message) {
	s.context.Logger().Info(fmt.Sprintf("received temperature data: %s", msg.Payload()))
}
