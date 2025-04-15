package v1

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/safe-homie/backend/internal/domain"
	"github.com/safe-homie/backend/internal/store"
)

// @Summary		Get a device profile
// @Description	Get device profile by ID
// @Tags			devices
// @Produce		json
// @Param			id	path		int	true	"Device ID"
// @Success		200	{object}	SuccessResponseWrapper{data=domain.GetDeviceResponse}
// @Failure		400	{object}	ErrorResponseWrapper{error=string}
// @Failure		500	{object}	ErrorResponseWrapper{error=string}
// @Router			/devices/{id} [get]
func (a *APIV1) GetDevice(ctx echo.Context) *Response {
	id, errResp := ParseIDParam(ctx, "id")
	if errResp != nil {
		return errResp
	}
	device, err := a.deviceService.GetDevice(int32(id))
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get device: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToDeviceDTO(device))
}

// @Summary		Get device profiles list
// @Description	Get device profile lists by location
// @Tags			devices
// @Produce		json
// @Param			room	query		string	false	"list devices by room"
// @Success		200			{object}	SuccessResponseWrapper{data=[]domain.GetDeviceResponse}
// @Failure		500			{object}	ErrorResponseWrapper{error=string}
// @Router			/devices [get]
func (a *APIV1) ListDevices(ctx echo.Context) *Response {
	room := ctx.QueryParam("room")
	devices, err := a.deviceService.ListDevices(room)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get device list: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToDevicesDTO(devices))
}

// @Summary		Update a device
// @Description	Update device information by ID
// @Tags			devices
// @Accept			json
// @Produce		json
// @Param			id		path		int							true	"Device ID"
// @Param			body	body		domain.UpdateDeviceRequest	true	"Updated device info"
// @Success		200		{object}	SuccessResponseWrapper{data=domain.UpdateDeviceResponse}
// @Failure		400		{object}	ErrorResponseWrapper{error=string}
// @Failure		500		{object}	ErrorResponseWrapper{error=string}
// @Router			/devices/{id} [put]
func (a *APIV1) UpdateDevice(ctx echo.Context) *Response {
	id, errResp := ParseIDParam(ctx, "id")
	if errResp != nil {
		return errResp
	}

	var req domain.UpdateDeviceRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorInvalidRequestBody
	}
	if err := a.context.Validator().Validate(req); err != nil {
		return ErrorInvalidRequestBody
	}

	device, err := a.deviceService.UpdateDevice(int32(id), &req)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to update device: "+err.Error())
	}

	return NewDataResponse(http.StatusOK, domain.ToUpdateDeviceDTO(device))
}

func (a *APIV1) DeviceStatus(ctx echo.Context) *Response {
	id, errResp := ParseIDParam(ctx, "id")
	if errResp != nil {
		return errResp
	}
	status, err := a.deviceService.GetDeviceStatus(int32(id))
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get device status: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToDeviceStatusDTO(status))
}

func (a *APIV1) ControlDevice(ctx echo.Context) *Response {

	// status, err := a.deviceService.GetDeviceStatus(int32(id))
	// if err != nil {
	// 	return err
	// }
	id, errResp := ParseIDParam(ctx, "id")
	if errResp != nil {
		return errResp
	}
	var req domain.ControlDeviceRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorInvalidRequestBody
	}
	if err := a.context.Validator().Validate(req); err != nil {
		return ErrorInvalidRequestBody
	}
	var cmd domain.DeviceCommand
	switch req.Action {
	case domain.DeviceActionTurnON:
		cmd = store.NewTurnOnCommand(id)
	case domain.DeviceActionTurnOFF:
		cmd = store.NewTurnOffCommand(id)
	case domain.DeviceActionSetLEVEL:
		if level, ok := req.State["level"].(string); ok {
			cmd = store.NewSetLevelCommand(id, level)
		}
	default:
		return NewErrorResponse(http.StatusInternalServerError, "failed to action device: ")
	}

	if err := a.deviceService.ExecuteCommand(cmd); err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to control device: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, "ĐÃ THÀNH CÔNG")
}

// func (h *RESTHandler) ControlDevice(c *gin.Context) {
// 	var request struct {
// 		Action string                 `json:"action"`
// 		State  map[string]interface{} `json:"state"`
// 	}

// 	if err := c.BindJSON(&request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	deviceID := c.Param("id")
// 	// Convert deviceID to int32

// 	var cmd domain.DeviceCommand

// 	switch request.Action {
// 	case "TURN_ON":
// 		cmd = application.NewTurnOnCommand(deviceID)
// 	case "TURN_OFF":
// 		cmd = application.NewTurnOffCommand(deviceID)
// 	case "SET_LEVEL":
// 		if level, ok := request.State["level"].(string); ok {
// 			cmd = application.NewSetLevelCommand(deviceID, level)
// 		}
// 	// Add more commands as needed
// 	default:
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown action"})
// 		return
// 	}

// 	if err := h.service.ExecuteCommand(c.Request.Context(), cmd); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"status": "success"})
// }

// @Summary		Get device history
// @Description	Get historical data of a device by ID
// @Tags			devices
// @Accept			json
// @Produce		json
// @Param			id		path		int								true	"Device ID"
// @Param			body	body		domain.GetDeviceHistoryRequest	true	"Filter for device history"
// @Success		200		{object}	SuccessResponseWrapper{data=[]domain.DeviceHistoryResponse}
// @Failure		400		{object}	ErrorResponseWrapper{error=string}
// @Failure		500		{object}	ErrorResponseWrapper{error=string}
// @Router			/devices/{id}/history [post]
func (a *APIV1) DeviceHistories(ctx echo.Context) *Response {
	id, errResp := ParseIDParam(ctx, "id")
	if errResp != nil {
		return errResp
	}
	var req domain.GetDeviceHistoryRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorInvalidRequestBody
	}
	if err := a.context.Validator().Validate(req); err != nil {
		return ErrorInvalidRequestBody
	}
	histories, err := a.deviceService.ListDeviceHistory(int32(id), &req)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get device history: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToDeviceHistoryDTO(histories))
}

// func (a *APIV1) DeviceSchedules(ctx echo.Context) *Response {
// 	id, errResp := ParseIDParam(ctx, "id")
// 	if errResp != nil {
// 		return errResp
// 	}
// 	schedules, err := a.deviceService.ListDeviceSchedule(int32(id))
// 	if err != nil {
// 		return NewErrorResponse(http.StatusInternalServerError, "failed to list device schedules: "+err.Error())
// 	}
// 	return NewDataResponse(http.StatusOK, domain.ToDeviceScheduleDTO(schedules))
// }

// func (a *APIV1) ScheduleDevice(ctx echo.Context) *Response {
// 	// idStr := ctx.Param("id")
// 	// if idStr == "" {
// 	// 	return ErrorInvalidPathParam
// 	// }
// 	// idInt, err := strconv.ParseInt(idStr, 10, 32)
// 	// if err != nil {
// 	// 	return ErrorInvalidPathParam
// 	// }

// 	// var req domain.ScheduleDeviceRequest
// 	// if err := ctx.Bind(&req); err != nil {
// 	// 	return ErrorInvalidRequestBody
// 	// }
// 	// if err := a.context.Validator().Validate(req); err != nil {
// 	// 	return ErrorInvalidRequestBody
// 	// }

// 	// schedule, err := a.deviceService.CreateDeviceSchedule(int32(idInt), &req)
// 	// if err != nil {
// 	// 	return NewErrorResponse(http.StatusInternalServerError, "failed to schedule device: "+err.Error())
// 	// }
// 	// return NewDataResponse(http.StatusCreated, domain.ToDeviceScheduleDTO(schedule))
// 	return nil
// }

// func (a *APIV1) UpdateScheduleDevice(ctx echo.Context) *Response {
// 	// idStr := ctx.Param("id")
// 	// if idStr == "" {
// 	// 	return ErrorInvalidPathParam
// 	// }
// 	// idInt, err := strconv.ParseInt(idStr, 10, 32)
// 	// if err != nil {
// 	// 	return ErrorInvalidPathParam
// 	// }

// 	// scheduleIdStr := ctx.Param("schedule_id")
// 	// if scheduleIdStr == "" {
// 	// 	return ErrorInvalidPathParam
// 	// }
// 	// scheduleIdInt, err := strconv.ParseInt(scheduleIdStr, 10, 32)
// 	// if err != nil {
// 	// 	return ErrorInvalidPathParam
// 	// }

// 	// var req domain.UpdateScheduleRequest
// 	// if err := ctx.Bind(&req); err != nil {
// 	// 	return ErrorInvalidRequestBody
// 	// }
// 	// if err := a.context.Validator().Validate(req); err != nil {
// 	// 	return ErrorInvalidRequestBody
// 	// }

// 	// schedule, err := a.deviceService.UpdateDeviceSchedule(int32(idInt), int32(scheduleIdInt), &req)
// 	// if err != nil {
// 	// 	return NewErrorResponse(http.StatusInternalServerError, "failed to update device schedule: "+err.Error())
// 	// }
// 	// return NewDataResponse(http.StatusOK, domain.ToDeviceScheduleDTO(schedule))
// 	return nil
// }

// func (a *APIV1) DeleteScheduleDevice(ctx echo.Context) *Response {
// 	// idStr := ctx.Param("id")
// 	// if idStr == "" {
// 	// 	return ErrorInvalidPathParam
// 	// }
// 	// idInt, err := strconv.ParseInt(idStr, 10, 32)
// 	// if err != nil {
// 	// 	return ErrorInvalidPathParam
// 	// }

// 	// scheduleIdStr := ctx.Param("schedule_id")
// 	// if scheduleIdStr == "" {
// 	// 	return ErrorInvalidPathParam
// 	// }
// 	// scheduleIdInt, err := strconv.ParseInt(scheduleIdStr, 10, 32)
// 	// if err != nil {
// 	// 	return ErrorInvalidPathParam
// 	// }

// 	// if err := a.deviceService.DeleteDeviceSchedule(int32(idInt), int32(scheduleIdInt)); err != nil {
// 	// 	return NewErrorResponse(http.StatusInternalServerError, "failed to delete device schedule: "+err.Error())
// 	// }
// 	// return NewDataResponse(http.StatusNoContent, nil)
// 	return nil
// }

// func (a *APIV1) AutoModeSchedule(ctx echo.Context) *Response {
// 	return nil
// }

// func (a *APIV1) PresetModeSchedule(ctx echo.Context) *Response {
// 	return nil
// }

func ParseIDParam(ctx echo.Context, paramName string) (int32, *Response) {
	idStr := ctx.Param(paramName)
	if idStr == "" {
		return 0, ErrorInvalidPathParam
	}
	idInt, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return 0, ErrorInvalidPathParam
	}
	return int32(idInt), nil
}
