package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/safe-homie/backend/infrastructure"
	"github.com/safe-homie/backend/internal/service"
)

type APIV1 struct {
	infra       infrastructure.AppContext
	demoService *service.Demo
}

func New(infra infrastructure.AppContext) *APIV1 {
	return &APIV1{infra: infra, demoService: new(service.Demo)}
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
