package infrastructure

import (
	"github.com/safe-homie/backend/internal/service/notify"
	"github.com/safe-homie/backend/internal/service/sensor"
	"github.com/safe-homie/backend/internal/store"
	"github.com/safe-homie/backend/pkg/event"
)

type AppService interface {
	SensorService() sensor.SensorService
	NotifyService() notify.NotifyService
}

type appService struct {
	sensorService sensor.SensorService
	notifyService notify.NotifyService
}

func NewAppService(store store.Store, evm event.EventManager) AppService {
	return &appService{
		sensorService: sensor.NewService(store, evm),
		notifyService: notify.NewService(store),
	}
}

func (ap *appService) SensorService() sensor.SensorService {
	return ap.sensorService
}

func (ap *appService) NotifyService() notify.NotifyService {
	return ap.notifyService
}
