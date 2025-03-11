package mqtt

import (
	"fmt"

	"github.com/safe-homie/backend/infrastructure"
	"github.com/safe-homie/backend/internal/transport/mqtt/handler"
	_mqtt "github.com/safe-homie/backend/pkg/mqtt"
)

type Server struct {
	context infrastructure.AppContext
	client  _mqtt.MQTTClient
	handler handler.MQTTHandler
}

func NewServer(context infrastructure.AppContext, infra infrastructure.AppInfra, srv infrastructure.AppService) *Server {
	return &Server{
		context: context,
		client:  infra.MQTT(),
		handler: handler.NewMQTTHandler(srv.SensorService()),
	}
}

func (s *Server) Start() error {
	s.context.Logger().Info("mqtt server started")
	if err := s.client.Connect(); err != nil {
		s.context.Logger().Error(fmt.Sprintf("failed to connect to mqtt broker: %s", err.Error()))
		return err
	}
	if err := s.registerHandlers(); err != nil {
		return err
	}
	return nil
}

func (s *Server) registerHandlers() (err error) {
	s.context.Logger().Info("subscribe sensors topic")
	err = s.client.Subscribe("sensors/+", s.handler.HandleSensorMessage)
	if err != nil {
		return
	}
	return nil
}
