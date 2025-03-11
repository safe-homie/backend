package infrastructure

import (
	"github.com/safe-homie/backend/internal/service/sensor"
	"github.com/safe-homie/backend/internal/store"
)

type AppService interface {
	SensorService() sensor.SensorService
}

type appService struct {
	sensorService sensor.SensorService
}

func NewAppService(store store.Store) AppService {
	return &appService{
		sensorService: sensor.NewService(store),
	}
}

func (ap *appService) SensorService() sensor.SensorService {
	return ap.sensorService
}
