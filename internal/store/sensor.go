package store

import "time"

type (
	Sensor struct {
		ID               int32
		Type             string
		Name             string
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
		Name             *string
		Location         *string
		Unit             *string
		ThresholdWarning *float64
		ThresholdDanger  *float64
	}

	FindSensor struct {
		ID       *int32
		Location *string
		Limit    *int
	}

	FindSensorData struct {
		SensorID  int32
		Location  string
		StartTime *time.Time
		EndTime   *time.Time
	}

	SensorDataWithProfile struct {
		SensorID int32
		Type     string
		Name     string
		Location string
		Unit     string
		Value    float64
		Time     time.Time
	}
)

func (s *store) CreateSensor(create *Sensor) (*Sensor, error) {
	return s.driver.CreateSensor(create)
}

func (s *store) UpdateSensor(update *UpdateSensor) (*Sensor, error) {
	return s.driver.UpdateSensor(update)
}

func (s *store) ListSensors(find *FindSensor) ([]*Sensor, error) {
	return s.driver.ListSensors(find)
}

func (s store) GetSensor(find *FindSensor) (*Sensor, error) {
	return s.driver.GetSensorByID(find)
}

func (s *store) InsertSensorData(insert *SensorData) (*SensorData, error) {
	return s.driver.InsertSensorData(insert)
}

func (s *store) GetLatestSensorData(find *FindSensorData) (*SensorData, error) {
	return s.driver.GetLatestSensorData(find)
}

func (s *store) ListLatestSensorDataByLocation(find *FindSensorData) ([]*SensorDataWithProfile, error) {
	return s.driver.ListLatestSensorDataByLocation(find)
}
