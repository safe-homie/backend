package main

import (
	"sync"

	"github.com/safe-homie/backend/infrastructure"
	"github.com/safe-homie/backend/internal/config"
	"github.com/safe-homie/backend/internal/store"
	"github.com/safe-homie/backend/internal/store/db"
	"github.com/safe-homie/backend/internal/transport/mqtt"
	"github.com/safe-homie/backend/internal/transport/rest"
	"github.com/safe-homie/backend/pkg/logger"
	mqtt_cliet "github.com/safe-homie/backend/pkg/mqtt"
	"github.com/safe-homie/backend/pkg/validator"
)

func main() {
	cfg := config.New()
	logger := logger.New()
	driver, err := db.New(cfg)
	if err != nil {
		logger.Error(err.Error())
	}
	store := store.New(driver)
	if err := store.Migrate(); err != nil {
		logger.Error(err.Error())
	} else {
		logger.Info("migrate complete")
	}
	validator := validator.New()
	infra := infrastructure.New(cfg, logger, store, validator)
	client := mqtt_cliet.New(cfg)

	restServer := rest.NewServer(infra)
	mqttServer := mqtt.NewServer(infra, client)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := restServer.Start(); err != nil {
			logger.Error(err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		if err := mqttServer.Start(); err != nil {
			logger.Error(err.Error())
		}
	}()
	wg.Wait()
}
