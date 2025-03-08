package store

import "time"

type (
	Sensor struct {
		ID               int32
		Type             string
		Location         string
		Unit             string
		ThresholdWarning float64
		ThresholdDanger  float64
	}

	SensorData struct {
		SensorID int32
		Time     time.Time
		Value    float64
	}

	UpdateSensor struct {
		ID               int32
		Type             *string
		Location         *string
		Unit             *string
		ThresholdWarning *float64
		ThresholdDanger  *float64
	}

	FindSensorData struct {
		SensorID  int32
		StartTime *time.Time
		EndTime   *time.Time
	}
)

func (s *store) CreateSensor(create *Sensor) (*Sensor, error) {
	return s.driver.CreateSensor(create)
}

func (s *store) UpdateSensor(update *UpdateSensor) (*Sensor, error) {
	return s.driver.UpdateSensor(update)
}

func (s *store) InsertSensorData(insert *SensorData) (*SensorData, error) {
	return s.driver.InsertSensorData(insert)
}

func (s *store) GetLatestSensorData(find *FindSensorData) (*SensorData, error) {
	return s.driver.GetLatestSensorData(find)
}
