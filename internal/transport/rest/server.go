package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/safe-homie/backend/infrastructure"
	v1 "github.com/safe-homie/backend/internal/transport/rest/api/v1"
)

type Server struct {
	api    *v1.APIV1
	router *echo.Echo
	infra  infrastructure.AppContext
}

func NewServer(infra infrastructure.AppContext) *Server {
	server := &Server{
		router: echo.New(),
		api:    v1.New(infra),
		infra:  infra,
	}
	server.setupMiddlewares()
	server.registerRoutes()
	return server
}

func (s *Server) Start() error {
	s.infra.Logger().Info("server started")
	return s.router.Start(s.infra.Config().AppAddress())
}
