package store

import (
	"fmt"
	"time"

	"github.com/safe-homie/backend/pkg/mqtt"
)

type (
	Device struct {
		ID           int32
		SerialDevice string
		Name         string
		Type         string
		Room         string
	}
	DeviceStatus struct {
		DeviceID       int32
		Active         bool                   // active = true; inactive = false
		ScheduleEnable bool                   // enable = true; disable = false
		State          map[string]interface{} //power ON/OF ; level MIN/MEDIUM/MAX
		Schedule       DeviceSchedule
		UpdatedAt      time.Time
	}
	DeviceHistory struct {
		ID        int32                  `json:"id"`
		DeviceID  int32                  `json:"device_id"`
		Action    string                 `json:"action"`
		State     map[string]interface{} `json:"state"` //power ON/OF ; level MIN/MEDIUM/MAX
		By        string                 `json:"by"`    // user or scheduling:%ID
		Timestamp time.Time              `json:"timestamp"`
	}
	DeviceSchedule struct {
		ID          int32
		DeviceID    int32
		Action      string // "TURN_ON";"TURN_OFF";"SET_LEVEL"
		ScheduledAt time.Time
		Recurring   string // "DAILY", "WEEKLY", "ONCE"
		IsActive    bool
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
	FindDevice struct {
		ID    *int32
		Room  *string
		Limit *int
	}
	FindDeviceHistory struct {
		DeviceID  *int32
		StartTime *time.Time
		EndTime   *time.Time
		Limit     *int
	}
	FindDeviceSchedule struct {
		DeviceID *int32
		Limit    *int
	}
	UpdateDevice struct {
		ID   int32
		Room *string
		Name *string
		Type *string
	}
	UpdateDeviceSchedule struct {
		ID          int32
		DeviceID    int32
		Action      *string
		ScheduledAt *time.Time
		Recurring   *string
		UpdatedAt   *time.Time
	}
)

func (s *store) GetDevice(find *FindDevice) (*Device, error) {
	return s.driver.GetDeviceByID(find)
}

func (s *store) ListDevices(find *FindDevice) ([]*Device, error) {
	return s.driver.ListDevices(find)
}

func (s *store) UpdateDevice(edit *UpdateDevice) (*UpdateDevice, error) {
	return s.driver.UpdateDevice(edit)
}

func (s *store) InsertStatusData(insert *DeviceStatus) error {
	return s.driver.InsertStatusData(insert)
}

func (s *store) GetStatus(find *FindDevice) (*DeviceStatus, error) {
	return s.driver.GetDeviceStatus(find)
}
func (s *store) UpdateStatus(status *DeviceStatus) (*DeviceStatus, error) {
	return s.driver.UpdateDeviceStatus(status)
}
func (s *store) RecordHistory(history *DeviceHistory) (*DeviceHistory, error) {
	return s.driver.CreateDeviceHistory(history)
}

func (s *store) ListDeviceHistory(find *FindDeviceHistory) ([]*DeviceHistory, error) {
	return s.driver.ListDeviceHistory(find)
}

// func (s *store) ListDeviceSchedule(find *FindDeviceSchedule) ([]*DeviceSchedule, error) {
// 	return s.driver.ListDeviceSchedule(find)
// }

// func (s *store) CreateDeviceSchedule(create *DeviceSchedule) (*DeviceSchedule, error) {
// 	return s.driver.CreateDeviceSchedule(create)
// }

// func (s *store) UpdateDeviceSchedule(edit *UpdateDeviceSchedule) (*DeviceSchedule, error) {
// 	return s.driver.UpdateDeviceSchedule(edit)
// }

// func (s *store) DeleteDeviceSchedule(find *FindDevice) error {
// 	return s.driver.DeleteDeviceSchedule(find)
// }

// type DeviceCommand interface {
// 	Execute() error
// 	GetDeviceID() int32
// 	GetCommandType() string
// 	GetState() map[string]interface{}
// }

type BaseCommand struct {
	DeviceID    int32
	CommandType string
	State       map[string]interface{}
	MQTT        mqtt.MQTTClient
}

func (c *BaseCommand) GetDeviceID() int32 {
	return c.DeviceID
}

func (c *BaseCommand) GetCommandType() string {
	return c.CommandType
}

func (c *BaseCommand) GetState() map[string]interface{} {
	return c.State
}

type TurnOnCommand struct {
	BaseCommand
}

func NewTurnOnCommand(deviceID int32, client mqtt.MQTTClient) *TurnOnCommand {
	return &TurnOnCommand{
		BaseCommand: BaseCommand{
			DeviceID:    deviceID,
			CommandType: "TURN_ON",
			State:       map[string]interface{}{"power": "ON"},
			MQTT:        client,
		},
	}
}

func (c *TurnOnCommand) Execute() error {
	topic := fmt.Sprintf("devices/%d/command", c.DeviceID)
	payload := map[string]interface{}{
		"command": c.CommandType,
		"state":   c.State,
	}
	return c.MQTT.Publish(topic, payload)
}

type TurnOffCommand struct {
	BaseCommand
}

func NewTurnOffCommand(deviceID int32, client mqtt.MQTTClient) *TurnOffCommand {
	return &TurnOffCommand{
		BaseCommand: BaseCommand{
			DeviceID:    deviceID,
			CommandType: "TURN_OFF",
			State:       map[string]interface{}{"power": "OFF"},
			MQTT:        client,
		},
	}
}

func (c *TurnOffCommand) Execute() error {
	topic := fmt.Sprintf("devices/%d/command", c.DeviceID)
	payload := map[string]interface{}{
		"command": c.CommandType,
		"state":   c.State,
	}
	return c.MQTT.Publish(topic, payload)
}

type SetLevelCommand struct {
	BaseCommand
	Level string
}

func NewSetLevelCommand(deviceID int32, level string, client mqtt.MQTTClient) *SetLevelCommand {
	return &SetLevelCommand{
		BaseCommand: BaseCommand{
			DeviceID:    deviceID,
			CommandType: "SET_LEVEL",
			State:       map[string]interface{}{"level": level},
			MQTT:        client,
		},
		Level: level,
	}
}

func (c *SetLevelCommand) Execute() error {
	topic := fmt.Sprintf("devices/%d/command", c.DeviceID)
	payload := map[string]interface{}{
		"command": c.CommandType,
		"state":   c.State,
	}
	return c.MQTT.Publish(topic, payload)
}
