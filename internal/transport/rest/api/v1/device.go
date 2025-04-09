package v1

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/safe-homie/backend/internal/domain"
)

func (a *APIV1) GetDevice(ctx echo.Context) *Response {
	idStr := ctx.Param("id")
	if idStr == "" {
		return ErrorInvalidPathParam
	}
	idInt, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return ErrorInvalidPathParam
	}
	device, err := a.deviceService.GetDevice(int32(idInt))
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get device: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToDeviceDTO(device))
}

func (a *APIV1) ListDevices(ctx echo.Context) *Response {
	location := ctx.QueryParam("location")
	devices, err := a.deviceService.ListDevices(location)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get device list: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToDevicesDTO(devices))
}

func (a *APIV1) DeviceStatus(ctx echo.Context) *Response {
	id, errResp := ParseIDParam(ctx, "id")
	if errResp != nil {
		return errResp
	}
	status, err := a.deviceService.GetDevice(int32(id))
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to get device status: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToDeviceStatusDTO(status))
}

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

func (a *APIV1) DeviceSchedules(ctx echo.Context) *Response {
	id, errResp := ParseIDParam(ctx, "id")
	if errResp != nil {
		return errResp
	}
	schedules, err := a.deviceService.ListDeviceSchedule(int32(id))
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to list device schedules: "+err.Error())
	}
	return NewDataResponse(http.StatusOK, domain.ToDeviceScheduleDTO(schedules))
}

func (a *APIV1) CreateDevice(ctx echo.Context) *Response {
	var req domain.CreateDeviceRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorInvalidRequestBody
	}
	if err := a.context.Validator().Validate(req); err != nil {
		return ErrorInvalidRequestBody
	}
	device, err := a.deviceService.CreateDevice(&req)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to create device: "+err.Error())
	}
	return NewDataResponse(http.StatusCreated, domain.ToDeviceDTO(device))

}

func (a *APIV1) UpdateDevice(ctx echo.Context) *Response {
	idStr := ctx.Param("id")
	if idStr == "" {
		return ErrorInvalidPathParam
	}
	idInt, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return ErrorInvalidPathParam
	}

	var req domain.UpdateDeviceRequest
	if err := ctx.Bind(&req); err != nil {
		return ErrorInvalidRequestBody
	}
	if err := a.context.Validator().Validate(req); err != nil {
		return ErrorInvalidRequestBody
	}

	device, err := a.deviceService.UpdateDevice(int32(idInt), &req)
	if err != nil {
		return NewErrorResponse(http.StatusInternalServerError, "failed to update device: "+err.Error())
	}

	return NewDataResponse(http.StatusOK, domain.ToUpdateDeviceDTO(device))
}

func (a *APIV1) ControlDevice(ctx echo.Context) *Response {
	// idStr := ctx.Param("id")
	// if idStr == "" {
	// 	return ErrorInvalidPathParam
	// }
	// idInt, err := strconv.ParseInt(idStr, 10, 32)
	// if err != nil {
	// 	return ErrorInvalidPathParam
	// }

	// var req domain.ControlDeviceRequest
	// if err := ctx.Bind(&req); err != nil {
	// 	return ErrorInvalidRequestBody
	// }
	// if err := a.context.Validator().Validate(req); err != nil {
	// 	return ErrorInvalidRequestBody
	// }

	// result, err := a.deviceService.ControlDevice(int32(idInt), &req)
	// if err != nil {
	// 	return NewErrorResponse(http.StatusInternalServerError, "failed to control device: "+err.Error())
	// }
	// return NewDataResponse(http.StatusOK, domain.ToControlResultDTO(result))
	return nil
}
func (a *APIV1) ScheduleDevice(ctx echo.Context) *Response {
	// idStr := ctx.Param("id")
	// if idStr == "" {
	// 	return ErrorInvalidPathParam
	// }
	// idInt, err := strconv.ParseInt(idStr, 10, 32)
	// if err != nil {
	// 	return ErrorInvalidPathParam
	// }

	// var req domain.ScheduleDeviceRequest
	// if err := ctx.Bind(&req); err != nil {
	// 	return ErrorInvalidRequestBody
	// }
	// if err := a.context.Validator().Validate(req); err != nil {
	// 	return ErrorInvalidRequestBody
	// }

	// schedule, err := a.deviceService.CreateDeviceSchedule(int32(idInt), &req)
	// if err != nil {
	// 	return NewErrorResponse(http.StatusInternalServerError, "failed to schedule device: "+err.Error())
	// }
	// return NewDataResponse(http.StatusCreated, domain.ToDeviceScheduleDTO(schedule))
	return nil
}

func (a *APIV1) UpdateScheduleDevice(ctx echo.Context) *Response {
	// idStr := ctx.Param("id")
	// if idStr == "" {
	// 	return ErrorInvalidPathParam
	// }
	// idInt, err := strconv.ParseInt(idStr, 10, 32)
	// if err != nil {
	// 	return ErrorInvalidPathParam
	// }

	// scheduleIdStr := ctx.Param("schedule_id")
	// if scheduleIdStr == "" {
	// 	return ErrorInvalidPathParam
	// }
	// scheduleIdInt, err := strconv.ParseInt(scheduleIdStr, 10, 32)
	// if err != nil {
	// 	return ErrorInvalidPathParam
	// }

	// var req domain.UpdateScheduleRequest
	// if err := ctx.Bind(&req); err != nil {
	// 	return ErrorInvalidRequestBody
	// }
	// if err := a.context.Validator().Validate(req); err != nil {
	// 	return ErrorInvalidRequestBody
	// }

	// schedule, err := a.deviceService.UpdateDeviceSchedule(int32(idInt), int32(scheduleIdInt), &req)
	// if err != nil {
	// 	return NewErrorResponse(http.StatusInternalServerError, "failed to update device schedule: "+err.Error())
	// }
	// return NewDataResponse(http.StatusOK, domain.ToDeviceScheduleDTO(schedule))
	return nil
}

func (a *APIV1) DeleteScheduleDevice(ctx echo.Context) *Response {
	// idStr := ctx.Param("id")
	// if idStr == "" {
	// 	return ErrorInvalidPathParam
	// }
	// idInt, err := strconv.ParseInt(idStr, 10, 32)
	// if err != nil {
	// 	return ErrorInvalidPathParam
	// }

	// scheduleIdStr := ctx.Param("schedule_id")
	// if scheduleIdStr == "" {
	// 	return ErrorInvalidPathParam
	// }
	// scheduleIdInt, err := strconv.ParseInt(scheduleIdStr, 10, 32)
	// if err != nil {
	// 	return ErrorInvalidPathParam
	// }

	// if err := a.deviceService.DeleteDeviceSchedule(int32(idInt), int32(scheduleIdInt)); err != nil {
	// 	return NewErrorResponse(http.StatusInternalServerError, "failed to delete device schedule: "+err.Error())
	// }
	// return NewDataResponse(http.StatusNoContent, nil)
	return nil
}

func (a *APIV1) AutoModeSchedule(ctx echo.Context) *Response {
	return nil
}

func (a *APIV1) PresetModeSchedule(ctx echo.Context) *Response {
	return nil
}

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
