package main

import (
	"github.com/safe-homie/backend/infrastructure"
	"github.com/safe-homie/backend/internal/config"
	"github.com/safe-homie/backend/internal/store"
	"github.com/safe-homie/backend/internal/store/db"
	"github.com/safe-homie/backend/internal/transport/rest"
	"github.com/safe-homie/backend/pkg/logger"
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
	server := rest.NewServer(infra)
	if err := server.Start(); err != nil {
		logger.Error(err.Error())
	}
}
