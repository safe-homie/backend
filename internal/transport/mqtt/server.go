package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
		handler: handler.NewMQTTHandler(srv.SensorService(), srv.DeviceService()),
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
	// s.context.Logger().Info("subscribe sensors topic")
	// err = s.client.Subscribe("sensors/+", s.handler.HandleSensorMessage)
	// if err != nil {
	// 	return
	// }
	// s.context.Logger().Info("subscribe devices topic")
	// if err := s.client.Subscribe("devices/+", s.handler.HandleDeviceMessage); err != nil {
	// 	return err
	// }
	// return nil
	topics := map[string]mqtt.MessageHandler{
		"sensors/+":         s.handler.HandleSensorMessage,
		"devices/+":         s.handler.HandleDeviceMessage,
		"devices/+/status":  s.handler.HandleDeviceMessage,
		"devices/+/control": s.handler.HandleDeviceMessage,
	}

	for topic, handler := range topics {
		s.context.Logger().Info(fmt.Sprintf("Subscribing to topic: %s", topic))
		if err := s.client.Subscribe(topic, handler); err != nil {
			return err
		}
	}
	return nil
}
