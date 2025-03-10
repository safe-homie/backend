package store

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Driver interface {
	Close()
	GetDB() *pgxpool.Pool
	Migrate() error

	ListSensors(find *FindSensor) ([]*Sensor, error)
	CreateSensor(create *Sensor) (*Sensor, error)
	UpdateSensor(update *UpdateSensor) (*Sensor, error)
	InsertSensorData(insert *SensorData) (*SensorData, error)
	GetLatestSensorData(find *FindSensorData) (*SensorData, error)
	ListLatestSensorDataByLocation(find *FindSensorData) ([]*SensorDataWithProfile, error)
}
