package domain

import (
	"time"

	"github.com/safe-homie/backend/internal/store"
)

type GetDeviceResponse struct {
	ID           int32  `json:"id"`
	SerialDevice string `json:"serial_device"`
	Room         string `json:"location"`
	Name         string `json:"name"`
	Type         string `json:"type"`
}
type GetDeviceStatusResponse struct {
	DeviceID       int32                  `json:"device_id"`
	Active         bool                   `json:"active"`
	ScheduleEnable bool                   `json:"schedule_enable"`
	State          map[string]interface{} `json:"state"`
	Schedule       map[string]interface{} `json:"schedule"`
	CreatedAt      time.Time              `json:"created_time"`
	UpdatedAt      time.Time              `json:"updated_time"`
}

type GetDeviceStatusWithProfileResponse struct {
	DeviceID       int32                  `json:"device_id"`
	SerialDevice   string                 `json:"serial_device"`
	Room           string                 `json:"location"`
	Name           string                 `json:"name"`
	Type           string                 `json:"type"`
	Active         bool                   `json:"active"`
	ScheduleEnable bool                   `json:"schedule_enable"`
	State          map[string]interface{} `json:"state"`
	Schedule       map[string]interface{} `json:"schedule"`
	CreatedAt      time.Time              `json:"created_time"`
	UpdatedAt      time.Time              `json:"updated_time"`
}
type TimeRange struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}
type GetDeviceHistoryResponse struct {
	TimeRange TimeRange              `json:"time_range"`
	Data      []*store.DeviceHistory `json:"data"`
}
type ScheduleResponse struct {
	ID          int32     `json:"id"`
	DeviceID    int32     `json:"device_id"`
	Action      string    `json:"action"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Recurring   string    `json:"recurring"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type GetDeviceScheduleResponse struct {
	Schedules []*ScheduleResponse `json:"schedules"`
	Total     int                 `json:"total"`
}
type UpdateDeviceResponse struct {
	ID   int32  `json:"id"`
	Room string `json:"room"`
	Name string `json:"name"`
	Type string `json:"type"`
	// UpdatedAt time.Time `json:"updated_at"`
}
type UpdateDeviceRequest struct {
	ID   int32  `json:"id"`
	Room string `json:"room"`
	Name string `json:"name"`
	Type string `json:"type"`
}
type GetDeviceHistoryRequest struct {
	StartTime string `query:"start_time"`
	EndTime   string `query:"end_time"`
}
type ControlDeviceRequest struct {
	Action string                 `json:"action" example:"TURN_ON" enums:"TURN_ON,TURN_OFF,SET_LEVEL"`
	State  map[string]interface{} `json:"state"`
}

func ToDeviceDTO(d *store.Device) *GetDeviceResponse {
	return &GetDeviceResponse{
		ID:           d.ID,
		SerialDevice: d.SerialDevice,
		Name:         d.Name,
		Room:         d.Room,
		Type:         d.Type,
	}
}

func ToDevicesDTO(in []*store.Device) []*GetDeviceResponse {
	out := make([]*GetDeviceResponse, 0)
	for _, item := range in {
		device := ToDeviceDTO(item)
		out = append(out, device)
	}
	return out
}

func ToUpdateDeviceDTO(d *store.UpdateDevice) *UpdateDeviceResponse {
	return &UpdateDeviceResponse{
		ID:   d.ID,
		Room: *d.Room,
		Name: *d.Name,
	}
}

func ToDeviceStatusDTO(in *store.DeviceStatus) *GetDeviceStatusResponse {
	return &GetDeviceStatusResponse{
		DeviceID:       in.DeviceID,
		Active:         in.Active,
		ScheduleEnable: in.ScheduleEnable,
		Schedule: map[string]interface{}{
			"id":           in.Schedule.ID,
			"device_id":    in.Schedule.DeviceID,
			"action":       in.Schedule.Action,
			"scheduled_at": in.Schedule.ScheduledAt,
			"recurring":    in.Schedule.Recurring,
			"is_active":    in.Schedule.IsActive,
			"created_at":   in.Schedule.CreatedAt,
			"updated_at":   in.Schedule.UpdatedAt,
		},
		State:     in.State,
		CreatedAt: in.Schedule.CreatedAt,
		UpdatedAt: in.UpdatedAt,
	}
}

func ToDeviceHistoryDTO(histories []*store.DeviceHistory) *GetDeviceHistoryResponse {
	if len(histories) == 0 {
		return &GetDeviceHistoryResponse{
			TimeRange: TimeRange{
				StartTime: time.Time{},
				EndTime:   time.Time{},
			},
			Data: []*store.DeviceHistory{},
		}
	}

	startTime := histories[0].Timestamp
	endTime := histories[0].Timestamp

	for _, history := range histories {
		if history.Timestamp.Before(startTime) {
			startTime = history.Timestamp
		}
		if history.Timestamp.After(endTime) {
			endTime = history.Timestamp
		}
	}

	return &GetDeviceHistoryResponse{
		TimeRange: TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Data: histories,
	}
}

// func ToDeviceScheduleDTO(in []*store.DeviceSchedule) *GetDeviceScheduleResponse {
// 	if len(in) == 0 {
// 		return &GetDeviceScheduleResponse{
// 			Schedules: []*ScheduleResponse{},
// 			Total:     0,
// 		}
// 	}
// 	schedules := make([]*ScheduleResponse, 0, len(in))
// 	for _, item := range in {
// 		schedule := &ScheduleResponse{
// 			ID:          item.ID,
// 			DeviceID:    item.DeviceID,
// 			Action:      item.Action,
// 			ScheduledAt: item.ScheduledAt,
// 			Recurring:   item.Recurring,
// 			IsActive:    item.IsActive,
// 			CreatedAt:   item.CreatedAt,
// 		}
// 		schedules = append(schedules, schedule)
// 	}
// 	return &GetDeviceScheduleResponse{
// 		Schedules: schedules,
// 		Total:     len(in),
// 	}
// }

// type CreateScheduleRequest struct {
// 	DeviceID  int32  `json:"device_id" validate:"required,gte=0"`
// 	Action    string `json:"action" validate:"required,device_status"`
// 	StartTime string `json:"start_time" validate:"required,time_format"`
// 	EndTime   string `json:"end_time" validate:"required,time_format"`
// 	Repeat    string `json:"repeat" validate:"required,schedule_repeat"`
// 	IsActive  bool   `json:"is_active"`
// }
