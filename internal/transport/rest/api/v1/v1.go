package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/safe-homie/backend/docs"
	"github.com/safe-homie/backend/infrastructure"
	"github.com/safe-homie/backend/internal/service/demo"
	"github.com/safe-homie/backend/internal/service/sensor"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type APIV1 struct {
	context       infrastructure.AppContext
	demoService   *demo.Demo
	sensorService sensor.SensorService
}

func New(context infrastructure.AppContext, infra infrastructure.AppInfra, srv infrastructure.AppService) *APIV1 {
	return &APIV1{
		context:       context,
		demoService:   new(demo.Demo),
		sensorService: srv.SensorService(),
	}
}

func (a *APIV1) RegisterHandlers(e *echo.Echo) {
	e.GET("/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/docs/index.html")
	})
	e.GET("/docs/*", echoSwagger.WrapHandler)

	v1 := e.Group("/api/v1")
	v1.GET("/demo", Wrap(a.Demo))

	sensor := v1.Group("/sensors")
	{
		sensor.GET("/:id", Wrap(a.GetSensor))
		sensor.GET("", Wrap(a.ListSensors))
		sensor.POST("", Wrap(a.CreateSensor))
		sensor.GET("/data/latest", Wrap(a.ListLatestSensorDataByLocation))
	}
}

func (a *APIV1) Demo(ctx echo.Context) *Response {
	data := a.demoService.Demo()
	return &Response{
		Code: http.StatusOK,
		Data: data,
	}
}
