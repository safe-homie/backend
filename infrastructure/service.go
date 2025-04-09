package infrastructure

import (
	"github.com/safe-homie/backend/internal/service/device"
	"github.com/safe-homie/backend/internal/service/sensor"
	"github.com/safe-homie/backend/internal/store"
)

type AppService interface {
	SensorService() sensor.SensorService
	DeviceService() device.DeviceService
}

type appService struct {
	sensorService sensor.SensorService
	deviceService device.DeviceService
}

func NewAppService(store store.Store) AppService {
	return &appService{
		sensorService: sensor.NewService(store),
		deviceService: device.NewService(store),
	}
}

func (ap *appService) SensorService() sensor.SensorService {
	return ap.sensorService
}

func (ap *appService) DeviceService() device.DeviceService {
	return ap.deviceService
}
