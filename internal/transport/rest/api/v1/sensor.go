package v1

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/safe-homie/backend/internal/domain"
)

func (a *APIV1) GetSensor(ctx echo.Context) *Response {
	idStr := ctx.Param("id")
	if idStr == "" {
		return ErrorInvalidPathParam
	}
	idInt, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return ErrorInvalidPathParam
	}
	sensor, err := a.sensorService.GetSensor(int32(idInt))
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get sensor: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToSensorDTO(sensor))
}

func (a *APIV1) ListSensors(ctx echo.Context) *Response {
	location := ctx.QueryParam("location")
	sensors, err := a.sensorService.ListSensors(location)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get sensor list: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToSensorsDTO(sensors))
}

func (a *APIV1) CreateSensor(ctx echo.Context) *Response {
	var req domain.CreateSensorRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorInvalidRequestBody
	}
	if err := a.context.Validator().Validate(req); err != nil {
		return ErrorInvalidRequestBody
	}
	sensor, err := a.sensorService.CreateSensor(&req)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to create sensor: "+err.Error())
	}
	return NewDataResponse(http.StatusCreated, domain.ToSensorDTO(sensor))
}

func (a *APIV1) GetLatestSensorData(ctx echo.Context) *Response {
	idStr := ctx.Param("id")
	if idStr == "" {
		return ErrorInvalidPathParam
	}
	idInt, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return ErrorInvalidPathParam
	}
	data, err := a.sensorService.GetLatestSensorData(int32(idInt))
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get latest sensor data: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToSensorDataDTO(data))
}

func (a *APIV1) ListLatestSensorDataByLocation(ctx echo.Context) *Response {
	location := ctx.QueryParam("location")
	data, err := a.sensorService.ListLatestSensorDataByLocation(location)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get latest sensor data by location: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToSensorDataListDTO(data))
}
