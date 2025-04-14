package app

import (
	"github.com/safe-homie/backend/infrastructure"
	"github.com/safe-homie/backend/internal/config"
	"github.com/safe-homie/backend/internal/domain"
	"github.com/safe-homie/backend/internal/store"
	"github.com/safe-homie/backend/internal/store/db"
	"github.com/safe-homie/backend/pkg/event"
	"github.com/safe-homie/backend/pkg/logger"
	_mqtt "github.com/safe-homie/backend/pkg/mqtt"
	"github.com/safe-homie/backend/pkg/validator"
)

func InitApp() (infrastructure.AppContext, infrastructure.AppInfra, infrastructure.AppService, error) {
	cfg := config.New()
	logger := logger.New()
	driver, err := db.New(cfg)
	if err != nil {
		return nil, nil, nil, err
	}
	store := store.New(driver)
	validator := validator.New()
	mqttClient := _mqtt.New(cfg)
	eventManager := event.New()

	context := infrastructure.NewAppContext(cfg, logger, validator)
	infra := infrastructure.NewAppInfra(store, mqttClient, eventManager)
	srv := infrastructure.NewAppService(store, eventManager)
	if err := store.Migrate(); err != nil {
		// TODO: Needs to refactor this + Add graceful shutdown
		Close(infra)
		return nil, nil, nil, err
	} else {
		context.Logger().Info("migrate completed")
	}
	eventManager.RegisterEvent(domain.SensorThresholdExceed, srv.NotifyService().Notify)
	eventManager.RegisterEvent(domain.SensorThresholdExceed, mqttClient.Publish)
	// TODO: Need method to control device from DeviceService
	// eventManager.RegisterEvent("sensor:thresholdexceeded", ...)
	return context, infra, srv, nil
}
