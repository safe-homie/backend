package infrastructure

import (
	"github.com/safe-homie/backend/internal/config"
	"github.com/safe-homie/backend/internal/store"
	"github.com/safe-homie/backend/pkg/logger"
	"github.com/safe-homie/backend/pkg/validator"
)

type AppContext interface {
	Config() *config.Config
	Logger() logger.Logger
	Validator() validator.Validator
	Store() store.Store
}

type appContext struct {
	config    *config.Config
	logger    logger.Logger
	validator validator.Validator
	store     store.Store
}

func New(cfg *config.Config, lg logger.Logger, st store.Store, vldt validator.Validator) *appContext {
	return &appContext{
		config:    cfg,
		logger:    lg,
		store:     st,
		validator: vldt,
	}
}

func (ac *appContext) Config() *config.Config {
	return ac.config
}

func (ac *appContext) Logger() logger.Logger {
	return ac.logger
}

func (ac *appContext) Validator() validator.Validator {
	return ac.validator
}

func (ac *appContext) Store() store.Store {
	return ac.store
}
