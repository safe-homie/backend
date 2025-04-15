package store

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Driver interface {
	Close()
	GetDB() *pgxpool.Pool
	Migrate() error

	ListSensors(find *FindSensor) ([]*Sensor, error)
	GetSensorByID(find *FindSensor) (*Sensor, error)
	CreateSensor(create *Sensor) (*Sensor, error)
	UpdateSensor(update *UpdateSensor) (*Sensor, error)
	InsertSensorData(insert *SensorData) (*SensorData, error)
	GetLatestSensorData(find *FindSensorData) (*SensorData, error)
	ListLatestSensorDataByLocation(find *FindSensorData) ([]*SensorDataWithProfile, error)

	GetDeviceByID(find *FindDevice) (*Device, error)
	ListDevices(find *FindDevice) ([]*Device, error)
	UpdateDevice(edit *UpdateDevice) (*UpdateDevice, error)
	InsertStatusData(insert *DeviceStatus) error
	GetDeviceStatus(find *FindDevice) (*DeviceStatus, error)
	UpdateDeviceStatus(status *DeviceStatus) (*DeviceStatus, error)
	CreateDeviceHistory(create *DeviceHistory) (*DeviceHistory, error)
	ListDeviceHistory(find *FindDeviceHistory) ([]*DeviceHistory, error)

	// GetStatus(deviceID int32) (*DeviceStatus, error)
	// UpdateStatus(status *DeviceStatus) error
	// RecordHistory(history *DeviceHistory) error

	// ListDeviceSchedule(find *FindDeviceSchedule) ([]*DeviceSchedule, error)
	// CreateDeviceSchedule(create *DeviceSchedule) (*DeviceSchedule, error)
	// UpdateDeviceSchedule(edit *UpdateDeviceSchedule) (*DeviceSchedule, error)
	// DeleteDeviceSchedule(find *FindDevice) error
	SaveToken(save *Token) (*Token, error)
	GetToken(find *FindToken) (*Token, error)
}
