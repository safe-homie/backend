package v1

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/safe-homie/backend/internal/domain"
)

// @Summary		Get a sensor profile
// @Description	Get sensor profile by ID
// @Tags			sensors
// @Produce		json
// @Param			id	path		int	true	"Sensor ID"
// @Success		200	{object}	SuccessResponseWrapper{data=domain.GetSensorResponse}
// @Failure		400	{object}	ErrorResponseWrapper{error=string}
// @Failure		500	{object}	ErrorResponseWrapper{error=string}
// @Router			/sensors/{id} [get]
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

// @Summary		Get sensor profiles list
// @Description	Get sensor profile lists by location
// @Tags			sensors
// @Produce		json
// @Param			location	query		string	false	"list sensors by location"
// @Success		200			{object}	SuccessResponseWrapper{data=[]domain.GetSensorResponse}
// @Failure		500			{object}	ErrorResponseWrapper{error=string}
// @Router			/sensors [get]
func (a *APIV1) ListSensors(ctx echo.Context) *Response {
	location := ctx.QueryParam("location")
	sensors, err := a.sensorService.ListSensors(location)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get sensor list: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToSensorsDTO(sensors))
}

// @Summary		Create new sensor profile
// @Description	Create a new sensor profile
// @Tags			sensors
// @Accept			json
// @Produce		json
// @Param			request	body		domain.CreateSensorRequest	true	"create sensor request body"
// @Success		200		{object}	SuccessResponseWrapper{data=domain.GetSensorResponse}
// @Failure		400		{object}	ErrorResponseWrapper{error=string}
// @Failure		500		{object}	ErrorResponseWrapper{error=string}
// @Router			/sensors [post]
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

// @Summary		List sensor profile and latest data
// @Description	List sensor profile and latest data by location
// @Tags			sensors
// @Produce		json
// @Param			location	query		string	false	"list sensors by location"
// @Success		200			{object}	SuccessResponseWrapper{data=[]domain.GetSensorDataWithProfileResponse}
// @Failure		500			{object}	ErrorResponseWrapper{error=string}
// @Router			/sensors/data/latest [get]
func (a *APIV1) ListLatestSensorDataByLocation(ctx echo.Context) *Response {
	location := ctx.QueryParam("location")
	data, err := a.sensorService.ListLatestSensorDataByLocation(location)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get latest sensor data by location: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToSensorDataListDTO(data))
}
