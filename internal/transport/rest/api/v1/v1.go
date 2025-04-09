package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/safe-homie/backend/docs"
	"github.com/safe-homie/backend/infrastructure"
	"github.com/safe-homie/backend/internal/service/demo"
	"github.com/safe-homie/backend/internal/service/device"
	"github.com/safe-homie/backend/internal/service/sensor"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type APIV1 struct {
	context       infrastructure.AppContext
	demoService   *demo.Demo
	sensorService sensor.SensorService
	deviceService device.DeviceService
}

func New(context infrastructure.AppContext, infra infrastructure.AppInfra, srv infrastructure.AppService) *APIV1 {
	return &APIV1{
		context:       context,
		demoService:   new(demo.Demo),
		sensorService: srv.SensorService(),
		deviceService: srv.DeviceService(),
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
	device := v1.Group("/devices")
	{
		device.GET("/:id", Wrap(a.GetDevice))
		device.GET("", Wrap(a.ListDevices))
		device.GET("/:id/status", Wrap(a.DeviceStatus))     // kiểm tra tính tiện ích
		device.GET("/:id/history", Wrap(a.DeviceHistories)) // bổ sung chuẩn hóa thời gian
		device.GET("/:id/schedule", Wrap(a.DeviceSchedules))

		/////////////////////////////////////////////////////////
		// Command Pattern
		device.PUT("/:id", Wrap(a.UpdateDevice))
		device.POST("", Wrap(a.CreateDevice)) // xem lại sửa đổi
		device.POST("/:id/control", Wrap(a.ControlDevice))
		device.POST("/:id/schedule", Wrap(a.ScheduleDevice))
		device.PUT("/:id/schedule/:schedule_id", Wrap(a.UpdateScheduleDevice))
		device.DELETE("/:id/schedule/:schedule_id", Wrap(a.DeleteScheduleDevice))

		// Update later
		device.POST("room/auto-mode", Wrap(a.AutoModeSchedule))
		device.POST("room/apply-preset", Wrap(a.PresetModeSchedule))

	}
}

func (a *APIV1) Demo(ctx echo.Context) *Response {
	data := a.demoService.Demo()
	return &Response{
		Code: http.StatusOK,
		Data: data,
	}
}
