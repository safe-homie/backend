package rest

import (
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) setupMiddlewares() {
	s.router.Use(middleware.CORS())

	logCfg := s.context.Config().GetEchoLogConfig()
	s.router.Use(middleware.LoggerWithConfig(logCfg))
}

func (s *Server) registerRoutes() {
	s.api.RegisterHandlers(s.router)
}
