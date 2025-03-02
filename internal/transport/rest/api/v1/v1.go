package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/safe-homie/backend/infrastructure"
	"github.com/safe-homie/backend/internal/service"
)

type APIV1 struct {
	context     infrastructure.AppContext
	demoService *service.Demo
}

func New(context infrastructure.AppContext, infra infrastructure.AppInfra) *APIV1 {
	// TODO: APIV1 would contains services and passing infra as dependency
	return &APIV1{context: context, demoService: new(service.Demo)}
}

func (a *APIV1) RegisterHandlers(e *echo.Echo) {
	v1 := e.Group("/api/v1")
	v1.GET("/demo", Wrap(a.Demo))
}

func (a *APIV1) Demo(ctx echo.Context) *Response {
	data := a.demoService.Demo()
	return &Response{
		Code: http.StatusOK,
		Data: data,
	}
}
