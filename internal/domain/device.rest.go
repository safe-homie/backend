package domain

import (
	"time"

	"github.com/safe-homie/backend/internal/store"
)

type OperationStatus struct {
	Power string `json:"power,omitempty"`
	Level string `json:"level,omitempty"`
	// Mode  *string `json:"mode,omitempty"`
}

type TimeRange struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type ScheduleResponse struct {
	ID          int32     `json:"id"`
	DeviceID    int32     `json:"device_id"`
	Action      string    `json:"action"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Recurring   string    `json:"recurring"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetDeviceResponse struct {
	ID         int32           `json:"id"`
	Location   string          `json:"location"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Attributes OperationStatus `json:"attributes"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

type GetDeviceStatusResponse struct {
	ID         int32           `json:"id"`
	Name       string          `json:"name"`
	Type       string          `json:"type"`
	Attributes OperationStatus `json:"attributes"`
}

type GetDeviceHistoryResponse struct {
	TimeRange TimeRange              `json:"time_range"`
	Data      []*store.DeviceHistory `json:"data"`
}

type GetDeviceScheduleResponse struct {
	Schedules []*ScheduleResponse `json:"schedules"`
	Total     int                 `json:"total"`
}

type PutDeviceResponse struct {
	ID        int32     `json:"id"`
	Location  string    `json:"location"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetDeviceHistoryRequest struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type CreateDeviceRequest struct {
	Location string `json:"location" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Power    string `json:"power" binding:"required"`
	Level    string `json:"level" binding:"required"`
}

type UpdateDeviceRequest struct {
	Location string `json:"location"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}

type CreateScheduleRequest struct {
	DeviceID  int32  `json:"device_id" validate:"required,gte=0"`
	Action    string `json:"action" validate:"required,device_status"`
	StartTime string `json:"start_time" validate:"required,time_format"`
	EndTime   string `json:"end_time" validate:"required,time_format"`
	Repeat    string `json:"repeat" validate:"required,schedule_repeat"`
	IsActive  bool   `json:"is_active"`
}

func ToDeviceDTO(d *store.Device) *GetDeviceResponse {
	return &GetDeviceResponse{
		ID:       d.ID,
		Location: d.Location,
		Name:     d.Name,
		Type:     d.Type,
		Attributes: OperationStatus{
			Power: d.Power,
			Level: d.Level,
		},
		CreatedAt: d.CreateAt,
		UpdatedAt: d.UpdateAt,
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

func ToUpdateDeviceDTO(d *store.UpdateDevice) *PutDeviceResponse {
	return &PutDeviceResponse{
		ID:        d.ID,
		Location:  *d.Location,
		Name:      *d.Name,
		Type:      *d.Type,
		UpdatedAt: *d.UpdatedAt,
	}
}

func ToDeviceStatusDTO(in *store.Device) *GetDeviceStatusResponse {
	return &GetDeviceStatusResponse{
		ID:   in.ID,
		Name: in.Name,
		Type: in.Type,
		Attributes: OperationStatus{
			Power: in.Power,
			Level: in.Level,
		},
	}
}

func ToDeviceHistoryDTO(in []*store.DeviceHistory) *GetDeviceHistoryResponse {
	if len(in) <= 0 {
		return &GetDeviceHistoryResponse{
			TimeRange: TimeRange{
				StartTime: time.Time{},
				EndTime:   time.Time{},
			},
			Data: []*store.DeviceHistory{},
		}
	}
	startTime := in[0].Timestamp
	endTime := in[0].Timestamp
	for _, item := range in {
		if item.Timestamp.Before(startTime) {
			startTime = item.Timestamp
		}
		if item.Timestamp.After(endTime) {
			endTime = item.Timestamp
		}
	}
	return &GetDeviceHistoryResponse{
		TimeRange: TimeRange{
			StartTime: startTime,
			EndTime:   endTime,
		},
		Data: in,
	}
}

func ToDeviceScheduleDTO(in []*store.DeviceSchedule) *GetDeviceScheduleResponse {
	if len(in) == 0 {
		return &GetDeviceScheduleResponse{
			Schedules: []*ScheduleResponse{},
			Total:     0,
		}
	}
	schedules := make([]*ScheduleResponse, 0, len(in))
	for _, item := range in {
		schedule := &ScheduleResponse{
			ID:          item.ID,
			DeviceID:    item.DeviceID,
			Action:      item.Action,
			ScheduledAt: item.ScheduledAt,
			Recurring:   item.Recurring,
			IsActive:    item.IsActive,
			CreatedAt:   item.CreatedAt,
		}
		schedules = append(schedules, schedule)
	}
	return &GetDeviceScheduleResponse{
		Schedules: schedules,
		Total:     len(in),
	}
}
