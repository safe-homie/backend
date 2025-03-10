package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/safe-homie/backend/infrastructure"
	v1 "github.com/safe-homie/backend/internal/transport/rest/api/v1"
)

type Server struct {
	api     *v1.APIV1
	router  *echo.Echo
	context infrastructure.AppContext
}

func NewServer(context infrastructure.AppContext, infra infrastructure.AppInfra, srv infrastructure.AppService) *Server {
	server := &Server{
		router:  echo.New(),
		api:     v1.New(context, infra, srv),
		context: context,
	}
	server.setupMiddlewares()
	server.registerRoutes()
	return server
}

func (s *Server) Start() error {
	s.context.Logger().Info("api server started")
	return s.router.Start(s.context.Config().AppAddress())
}
